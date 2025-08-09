package domain

import "context"

// AIContentService defines AI-powered content assistance capabilities.
type IAIContentService interface {
    SuggestTags(ctx context.Context, title, content string, maxTags int) ([]string, error)
    Summarize(ctx context.Context, content string, maxWords int) (string, error)
    GenerateTitle(ctx context.Context, content, style string) (string, error)
    SuggestContent(ctx context.Context, keywords, style string, wordCount int) (string, error)
    ImproveContent(ctx context.Context, content, focus string) (ImprovementResult, error)
    Chat(ctx context.Context, messages []AIMessage) (AIMessage, error)
}

// IAIUseCase exposes AI helpers to the delivery layer.
type IAIUseCase interface {
    SuggestTags(ctx context.Context, title, content string, maxTags int) ([]string, error)
    Summarize(ctx context.Context, content string, maxWords int) (string, error)
    GenerateTitle(ctx context.Context, content, style string) (string, error)
    SuggestContent(ctx context.Context, keywords, style string, wordCount int) (string, error)
    ImproveContent(ctx context.Context, content, focus string) (ImprovementResult, error)
    Chat(ctx context.Context, messages []AIMessage) (AIMessage, error)
}

// AIMessage models a chat message for a simple AI chat feature.
type AIMessage struct {
    Role    string // "system" | "user" | "assistant"
    Content string
}

// ImprovementResult represents improved content plus suggestions.
type ImprovementResult struct {
    ImprovedContent string
    Suggestions     []string
}


// AIModelClient defines the low-level AI API communication methods.
// Generate sends a prompt string to the AI provider and returns the raw text response.
type IAIModelClient interface {
    Generate(ctx context.Context, prompt string) (string, error)
}
