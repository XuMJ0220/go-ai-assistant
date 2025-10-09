package service

import (
	"fmt"
	"go-ai-assistant/core"
	"go-ai-assistant/models"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
)

// CreateKnowledgeBaseInput 定义了创建知识库的输入
type CreateKnowledgeBaseInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// CreateKnowledgeBase 处理创建知识库的逻辑
func CreateKnowledgeBase(input CreateKnowledgeBaseInput, userID uint) (models.KnowledgeBase, error) {
	// 创建一个唯一的集合名称，用于向量数据库
	collectionName := fmt.Sprintf("kb_%s", uuid.New().String())

	kb := models.KnowledgeBase{
		UserID:                 userID,
		Name:                   input.Name,
		Description:            input.Description,
		VectorDBCollectionName: collectionName,
	}

	// 步骤 1: 先创建知识库记录
	if err := core.DB.Create(&kb).Error; err != nil {
		return models.KnowledgeBase{}, err
	}

	// 步骤 2: 创建成功后，使用 Preload 重新加载这条记录及其关联的 User 数据
	if err := core.DB.Preload("User").First(&kb, kb.ID).Error; err != nil {
		// 这里的错误通常是数据库连接问题，因为记录肯定存在
		return models.KnowledgeBase{}, err
	}

	return kb, nil
}

// ListKnowledgeBases 获取指定用户的所有知识库，并预加载用户信息
func ListKnowledgeBases(userID uint) ([]models.KnowledgeBase, error) {
	var knowledgeBases []models.KnowledgeBase

	// 在这里添加 .Preload("User")
	if err := core.DB.Preload("User").Where("user_id = ?", userID).Find(&knowledgeBases).Error; err != nil {
		return nil, err
	}

	return knowledgeBases, nil
}

func ProcessAndEmbedDocument(doc models.Document) {
	// 记录后台任务开始
	log.Printf("[+] Starting to process document ID: %d, Path: %s", doc.ID, doc.FilePath)
	// 更新文档状态为“处理中”
	core.DB.Model(&doc).Update("Status", "PROCESSING")

	// 1. 读取文件内容
	textContent, err := readPdfText(doc.FilePath)
	if err != nil {
		log.Printf("[!] Error reading document content for doc ID %d: %v", doc.ID, err)
		core.DB.Model(&doc).Update("Status", "FAILED")
		return
	}
	if strings.TrimSpace(textContent) == "" {
		log.Printf("[-] Document ID %d has no content.", doc.ID)
		core.DB.Model(&doc).Update("Status", "DONE")
		return
	}

	// 2. **【在这里定义 splitter】** 创建文本切分器
	// 使用 LangChainGo 的递归字符切分器
	splitter := textsplitter.NewRecursiveCharacter(
		textsplitter.WithChunkSize(500),
		textsplitter.WithChunkOverlap(50),
	)
	// 使用 splitter 对文本进行切块
	chunks, err := splitter.SplitText(textContent)
	if err != nil || len(chunks) == 0 {
		log.Printf("[!] Error splitting text for document ID %d: %v", doc.ID, err)
		core.DB.Model(&doc).Update("Status", "FAILED")
		return
	}
	log.Printf("[+] Document ID %d split into %d chunks.", doc.ID, len(chunks))

	// 3. 调用 AI 服务创建向量
	_, err = CreateEmbeddings(chunks) // 我们暂时只调用，不使用返回值
	if err != nil {
		log.Printf("[!] Error creating embeddings for document ID %d: %v", doc.ID, err)
		core.DB.Model(&doc).Update("Status", "FAILED")
		return
	}
	log.Printf("[+] Successfully created embeddings for %d chunks.", len(chunks))

	// 向量已生成，但我们暂时不在这里做任何存储操作
	// 因为我们已经删除了 ChromaDB 的逻辑，等待下一步集成 pgvector

	// 4. 暂时直接标记为完成
	core.DB.Model(&doc).Update("Status", "DONE")
	log.Printf("[+] (Placeholder) Successfully processed document ID: %d", doc.ID)
}

// readPdfText 从指定路径的PDF文件中提取所有文本内容
// 参数:
//   - path: PDF文件的完整路径
//
// 返回值:
//   - string: 提取的文本内容，各页面之间用换行符分隔
//   - error: 如果读取过程中出现错误则返回错误信息
func readPdfText(path string) (string, error) {
	// 打开PDF文件
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close() // 确保文件在函数结束时关闭

	// 创建PDF读取器
	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		return "", err
	}

	// 获取PDF文件的总页数
	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return "", err
	}

	// 使用字符串构建器来高效拼接文本
	var textBuilder strings.Builder

	// 遍历PDF的每一页
	for i := 1; i <= numPages; i++ {
		// 获取指定页面
		page, err := pdfReader.GetPage(i)
		if err != nil {
			return "", err
		}

		// 创建文本提取器
		ex, err := extractor.New(page)
		if err != nil {
			return "", err
		}

		// 从当前页面提取文本
		text, err := ex.ExtractText()
		if err != nil {
			return "", err
		}

		// 将提取的文本添加到构建器中
		textBuilder.WriteString(text)
		textBuilder.WriteString("\n") // 在每页之间添加换行符
	}

	// 返回完整的文本内容
	return textBuilder.String(), nil
}

// UploadDocument 处理文档上传的逻辑
func UploadDocument(kbID uint, fileHeader *multipart.FileHeader) (models.Document, error) {
	// 1.检查知识库是否存在
	var kb models.KnowledgeBase
	if err := core.DB.First(&kb, kbID).Error; err != nil {
		return models.Document{}, fmt.Errorf("knowledge base with id %d not found", kbID)
	}
	// 2.创建用于储存知识库文件的目录（如果不存在）
	// 目录结构: ./uploads/{knowledge_base_id}/
	uploadDir := filepath.Join("uploads", strconv.FormatUint(uint64(kbID), 10))
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return models.Document{}, err
	}
	// 3.将文件保存到服务器
	filePath := filepath.Join(uploadDir, fileHeader.Filename)
	if err := os.WriteFile(filePath, []byte{}, 0644); err != nil {
		// 先创建一个空文件
		// ... (这里省略了实际的文件写入逻辑，因为我们需要 gin.Context 来保存文件)
		// 完整的逻辑将在 API 层实现
	}
	// 4.在数据库中创建文档记录
	doc := models.Document{
		KnowledgeBaseID: kbID,
		FileName:        fileHeader.Filename,
		FilePath:        filePath,
		FileSize:        fileHeader.Size,
		Status:          "UPLOADED", // 初始化状态为已上传
	}
	if err := core.DB.Create(&doc).Error; err != nil {
		return models.Document{}, err
	}

	// (未来) 在这里触发异步任务，对文档进行切块和向量化
	return doc, nil
}
