package api

import (
	"go-ai-assistant/core"
	"go-ai-assistant/models"
	"go-ai-assistant/service"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateKnowledgeBase 是创建知识库的 Handler
func CreateKnowledgeBase(c *gin.Context) {
	var input service.CreateKnowledgeBaseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从认证中间件中获取 userID
	userIDFloat, _ := c.Get("user_id")
	userID := uint(userIDFloat.(float64))

	kb, err := service.CreateKnowledgeBase(input, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create knowledge base"})
		return
	}
	c.JSON(http.StatusCreated, kb) // 返回 201 Created 状态码和创建的实体
}

// ListKnowledgeBases 是获取知识库列表的 Handler
func ListKnowledgeBases(c *gin.Context) {
	userIDFloat, _ := c.Get("user_id")
	userID := uint(userIDFloat.(float64))

	kbs, err := service.ListKnowledgeBases(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list knowledge bases"})
		return
	}

	c.JSON(http.StatusOK, kbs)
}

// UploadDocumentHandler 处理文档上传请求
func UploadDocumentHandler(c *gin.Context) {
	// 1. 从URL路径中获取knowledge_base_id
	kbIDStr := c.Param("kb_id")
	kbID, err := strconv.ParseUint(kbIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid knowledge base ID"})
		return
	}
	// 2. 检查知识库是否已经存在并且属于当前用户
	userID, _ := c.Get("user_id")
	var kb models.KnowledgeBase
	if err := core.DB.Preload("User").Where("id = ? AND user_id = ?", kbID, uint(userID.(float64))).First(&kb).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Knowledge base not found or you don't have permission"})
		return
	}
	// 3.从表单中获取上传的文件
	file, err := c.FormFile("document")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Document file is required"})
		return
	}
	// 4. 创建上传目录（如果不存在）
	uploadDir := filepath.Join("uploads", kbIDStr)
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	// 5. 保存文件到磁盘
	// (为了安全，可以对文件名做处理，防止路径遍历等问题，这里暂时简化)
	filePath := filepath.Join(uploadDir, filepath.Base(file.Filename))
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// 6. 在数据库中创建文档记录
	doc := models.Document{
		KnowledgeBaseID: uint(kbID),
		FileName:        file.Filename,
		FilePath:        filePath,
		FileSize:        file.Size,
		Status:          "UPLOADED",
	}
	if err := core.DB.Create(&doc).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create document record in db"})
		return
	}

	// 触发异步向量化任务...
	go service.ProcessAndEmbedDocument(doc)
	log.Printf("Dispatched background processing job for document ID: %d", doc.ID)

	c.JSON(http.StatusOK, gin.H{"message": "Document uploaded successfully", "document": doc})
}
