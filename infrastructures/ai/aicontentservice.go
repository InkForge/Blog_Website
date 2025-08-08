package infrastructures

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/InkForge/Blog_Website/delivery/controllers/dto"
	"github.com/InkForge/Blog_Website/domain"
)

type AIContentService struct {
	client domain.IAIModelClient 
}

func NewAIContentService(client domain.IAIModelClient) domain.IAIContentService {
	return &AIContentService{client: client}
}

// SuggestTags generates relevant tags for an article
func (s *AIContentService) SuggestTags(ctx context.Context, title, content string, maxTags int) ([]string, error) {
	if maxTags <= 0 || maxTags > 10 {
		maxTags = 10
	}

	prompt := fmt.Sprintf(`
		You are an AI tag generator for blog articles. 

		Given the title and content below, respond with ONLY a valid JSON object in this exact format:
		{"tags": ["tag1", "tag2", "tag3"]}

		Rules:
		- Maximum %d tags
		- Each tag should be 1–3 words
		- Tags must be factual, relevant to the content, and non-offensive
		- No personal identifiers or sensitive data
		- No political, sexual, or discriminatory terms unless explicitly part of the article content
		- Do not include backticks, code fences, or explanations

		Title: %q
		Content: %q
	`, maxTags, title, content)

	rawResp, err := s.client.Generate(ctx, prompt)
	if err != nil {
		return nil, err
	}

	var res dto.SuggestTagsResponse
	if err := json.Unmarshal([]byte(rawResp), &res); err != nil {
		return nil, err
	}

	return res.Tags, nil
}

// Summarize condenses content into a short summary
func (s *AIContentService) Summarize(ctx context.Context, content string, maxWords int) (string, error) {
	if maxWords <= 0 || maxWords > 500 {
		maxWords = 200
	}

	prompt := fmt.Sprintf(`
		You are an AI summarizer.

		Summarize the given content into at most %d words.
		Keep meaning intact, avoid bias, and use a neutral tone.

		Respond with ONLY a valid JSON object in this exact format:
		{"summary": "your summary here"}

		Rules:
		- Maximum %d words
		- Summary must be faithful to the original meaning
		- Avoid personal opinions, exaggerations, or speculative claims
		- Do not add facts not in the source
		- Do not include backticks, code fences, or any text outside the JSON

		Content: %q
	`, maxWords, maxWords, content)

	rawResp, err := s.client.Generate(ctx, prompt)
	if err != nil {
		return "", err
	}

	var res dto.SummarizeResponse
	if err := json.Unmarshal([]byte(rawResp), &res); err != nil {
		return "", err
	}

	return res.Summary, nil
}

// GenerateTitle creates a title in the requested style
func (s *AIContentService) GenerateTitle(ctx context.Context, content, style string) (string, error) {
	prompt := fmt.Sprintf(`
		You are an AI title generator. 
		Generate a single title in the "%s" style for the content provided.
		Return ONLY a JSON object: {"title": "..."}.

		Content: "%s"

		Rules:
		- Must be under 80 characters.
		- Avoid clickbait or misleading phrasing.
		- No inappropriate language.
		- Must match style requested.
	`, style, content)

	rawResp, err := s.client.Generate(ctx, prompt)
	if err != nil {
		return "", err
	}

	var res dto.GenerateTitleResponse
	if err := json.Unmarshal([]byte(rawResp), &res); err != nil {
		return "", err
	}

	return res.Title, nil
}

// SuggestContent drafts an article from keywords and style
func (s *AIContentService) SuggestContent(ctx context.Context, keywords, style string, wordCount int) (string, error) {
	if wordCount <= 0 {
		wordCount = 250
	}

	prompt := fmt.Sprintf(`
		You are an AI content writer.
		Using the keywords "%s" and style "%s", write a short article of about %d words.
		Keep it factual, unbiased, and relevant.
		Return ONLY the article content as plain text, without JSON or any other formatting.

		Do not include any extra commentary or markdown formatting.

		Rules:
		- Stay strictly on topic of keywords.
		- Avoid sensitive, harmful, or adult content unless explicitly requested for legitimate purposes.
		- No promoting illegal activities.
	`, keywords, style, wordCount)


	rawResp, err := s.client.Generate(ctx, prompt)
	if err != nil {
		return "", err
	}

	return rawResp, nil
}

// ImproveContent enhances an article based on focus
func (s *AIContentService) ImproveContent(ctx context.Context, content, focus string) (domain.ImprovementResult, error) {
	prompt := fmt.Sprintf(`
		You are an AI content improver.
		Given the content below, improve it focusing on "%s".
		Return ONLY a JSON object:
		{
		"improved_content": "...",
		"suggestions": ["...", "..."]
		}

		Content: "%s"

		Rules:
		- Preserve the original meaning.
		- Suggestions must be constructive and safe.
		- No introducing new opinions or factual errors.
		- Avoid changing writing style to offensive or political.
	`, focus, content)

	rawResp, err := s.client.Generate(ctx, prompt)
	if err != nil {
		return domain.ImprovementResult{}, err
	}

	var res dto.ImproveContentResponse
	if err := json.Unmarshal([]byte(rawResp), &res); err != nil {
		return domain.ImprovementResult{}, err
	}

	// convert DTO to domain model before returning
	return domain.ImprovementResult{
		ImprovedContent: res.ImprovedContent,
		Suggestions:     res.Suggestions,
	}, nil
}


// Chat provides a controlled AI conversation
func (s *AIContentService) Chat(ctx context.Context, messages []domain.AIMessage) (domain.AIMessage, error) {
    // convert domain.AIMessage slice to dto.ChatMessage slice
    dtoMessages := make([]dto.ChatMessage, len(messages))
    for i, m := range messages {
        dtoMessages[i] = dto.ChatMessage{
            Role:    m.Role,
            Content: m.Content,
        }
    }

    messagesJSON, _ := json.Marshal(dtoMessages)

    prompt := fmt.Sprintf(`
        You are an AI assistant for content creation.
        Respond to the user's messages based strictly on the conversation context provided.
        Do not answer outside the topic of content creation, writing, or AI assistance.
        Return ONLY a JSON object: 
        {"message": {"role": "assistant", "content": "..."}}

        Messages: %s

        Rules:
        - Only answer within allowed domain (content creation, writing help).
        - Refuse to answer political, harmful, or personal-identifying requests.
    `, string(messagesJSON))

    rawResp, err := s.client.Generate(ctx, prompt)
    if err != nil {
        return domain.AIMessage{}, err
    }

    var res dto.ChatResponse
    if err := json.Unmarshal([]byte(rawResp), &res); err != nil {
        return domain.AIMessage{}, err
    }

    return domain.AIMessage{
        Role:    res.Message.Role,
        Content: res.Message.Content,
    }, nil
}
