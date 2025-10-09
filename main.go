package main

import (
	"fmt"
	"go-ai-assistant/config"
	"go-ai-assistant/core"
	"go-ai-assistant/routes"
	"os"
)

func main() {
	wd, _ := os.Getwd()
	fmt.Println("Current Working Directory:", wd)
	config.LoadConfig()

	core.InitDB()

	core.InitLLMClient()

	router := routes.SetupRouter()

	router.Run(":8080")
}
