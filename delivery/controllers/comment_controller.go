package controllers

import (
	"errors"
	"net/http"

	"github.com/InkForge/Blog_Website/delivery/controllers/dto"
	"github.com/InkForge/Blog_Website/domain"
	"github.com/gin-gonic/gin"
)

// CommentController handles HTTP requests for comment operations
type CommentController struct {
	commentUseCase domain.ICommentUsecase
}

// NewCommentController creates a new comment controller
func NewCommentController(commentUseCase domain.ICommentUsecase) *CommentController {
	return &CommentController{
		commentUseCase: commentUseCase,
	}
}



// AddComment handles POST /blogs/:blogID/comments
func (cc *CommentController) AddComment(c *gin.Context) {
	blogID := c.Param("blogID")
	if blogID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Blog ID is required"})
		return
	}

	// Get user ID from authenticated context (set by auth middleware)
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req dto.CommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	// Create comment domain object
	comment := &domain.Comment{
		Blog_id: blogID,
		User_id: userID,
		Content: req.Content,
	}

	// Add comment via use case
	commentID, err := cc.commentUseCase.AddComment(c.Request.Context(), blogID, comment)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrCommentRequired):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Comment is required"})
		case errors.Is(err, domain.ErrInvalidBlogID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID"})
		case errors.Is(err, domain.ErrEmptyCommentContent):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Comment content cannot be empty"})
		case errors.Is(err, domain.ErrInvalidUserID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		case errors.Is(err, domain.ErrBlogNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add comment", "details": err.Error()})
		}
		return
	}

	// Prepare response
	response := gin.H{
		"message":     "Comment added successfully",
		"comment_id":  commentID,
		"blog_id":     blogID,
		"user_id":     userID,
		"content":     req.Content,
	}

	c.JSON(http.StatusCreated, response)
}

// GetBlogComments handles GET /blogs/:blogID/comments
func (cc *CommentController) GetBlogComments(c *gin.Context) {
	blogID := c.Param("blogID")
	if blogID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Blog ID is required"})
		return
	}

	// Get comments via use case
	comments, err := cc.commentUseCase.GetBlogComments(c.Request.Context(), blogID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrBlogNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments", "details": err.Error()})
		}
		return
	}

	// Convert domain comments to response format using DTO
	response := dto.FromDomainComments(comments)
	c.JSON(http.StatusOK, response)
}

// UpdateComment handles PUT /comments/:commentID
func (cc *CommentController) UpdateComment(c *gin.Context) {
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

	var req dto.CommentUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	// Create comment domain object for update
	comment := &domain.Comment{
		Comment_id: commentID,
		User_id:    userID,
		Content:    req.Content,
	}

	// Update comment via use case
	err := cc.commentUseCase.UpdateComment(c.Request.Context(), commentID, comment)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrCommentRequired):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Comment is required"})
		case errors.Is(err, domain.ErrCommentNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Comment updated successfully",
		"comment_id": commentID,
	})
}

// RemoveComment handles DELETE /blogs/:blogID/comments/:commentID
func (cc *CommentController) RemoveComment(c *gin.Context) {
	blogID := c.Param("blogID")
	commentID := c.Param("commentID")
	
	if blogID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Blog ID is required"})
		return
	}
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

	// Remove comment via use case
	err := cc.commentUseCase.RemoveComment(c.Request.Context(), blogID, commentID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidBlogID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID"})
		case errors.Is(err, domain.ErrInvalidCommentID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		case errors.Is(err, domain.ErrCommentNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove comment", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Comment removed successfully",
		"comment_id": commentID,
	})
} 