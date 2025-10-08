package api

import (
	"go-ai-assistant/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SimpleChatHandler(c *gin.Context) {
	var input service.SimpleChatInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 从Context中获取用户信息(由 AuthMiddleware设置)
	// 这样我们就知道是那个用户在提问了
	//userID, _ := c.Get("user_id")

	// 调用AI服务
	response, err := service.SimpleChat(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 可以在这里记录是哪个用户(userID)发起了请求和得到了什么回答
	// 这是后续做聊天记录功能的基础

	c.JSON(http.StatusOK, gin.H{
		"response": response,
	})
}
