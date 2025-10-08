package service

import (
	"fmt"
	"go-ai-assistant/core"
	"go-ai-assistant/models"

	"github.com/google/uuid"
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
