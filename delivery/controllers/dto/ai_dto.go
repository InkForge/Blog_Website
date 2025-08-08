package dto

type SuggestTagsRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	MaxTags int    `json:"max_tags"`
}

type SummarizeRequest struct {
	Content  string `json:"content"`
	MaxWords int    `json:"max_words"`
}

type GenerateTitleRequest struct {
	Content string `json:"content"`
	Style   string `json:"style"`
}

// Responses
type SuggestTagsResponse struct {
	Tags []string `json:"tags"`
}

type SummarizeResponse struct {
	Summary string `json:"summary"`
}

type GenerateTitleResponse struct {
	Title string `json:"title"`
}

// New AI features per PRD
type SuggestContentRequest struct {
	Keywords  string `json:"keywords"`
	Style     string `json:"style"`      // e.g., "technical", "casual"
	WordCount int    `json:"word_count"` // optional, default 250
}

type SuggestContentResponse struct {
	Content string `json:"content"`
}

type ImproveContentRequest struct {
	Content string `json:"content"`
	Focus   string `json:"focus"` // e.g., "clarity", "tone", "grammar"
}

type ImproveContentResponse struct {
	ImprovedContent string   `json:"improved_content"`
	Suggestions     []string `json:"suggestions"`
}

type ChatMessage struct {
	Role    string `json:"role"` // system | user | assistant
	Content string `json:"content"`
}

type ChatRequest struct {
	Messages []ChatMessage `json:"messages"`
}

type ChatResponse struct {
	Message ChatMessage `json:"message"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}