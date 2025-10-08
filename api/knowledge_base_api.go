package api

import (
	"go-ai-assistant/service"
	"net/http"

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