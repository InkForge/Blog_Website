package controllers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/InkForge/Blog_Website/delivery/controllers/dto"
	"github.com/InkForge/Blog_Website/domain"
	"github.com/gin-gonic/gin"
)

type CommentController struct {
	commentUseCase domain.ICommentUsecase
}

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

	comment := &domain.Comment{
		Blog_id: blogID,
		User_id: userID,
		Content: req.Content,
	}

	ogCtx := c.Request.Context()
	ctx, cancel := context.WithTimeout(ogCtx, 5*time.Second)
	defer cancel()

	commentID, err := cc.commentUseCase.AddComment(ctx, blogID, comment)
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

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Comment added successfully",
		"comment_id": commentID,
		"blog_id":    blogID,
		"user_id":    userID,
		"content":    req.Content,
	})
}

// GetBlogComments handles GET /blogs/:blogID/comments
func (cc *CommentController) GetBlogComments(c *gin.Context) {
	blogID := c.Param("blogID")
	if blogID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Blog ID is required"})
		return
	}

	ogCtx := c.Request.Context()
	ctx, cancel := context.WithTimeout(ogCtx, 5*time.Second)
	defer cancel()

	comments, err := cc.commentUseCase.GetBlogComments(ctx, blogID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrBlogNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments", "details": err.Error()})
		}
		return
	}

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

	comment := &domain.Comment{
		Comment_id: commentID,
		User_id:    userID,
		Content:    req.Content,
	}

	ogCtx := c.Request.Context()
	ctx, cancel := context.WithTimeout(ogCtx, 5*time.Second)
	defer cancel()

	err := cc.commentUseCase.UpdateComment(ctx, commentID, comment)
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
		"success":    true,
		"message":    "Comment updated successfully",
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

	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	ogCtx := c.Request.Context()
	ctx, cancel := context.WithTimeout(ogCtx, 5*time.Second)
	defer cancel()

	err := cc.commentUseCase.RemoveComment(ctx, blogID, commentID)
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
		"success":    true,
		"message":    "Comment removed successfully",
		"comment_id": commentID,
	})
}
