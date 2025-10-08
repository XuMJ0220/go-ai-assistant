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
		userRoutes := apiV1.Group("/users")
		{
			userRoutes.POST("/register", api.Register)
			userRoutes.POST("/login", api.Login) // 添加登录路由
		}

		// 受保护的路由组
		protectedRoutes:=apiV1.Group("")
		protectedRoutes.Use(api.AuthMiddleware())
		{
			// 知识库路由组
			kbRoutes := protectedRoutes.Group("/knowledge-bases")
			{
				kbRoutes.POST("", api.CreateKnowledgeBase)    // POST /api/v1/knowledge-bases
				kbRoutes.GET("", api.ListKnowledgeBases)      // GET /api/v1/knowledge-bases
			}
		}
	}

	return router
}
