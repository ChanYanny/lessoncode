package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

type LLMClient struct {
	client *openai.Client
}

func main() {
	// question := "Eino ADK 中的 DeepAgent 是什么？一句话总结"
	llmClient, err := NewLLMClient()
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx := context.Background()
	// res, totalTokens, err := llmClient.Ask(ctx, question)
	_ = llmClient.Chat(ctx)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println("问题：", question, "\n回答：", res)
	// fmt.Println("消耗的总 Token 数：", totalTokens)
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

func (c *LLMClient) Ask(ctx context.Context, question string) (string, int, error) {
	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: os.Getenv("MODEL_ID"),
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "你是一个精通 Eino 框架的专家",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: question,
				},
			},
			Temperature:         0.3,
			MaxCompletionTokens: 1024,
		},
	)
	if err != nil {
		return "", 0, err
	}
	return resp.Choices[0].Message.Content, resp.Usage.TotalTokens, nil
}

func (c *LLMClient) Chat(ctx context.Context) error {
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "你是一个精通 Eino 框架的专家",
		},
	}
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("开始对话（输入 quit 退出）：")
	for {
		fmt.Print("\033[36m你 >> \033[0m")
		if !scanner.Scan() {
			return nil
		}
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		if input == "quit" {
			fmt.Println("再见👋")
			return nil
		}
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: input,
		})
		resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
			Model:    os.Getenv("MODEL_ID"),
			Messages: messages,
		})
		if err != nil {
			fmt.Println(err)
			messages = messages[:len(messages)-1]
			continue
		}
		reply := resp.Choices[0].Message.Content
		fmt.Println("\033[36m大模型 >> \033[0m", reply)
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: reply,
		})
	}
}
