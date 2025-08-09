package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/InkForge/Blog_Website/delivery/controllers/dto"
)

type DeepSeekClient struct {
	apiKey   string
	baseURL  string
	model    string
	client   *http.Client
}

func NewDeepSeekClient(apiKey, baseURL, model string) *DeepSeekClient {
	if baseURL == "" {
		baseURL = "https://api.deepseek.com/v1"
	}
	return &DeepSeekClient{
		apiKey:  apiKey,
		baseURL: baseURL,
		model:   model,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (d *DeepSeekClient) Generate(ctx context.Context, prompt string) (string, error) {
	// wrap prompt into the chat messages expected by DeepSeek (user role)
	messages := []dto.ChatMessage{
		{Role: "user", Content: prompt},
	}

	reqBody := dto.ChatCompletionRequest{
		Model:    d.model,
		Messages: messages,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", d.baseURL+"/chat/completions", bytes.NewReader(jsonData))
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+d.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", errors.New(fmt.Sprintf("deepseek API error: status %d, body: %s", resp.StatusCode, string(bodyBytes)))
	}

	var res dto.ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if len(res.Choices) == 0 {
		return "", errors.New("deepseek API error: no choices returned")
	}

	return res.Choices[0].Message.Content, nil
}
