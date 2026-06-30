package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
)

type LLMStreamClient struct {
	client *openai.Client
}

func main() {
	ctx := context.Background()
	question := "Eino ADK 中的 DeepAgent 是什么？一句话总结"
	llmStreamClient, err := NewLLMStreamClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = llmStreamClient.StreamAsk(ctx, question)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func NewLLMStreamClient() (*LLMStreamClient, error) {
	apiKey := strings.TrimSpace(os.Getenv("DEEPSEEK_API_KEY"))
	if apiKey == "" {
		return nil, errors.New("apikey is required")
	}
	cfg := openai.DefaultConfig(apiKey)
	cfg.BaseURL = os.Getenv("BASE_URL")

	return &LLMStreamClient{
		client: openai.NewClientWithConfig(cfg),
	}, nil
}

func (c *LLMStreamClient) StreamAsk(ctx context.Context, question string) error {
	streamResp, err := c.client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
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
		Stream: true,
	})
	if err != nil {
		return err
	}
	defer streamResp.Close()
	fmt.Println("=== 大模型回复 ===")
	var (
		full         strings.Builder
		firstTokenAt time.Time
		startAt      = time.Now()
	)
	for {
		chunk, err := streamResp.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("streamResp.Recv: %w", err)
		}
		if len(chunk.Choices) == 0 {
			continue
		}
		content := chunk.Choices[0].Delta.Content
		if content == "" {
			continue
		}
		if firstTokenAt.IsZero() {
			firstTokenAt = time.Now()
		}
		fmt.Println(content)
		full.WriteString(content)
	}
	if !firstTokenAt.IsZero() {
		fmt.Printf("首 token 输出耗时：%s\n", firstTokenAt.Sub(startAt).Round(time.Millisecond))
	}
	fmt.Printf("总耗时：%s\n", time.Since(startAt).Round(time.Millisecond))
	fmt.Println(full.String())
	return nil
}
