package client

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

type OpenAIClient struct {
    client *openai.Client
	aimodel string
}

func NewOpenAIClient(apiKey, aimodel string) (*OpenAIClient, error) {
    return &OpenAIClient{
        client: openai.NewClient(apiKey),
		aimodel: aimodel,
    }, nil
}

func (c *OpenAIClient) Generate(ctx context.Context, prompt string) (string, error) {
    resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
        Model: c.aimodel, 
        Messages: []openai.ChatCompletionMessage{
            {
                Role:    "system",
                Content: "You are a helpful AI assistant focused on blog content creation.",
            },
            {
                Role:    "user",
                Content: prompt,
            },
        },
        Temperature: 0.7,
    })
    if err != nil {
        return "", fmt.Errorf("OpenAI API error: %w", err)
    }

    if len(resp.Choices) == 0 {
        return "", fmt.Errorf("OpenAI API returned no choices")
    }

    return resp.Choices[0].Message.Content, nil
}