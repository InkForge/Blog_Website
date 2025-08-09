package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GroqClient struct {
	apiKey string
	model  string
}

func NewGroqClient(apiKey, model string) *GroqClient {
	return &GroqClient{
		apiKey: apiKey,
		model:  model,
	}
}

func (c *GroqClient) Generate(ctx context.Context, prompt string) (string, error) {
	url := "https://api.groq.com/openai/v1/chat/completions"

	reqBody := map[string]interface{}{
		"model": c.model,
		"messages": []map[string]string{
			{"role": "system", "content": "You are an AI assistant."},
			{"role": "user", "content": prompt},
		},
		"temperature": 0.7,
	}

	data, _ := json.Marshal(reqBody)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("groq API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no completion returned from groq")
	}

	return result.Choices[0].Message.Content, nil
}
