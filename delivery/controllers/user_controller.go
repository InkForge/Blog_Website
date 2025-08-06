package controllers

// imports
import (
	"context"
	"net/http"
	"time"

	"github.com/InkForge/Blog_Website/domain"
	"github.com/gin-gonic/gin"
)

// user controller
type UserController struct {
	UserUseCase domain.IUserUseCase      // user usecase for user operations
}
func NewUserController(userUsecase domain.IUserUseCase)*UserController{
	return &UserController{
		UserUseCase: userUsecase,
	}
}

// user promote to admin role controller
func (uc *UserController) PromoteToAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5 * time.Second)
    defer cancel()


	userID := c.Param("id")

	err := uc.UserUseCase.PromoteToAdmin(ctx, userID)
	if err != nil {
		switch err {
		case domain.ErrInvalidUserID, domain.ErrInvalidRole:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case domain.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User promoted to admin successfully"})
}


// user demote from admin role controller
func (uc *UserController) DemoteFromAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5 * time.Second)
    defer cancel()
	userID := c.Param("id")

	err := uc.UserUseCase.DemoteFromAdmin(ctx, userID)
	if err != nil {
		switch err {
		case domain.ErrInvalidUserID, domain.ErrInvalidRole:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case domain.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User demoted to regular user successfully"})
}

// GetUserByID handles GET /users/:id
func (uc *UserController) GetUserByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5 * time.Second)
    defer cancel()
	userID := c.Param("id")

	user, err := uc.UserUseCase.GetUserByID(ctx, userID)
	if err != nil {
		switch err {
		case domain.ErrInvalidUserID:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		case domain.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetUsers handles GET /users
func (uc *UserController) GetUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5 * time.Second)
    defer cancel()
	users, err := uc.UserUseCase.GetUsers(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

 //DeleteUser handles DELETE /users/:id
func (uc *UserController) DeleteUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5 * time.Second)
    defer cancel()
	userID := c.Param("id")

	err := uc.UserUseCase.DeleteUserByID(ctx, userID)
	if err != nil {
		switch err {
		case domain.ErrInvalidUserID:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		case domain.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// SearchUsers handles GET /users/search?q=query
func (uc *UserController) SearchUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5 * time.Second)
    defer cancel()
	query := c.Query("q")

	users, err := uc.UserUseCase.SearchUsers(ctx, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetMyProfile handles GET /users/me
func (uc *UserController) GetMyProfile(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5 * time.Second)
    defer cancel()
	userID:= c.GetString("userID")
	if  userID=="" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := uc.UserUseCase.GetMyData(ctx, userID)
	if err != nil {
		switch err {
		case domain.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get profile"})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateProfile handles PUT /users/me
func (uc *UserController) UpdateProfile(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5 * time.Second)
    defer cancel()
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Ensure the user can only update their own profile
	user.UserID = userID.(string)

	err := uc.UserUseCase.UpdateProfile(ctx, &user)
	if err != nil {
		switch err {
		case domain.ErrInvalidUserID:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

