package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

// 弱 Prompt：模型只知道"你是个代码审查专家"，其它全靠自由发挥
const weakPrompt = `你是一个代码审查专家，请回答用户问题。`

// 结构化 Prompt：把角色、目标、回答要求都写清楚
const structuredPrompt = `你是一个面向本机代码仓库的分析助手。    
    
回答要求：  1. 先给结论，再展开细节。  2. 涉及代码的地方贴关键片段，并说明这段代码在做什么。  3. 看不准的地方明确写"看不准"，附上需要再确认的文件路径。  4. 报告结构固定：概述 / 技术栈 / 目录分层 / 调用链路 / 业务清单 / 改进建议。`

type LLMClient struct {
	client *openai.Client
}

func main() {
	llmClient, err := NewLLMClient()
	if err != nil {
		fmt.Println(err)
		return
	}
	question := `我有一个 Go 项目放在 [xxxx]，希望你帮我梳理：  1. 这是一个什么类型的项目，用了哪些主要技术栈？  2. 一次 HTTP 请求从入口到数据库大概经过哪几层？  3. 我现在要接手维护，建议从哪个目录、哪个文件开始读？`
	fmt.Println("=== 弱 Prompt ===")
	printAnswer(llmClient, weakPrompt, question)
	fmt.Println("=== 结构化 Prompt ===")
	printAnswer(llmClient, structuredPrompt, question)

}

func NewLLMClient() (*LLMClient, error) {
	apiKey := strings.TrimSpace(os.Getenv("DEEPSEEK_API_KEY"))
	if apiKey == "" {
		return nil, errors.New("apikey is required")
	}
	cfg := openai.DefaultConfig(apiKey)
	cfg.BaseURL = os.Getenv("BASE_URL")

	return &LLMClient{
		client: openai.NewClientWithConfig(cfg),
	}, nil
}

func printAnswer(c *LLMClient, systemPrompt, question string) {
	resp, err := c.client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: os.Getenv("MODEL_ID"),
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemPrompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: question,
			},
		},
		Temperature:         0.3,
		MaxCompletionTokens: 1024,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.Choices[0].Message.Content)
}
