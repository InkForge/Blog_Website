package controllers

// imports
import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/InkForge/Blog_Website/domain"
	"github.com/InkForge/Blog_Website/usecases"
	"github.com/gin-gonic/gin"
)

// user register controller
type UserController struct {
	userUseCase usecases.UserUseCase       // user usecase for user operations 
}

func (uc *UserController) Register(c *gin.Context) {
	
	ctx := c.Request.Context()

	var input domain.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	// context timeout handling
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// create user via usecase
	user, err := uc.userUseCase.Register(&input, nil)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidEmailFormat):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		case errors.Is(err, domain.ErrEmailAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		case errors.Is(err, domain.ErrPasswordHashingFailed):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		case errors.Is(err, domain.ErrTokenGenerationFailed):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate verification token"})
		case errors.Is(err, domain.ErrUserCreationFailed):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User could not be created"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong", "details": err.Error()})
		}
		return
	}

	// prepare safe response - omit sensitive info
	response := gin.H{
		"message":  "User registered successfully",
		"user_id":  user.UserID,
		"email":    user.Email,
		"username": user.Username,
		"role":     user.Role,
		"provider": user.Provider,
	}

	c.JSON(http.StatusCreated, response)
}
