package controllers

// imports
import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/InkForge/Blog_Website/domain"
	"github.com/gin-gonic/gin"
)

// user register controller
type UserController struct {
	userUseCase domain.IUserUseCase // user usecase for user operations
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

	user, err := uc.userUseCase.Register(ctx, &input, nil)

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

// user login controller
func (uc *UserController) Login(c *gin.Context) {

	ctx := c.Request.Context()

	// use separate struct for login input
	var loginInput struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	// bind and validate input
	if err := c.ShouldBindJSON(&loginInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// create domain user for login use case
	input := domain.User{
		Email:    loginInput.Email,
		Password: &loginInput.Password,
	}

	// perform login
	accessToken, refreshToken, user, err := uc.userUseCase.Login(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidEmailFormat):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		case errors.Is(err, domain.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		case errors.Is(err, domain.ErrOAuthUserCannotLoginWithPassword):
			c.JSON(http.StatusBadRequest, gin.H{"error": "This account uses OAuth login only"})
		case errors.Is(err, domain.ErrEmailNotVerified):
			c.JSON(http.StatusForbidden, gin.H{"error": "Please verify your email address first"})
		case errors.Is(err, domain.ErrTokenGenerationFailed):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Login failed",
				"details": err.Error(),
			})
		}
		return
	}

	// prepare sanitized user response
	safeUser := gin.H{
		"user_id":   user.UserID,
		"email":     user.Email,
		"username":  user.Username,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"role":      user.Role,
		"provider":  user.Provider,
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Login successful",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          safeUser,
	})
}


