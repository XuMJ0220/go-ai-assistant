package service

import (
	"context"
	"fmt"
	"go-ai-assistant/core"

	"github.com/sashabaranov/go-openai"
)

// SimpleChat 定义了简单聊天的输入
type SimpleChatInput struct {
	Prompt string `json:"prompt" binding:"required"`
}

func SimpleChat(input SimpleChatInput) (string, error) {
	resp, err := core.OpenAIClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: "qwen-plus",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: input.Prompt,
				},
			},
		},
	)
	if err != nil {
		// 可以在这里处理更详细的错误，比如 API Key 无效等
		fmt.Printf("OpenAI API Error: %v\n", err)
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}
