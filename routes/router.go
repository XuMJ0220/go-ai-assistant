package routes

import (
	"go-ai-assistant/api"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/ping", api.Ping)

	// 创建一个 API v1 路由组
	apiV1 := router.Group("/api/v1")
	{
		// 用户路由组
		userRouters := apiV1.Group("/users")
		{
			userRouters.POST("/register", api.Register)
		}
	}

	return router
}
