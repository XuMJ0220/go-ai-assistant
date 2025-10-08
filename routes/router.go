package routes

import (
	"go-ai-assistant/api"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router:=gin.Default()

	router.GET("/ping",api.Ping)

	return router
}
