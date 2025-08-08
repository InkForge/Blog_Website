package usecases

import (
	"context"
	"strings"
	"unicode/utf8"

	"github.com/InkForge/Blog_Website/domain"
)

// AIUseCase implements domain.IAiUseCase
type AIUseCase struct{
	AIService domain.IAIUseCase
}

// NewAIUseCase returns an instance of new AIUsecase
func NewAIUsecase(AIService domain.IAIContentService) domain.IAIUseCase {
	return &AIUseCase{
		AIService: AIService,
	}
}

// trimToMax is a helper function that follows utf8 encoding to trim and return
// the string if it is aboved the required maxChars
func trimToMax(input string, maxChars int) string {
	if utf8.RuneCountInString(input) <= maxChars {
		return input
	}
	return string([]rune(input)[:maxChars])
}

// SuggestTags makes sure the title, content and maxTags are in the appropriate size.
// It returns domain.ErrContentMissing if content length is zero.
func (aiu *AIUseCase) SuggestTags(ctx context.Context, title, content string, maxTags int) ([]string, error) {
	content = strings.TrimSpace(content)
	content = trimToMax(content, 500)
	title = strings.TrimSpace(title)
	title = trimToMax(title, 50)

	if content == "" {
		return nil, domain.ErrContentMissing
	}
	if maxTags <= 0 || maxTags > 5 {
		maxTags = 5
	}
	return aiu.AIService.SuggestTags(ctx, title, content, maxTags)
}

// Summarize makes sure the content length and maxWords are in the appropriate range.
// It returns domain.ErrContentMissing if the content length is zero.
func (aiu *AIUseCase) Summarize(ctx context.Context, content string, maxWords int) (string, error) {
	content = strings.TrimSpace(content)
	content = trimToMax(content, 500)

	if content == "" {
		return "", domain.ErrContentMissing
	}

	if maxWords <= 0 || maxWords > 200 {
		maxWords = 200
	}

	return aiu.AIService.Summarize(ctx, content, maxWords)
}

// GenerateTitle makes sure the style is appropriate and the content length is withing the required range.
// It returns domain.ErrContentMissing if content length is zeor.
func (aiu *AIUseCase) GenerateTitle(ctx context.Context, content, style string) (string, error) {
	content = strings.TrimSpace(content)
	content = trimToMax(content, 500)

	if content == "" {
		return "", domain.ErrContentMissing
	}

	style = strings.TrimSpace(style)
	if style == "" {
		style = "formal"
	}

	// limit the styles to avoid unnecessary length limit
	allowedStyles := map[string]bool{
		"formal":    true,
		"science":   true,
		"casual":    true,
		"clickbait": true,
		"poetic":    true,
		"seo":       true,
		"humorous":  true,
	}

	if !allowedStyles[style] {
		style = "formal"
	}

	return aiu.AIService.GenerateTitle(ctx, content, style)
}

// SuggestContent makes sure the style and wordCount appropriate and within the allowed style and count.
func (aiu *AIUseCase) SuggestContent(ctx context.Context, keywords, style string, wordCount int) (string, error) {
	if wordCount <= 0 || wordCount > 500 {
		wordCount = 500
	}
	keywords = strings.TrimSpace(keywords)
	if keywords == "" {
		keywords = "interesting"
	}
	style = strings.TrimSpace(style)
	if style == "" {
		style = "formal"
	}

	// limit the styles to avoid unnecessary length limit
	allowedStyles := map[string]bool{
		"formal":    true,
		"science":   true,
		"casual":    true,
		"clickbait": true,
		"poetic":    true,
		"seo":       true,
		"humorous":  true,
	}

	if !allowedStyles[style] {
		style = "formal"
	}

	return aiu.AIService.SuggestContent(ctx, keywords, style, wordCount)
}

// ImproveContent makes sure the content length is not off limit and set the focus to common areas.
// It returns domain.ErrContentMissing if the content length is zero.
func (aiu *AIUseCase) ImproveContent(ctx context.Context, content, focus string) (domain.ImprovementResult, error) {
	content = strings.TrimSpace(content)
	content = trimToMax(content, 500)

	if content == "" {
		return domain.ImprovementResult{}, domain.ErrContentMissing
	}

	focus = strings.TrimSpace(focus)
	if focus == "" {
		focus = "grammar"
	}
	var commonFocusAreas = map[string]bool{
		"clarity":        true,
		"tone":           true,
		"grammar":        true,
		"conciseness":    true,
		"persuasiveness": true,
		"readability":    true,
		"engagement":     true,
		"style":          true,
		"flow":           true,
		"spelling":       true,
	}
	if !commonFocusAreas[focus] {
		focus = "grammar"
	}
	return aiu.AIService.ImproveContent(ctx, content, focus)
}

// Chat makes sure the AI messages history are not off limit and contains a valid message.
// Top of the slice is considered the most recent message between the client and AI
func (aiu *AIUseCase) Chat(ctx context.Context, messages []domain.AIMessage) (domain.AIMessage, error) {
	if len(messages) > 5 {
		messages = messages[len(messages)-5:]
	}
	for idx, message := range messages {
		messages[idx].Content = trimToMax(message.Content, 200)
	}
	return aiu.AIService.Chat(ctx, messages)
}

