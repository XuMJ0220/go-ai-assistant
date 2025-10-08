package core

import (
	"go-ai-assistant/config"

	"github.com/sashabaranov/go-openai"
)

// OpenAIClient 全局的OpenAI客户端实例
var OpenAIClient *openai.Client

func InitLLMClient() {
	cfg := openai.DefaultConfig(config.AppConfig.DashScope.ApiKey)
	cfg.BaseURL = "https://dashscope.aliyuncs.com/compatible-mode/v1"

	OpenAIClient = openai.NewClientWithConfig(cfg)
}
