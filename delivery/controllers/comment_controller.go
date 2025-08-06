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

func NewCommentController(usecase domain.ICommentUsecase) *CommentController {
	return &CommentController{commentUseCase: usecase}
}

func (cc *CommentController) AddComment(c *gin.Context) {
	blogID := c.Param("blogID")
	userID := c.GetString("userID")
	role := c.GetString("userRole")

	var req dto.CommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload", "details": err.Error()})
		return
	}

	comment := &domain.Comment{
		User_id: userID,
		Content: req.Content,
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	commentID, err := cc.commentUseCase.AddComment(ctx, blogID, comment, role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Comment added successfully",
		"comment_id": commentID,
	})
}

func (cc *CommentController) UpdateComment(c *gin.Context) {
	commentID := c.Param("commentID")
	userID := c.GetString("userID")
	role := c.GetString("userRole")

	var req dto.CommentUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	comment := &domain.Comment{
		Comment_id: commentID,
		User_id:    userID,
		Content:    req.Content,
	}

	err := cc.commentUseCase.UpdateComment(ctx, commentID, comment, role)
	if err != nil {
		if errors.Is(err, domain.ErrCommentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		} else if errors.Is(err, domain.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to update this comment"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment updated successfully"})
}

func (cc *CommentController) RemoveComment(c *gin.Context) {
	blogID := c.Param("blogID")
	commentID := c.Param("commentID")
	userID := c.GetString("userID")
	role := c.GetString("userRole")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err := cc.commentUseCase.RemoveComment(ctx, blogID, commentID, userID, role)
	if err != nil {
		if errors.Is(err, domain.ErrCommentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		} else if errors.Is(err, domain.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to delete this comment"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment removed successfully"})
}

func (cc *CommentController) GetBlogComments(c *gin.Context) {
	blogID := c.Param("blogID")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	comments, err := cc.commentUseCase.GetBlogComments(ctx, blogID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.FromDomainComments(comments)
	c.JSON(http.StatusOK, response)
}
