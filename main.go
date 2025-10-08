package main

import (
	"go-ai-assistant/config"
	"go-ai-assistant/core"
	"go-ai-assistant/routes"
)

func main() {
	config.LoadConfig()

	core.InitDB()

	router := routes.SetupRouter()

	router.Run(":8080")
}
