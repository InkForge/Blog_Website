package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/InkForge/Blog_Website/delivery/controllers/dto"
	"github.com/InkForge/Blog_Website/domain"
	"github.com/gin-gonic/gin"
)

type AIController struct {
	aiUsecase domain.IAIUseCase
}

// NewAIController creates a controller exposing AI helper endpoints.
func NewAIController(aiUsecase domain.IAIUseCase) *AIController {
	return &AIController{aiUsecase: aiUsecase}
}

// SuggestTags generates up to N tags from a given title/content.
func (ac *AIController) SuggestTags(c *gin.Context) {
	if ac.aiUsecase == nil {
		c.JSON(http.StatusServiceUnavailable, dto.ErrorResponse{Error: "ServiceUnavailable", Message: "AI service not configured", Code: http.StatusServiceUnavailable})
		return
	}
	var req dto.SuggestTagsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "InvalidPayload", Message: err.Error(), Code: http.StatusBadRequest})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 20*time.Second)
	defer cancel()
	tags, err := ac.aiUsecase.SuggestTags(ctx, req.Title, req.Content, req.MaxTags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "SuggestTagsFailed", Message: err.Error(), Code: http.StatusInternalServerError})
		return
	}
	c.JSON(http.StatusOK, dto.SuggestTagsResponse{Tags: tags})
}

// Summarize returns a concise summary of the provided content.
func (ac *AIController) Summarize(c *gin.Context) {
	if ac.aiUsecase == nil {
		c.JSON(http.StatusServiceUnavailable, dto.ErrorResponse{Error: "ServiceUnavailable", Message: "AI service not configured", Code: http.StatusServiceUnavailable})
		return
	}
	var req dto.SummarizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "InvalidPayload", Message: err.Error(), Code: http.StatusBadRequest})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 20*time.Second)
	defer cancel()
	summary, err := ac.aiUsecase.Summarize(ctx, req.Content, req.MaxWords)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "SummarizeFailed", Message: err.Error(), Code: http.StatusInternalServerError})
		return
	}
	c.JSON(http.StatusOK, dto.SummarizeResponse{Summary: summary})
}

// GenerateTitle produces a single title for the given content in the requested style.
func (ac *AIController) GenerateTitle(c *gin.Context) {
	if ac.aiUsecase == nil {
		c.JSON(http.StatusServiceUnavailable, dto.ErrorResponse{Error: "ServiceUnavailable", Message: "AI service not configured", Code: http.StatusServiceUnavailable})
		return
	}
	var req dto.GenerateTitleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "InvalidPayload", Message: err.Error(), Code: http.StatusBadRequest})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 20*time.Second)
	defer cancel()
	title, err := ac.aiUsecase.GenerateTitle(ctx, req.Content, req.Style)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "GenerateTitleFailed", Message: err.Error(), Code: http.StatusInternalServerError})
		return
	}
	c.JSON(http.StatusOK, dto.GenerateTitleResponse{Title: title})
}

// SuggestContent generates new content from user-provided keywords with optional style and length.
func (ac *AIController) SuggestContent(c *gin.Context) {
	if ac.aiUsecase == nil {
		c.JSON(http.StatusServiceUnavailable, dto.ErrorResponse{Error: "ServiceUnavailable", Message: "AI service not configured", Code: http.StatusServiceUnavailable})
		return
	}
	var req dto.SuggestContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "InvalidPayload", Message: err.Error(), Code: http.StatusBadRequest})
		return
	}
	if req.WordCount <= 0 {
		req.WordCount = 250
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 25*time.Second)
	defer cancel()
	content, err := ac.aiUsecase.SuggestContent(ctx, req.Keywords, req.Style, req.WordCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "SuggestContentFailed", Message: err.Error(), Code: http.StatusInternalServerError})
		return
	}
	c.JSON(http.StatusOK, dto.SuggestContentResponse{Content: content})
}

// ImproveContent returns an improved version of draft content plus targeted suggestions.
func (ac *AIController) ImproveContent(c *gin.Context) {
	if ac.aiUsecase == nil {
		c.JSON(http.StatusServiceUnavailable, dto.ErrorResponse{Error: "ServiceUnavailable", Message: "AI service not configured", Code: http.StatusServiceUnavailable})
		return
	}
	var req dto.ImproveContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "InvalidPayload", Message: err.Error(), Code: http.StatusBadRequest})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 25*time.Second)
	defer cancel()
	res, err := ac.aiUsecase.ImproveContent(ctx, req.Content, req.Focus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "ImproveContentFailed", Message: err.Error(), Code: http.StatusInternalServerError})
		return
	}
	c.JSON(http.StatusOK, dto.ImproveContentResponse{ImprovedContent: res.ImprovedContent, Suggestions: res.Suggestions})
}

// Chat performs a stateless chat turn using the provided message history and returns the next reply.
func (ac *AIController) Chat(c *gin.Context) {
	if ac.aiUsecase == nil {
		c.JSON(http.StatusServiceUnavailable, dto.ErrorResponse{Error: "ServiceUnavailable", Message: "AI service not configured", Code: http.StatusServiceUnavailable})
		return
	}
	var req dto.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "InvalidPayload", Message: err.Error(), Code: http.StatusBadRequest})
		return
	}
	// Map DTO to domain
	msgs := make([]domain.AIMessage, 0, len(req.Messages))
	for _, m := range req.Messages {
		msgs = append(msgs, domain.AIMessage{Role: m.Role, Content: m.Content})
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()
	out, err := ac.aiUsecase.Chat(ctx, msgs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "ChatFailed", Message: err.Error(), Code: http.StatusInternalServerError})
		return
	}
	c.JSON(http.StatusOK, dto.ChatResponse{Message: dto.ChatMessage{Role: out.Role, Content: out.Content}})
}
