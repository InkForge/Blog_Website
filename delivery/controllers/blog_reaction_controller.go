package controllers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/InkForge/Blog_Website/domain"
	"github.com/gin-gonic/gin"
)

type BlogReactionController struct {
	BlogReactionUsecase domain.IBlogReactionUsecase
}

func NewBlogReactionController(usecase domain.IBlogReactionUsecase) *BlogReactionController {
	return &BlogReactionController{
		BlogReactionUsecase: usecase,
	}
}

func (bc *BlogReactionController) LikeBlog(c *gin.Context) {
	ogCtx := c.Request.Context()
	blogID := c.Param("id")

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: user ID is of wrong type"})
		return
	}

	ctx, cancel := context.WithTimeout(ogCtx, 5*time.Second)
	defer cancel()
	err := bc.BlogReactionUsecase.LikeBlog(ctx, blogID, userIDStr)
	if err != nil {
		switch {
		case err == context.DeadlineExceeded:
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Request timed out during like operation"})
		case errors.Is(err, domain.ErrBlogReactionNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "Blog reaction not found"})
		case errors.Is(err, domain.ErrCheckBlogReactionFailed):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check blog reaction", "details": err.Error()})
		case errors.Is(err, domain.ErrCreateBlogReactionFailed):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create blog reaction", "details": err.Error()})
		case errors.Is(err, domain.ErrUpdateBlogReactionFailed):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update blog reaction", "details": err.Error()})
		case errors.Is(err, domain.ErrIncrementLikeFailed):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to increment like count", "details": err.Error()})
		case errors.Is(err, domain.ErrToggleLikeDislikeFailed):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle like/dislike counts", "details": err.Error()})
		case errors.Is(err, domain.ErrInvalidBlogIdFormat):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID format"})
		case errors.Is(err, domain.ErrInsertingDocuments):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert document(s)", "details": err.Error()})
		case errors.Is(err, domain.ErrRetrievingDocuments):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve document(s)", "details": err.Error()})
		case errors.Is(err, domain.ErrDecodingDocument):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode document", "details": err.Error()})
		case errors.Is(err, domain.ErrDeletingBlogReaction):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete blog reaction", "details": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like blog", "details": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Blog liked successfully"})
}

func (bc *BlogReactionController) DislikeBlog(c *gin.Context) {
	ogCtx := c.Request.Context()
	blogID := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: user ID is of wrong type"})
		return
	}
	ctx, cancel := context.WithTimeout(ogCtx, 5*time.Second)
	defer cancel()
	err := bc.BlogReactionUsecase.DislikeBlog(ctx, blogID, userIDStr)
	if err != nil {
		switch {
		case err == context.DeadlineExceeded:
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Request timed out during dislike operation"})
		case errors.Is(err, domain.ErrBlogReactionNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "Blog reaction not found"})
		case errors.Is(err, domain.ErrCheckBlogReactionFailed):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check blog reaction", "details": err.Error()})
		case errors.Is(err, domain.ErrCreateBlogReactionFailed):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create blog reaction", "details": err.Error()})
		case errors.Is(err, domain.ErrUpdateBlogReactionFailed):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update blog reaction", "details": err.Error()})
		case errors.Is(err, domain.ErrIncrementLikeFailed):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to increment dislike count", "details": err.Error()})
		case errors.Is(err, domain.ErrToggleLikeDislikeFailed):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle like/dislike counts", "details": err.Error()})
		case errors.Is(err, domain.ErrInvalidBlogIdFormat):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID format"})
		case errors.Is(err, domain.ErrInsertingDocuments):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert document(s)", "details": err.Error()})
		case errors.Is(err, domain.ErrRetrievingDocuments):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve document(s)", "details": err.Error()})
		case errors.Is(err, domain.ErrDecodingDocument):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode document", "details": err.Error()})
		case errors.Is(err, domain.ErrDeletingBlogReaction):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete blog reaction", "details": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to dislike blog", "details": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Blog disliked successfully"})
}

func (bc *BlogReactionController) UnlikeBlog(c *gin.Context) {
	ogCtx := c.Request.Context()
	blogID := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: user ID is of wrong type"})
		return
	}
	ctx, cancel := context.WithTimeout(ogCtx, 5*time.Second)
	defer cancel()
	err := bc.BlogReactionUsecase.UnlikeBlog(ctx, blogID, userIDStr)
	if err != nil {
		switch {
		case err == context.DeadlineExceeded:
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Request timed out during unlike operation"})
		case errors.Is(err, domain.ErrBlogReactionNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "Blog reaction not found"})
		case errors.Is(err, domain.ErrCheckBlogReactionFailed):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check blog reaction", "details": err.Error()})
		case errors.Is(err, domain.ErrDeletingBlogReaction):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete blog reaction", "details": err.Error()})
		case errors.Is(err, domain.ErrToggleLikeDislikeFailed):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle like/dislike counts", "details": err.Error()})
		case errors.Is(err, domain.ErrInvalidBlogIdFormat):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID format"})
		case errors.Is(err, domain.ErrRetrievingDocuments):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve document(s)", "details": err.Error()})
		case errors.Is(err, domain.ErrDecodingDocument):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode document", "details": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unlike blog", "details": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Blog unliked successfully"})
}

func (bc *BlogReactionController) UndislikeBlog(c *gin.Context) {
	ogCtx := c.Request.Context()
	blogID := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: user ID is of wrong type"})
		return
	}
	ctx, cancel := context.WithTimeout(ogCtx, 5*time.Second)
	defer cancel()
	err := bc.BlogReactionUsecase.UndislikeBlog(ctx, blogID, userIDStr)
	if err != nil {
		switch {
		case err == context.DeadlineExceeded:
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Request timed out during undislike operation"})
		case errors.Is(err, domain.ErrBlogReactionNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "Blog reaction not found"})
		case errors.Is(err, domain.ErrCheckBlogReactionFailed):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check blog reaction", "details": err.Error()})
		case errors.Is(err, domain.ErrDeletingBlogReaction):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete blog reaction", "details": err.Error()})
		case errors.Is(err, domain.ErrToggleLikeDislikeFailed):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle like/dislike counts", "details": err.Error()})
		case errors.Is(err, domain.ErrInvalidBlogIdFormat):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID format"})
		case errors.Is(err, domain.ErrRetrievingDocuments):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve document(s)", "details": err.Error()})
		case errors.Is(err, domain.ErrDecodingDocument):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode document", "details": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to undislike blog", "details": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Blog undisliked successfully"})
}
