package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/InkForge/Blog_Website/delivery/controllers/dto"
	"github.com/InkForge/Blog_Website/domain"
	"github.com/gin-gonic/gin"
)

// CommentReactionController handles HTTP requests for comment reaction operations
type CommentReactionController struct {
	reactionUseCase domain.ICommentReactionUsecase
}

// NewCommentReactionController creates a new comment reaction controller
func NewCommentReactionController(reactionUseCase domain.ICommentReactionUsecase) *CommentReactionController {
	return &CommentReactionController{
		reactionUseCase: reactionUseCase,
	}
}

// ReactToComment handles POST /comments/:commentID/react/:status
func (crc *CommentReactionController) ReactToComment(c *gin.Context) {
	commentID := c.Param("commentID")
	if commentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Comment ID is required"})
		return
	}

	statusStr := c.Param("status")
	status, err := strconv.Atoi(statusStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status parameter. Must be -1, 0, or 1"})
		return
	}

	// Validate status values
	if status != -1 && status != 0 && status != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status must be -1 (dislike), 0 (neutral), or 1 (like)"})
		return
	}

	// Get user ID from authenticated context (set by auth middleware)
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// React to comment via use case
	var reactionErr error
	switch status {
	case 1:
		reactionErr = crc.reactionUseCase.LikeComment(c.Request.Context(), commentID, userID)
	case -1:
		reactionErr = crc.reactionUseCase.DislikeComment(c.Request.Context(), commentID, userID)
	case 0:
		reactionErr = crc.reactionUseCase.RemoveReaction(c.Request.Context(), commentID, userID)
	}

	if reactionErr != nil {
		switch {
		case errors.Is(reactionErr, domain.ErrInvalidCommentID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		case errors.Is(reactionErr, domain.ErrInvalidUserID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		case errors.Is(reactionErr, domain.ErrCommentNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to react to comment", "details": reactionErr.Error()})
		}
		return
	}

	// Create appropriate success message
	var message string
	switch status {
	case 1:
		message = "Comment liked successfully"
	case -1:
		message = "Comment disliked successfully"
	case 0:
		message = "Reaction removed successfully"
	}

	response := dto.CreateReactionResponse(commentID, userID, status, message)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// GetUserReaction handles GET /comments/:commentID/reaction
func (crc *CommentReactionController) GetUserReaction(c *gin.Context) {
	commentID := c.Param("commentID")
	if commentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Comment ID is required"})
		return
	}

	// Get user ID from authenticated context
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get user reaction via use case
	action, err := crc.reactionUseCase.GetUserReaction(c.Request.Context(), commentID, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidCommentID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		case errors.Is(err, domain.ErrInvalidUserID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		case errors.Is(err, domain.ErrCommentNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user reaction", "details": err.Error()})
		}
		return
	}

	response := dto.CreateUserReactionResponse(commentID, userID, action)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
} 