package controllers

// imports
import (
	"net/http"

	"github.com/InkForge/Blog_Website/domain"
	"github.com/gin-gonic/gin"
)

// user controller
type UserController struct {
	userUseCase domain.IUserUseCase      // user usecase for user operations
}

// user promote to admin role controller
func (uc *UserController) PromoteToAdmin(c *gin.Context) {
	userID := c.Param("id")

	err := uc.userUseCase.PromoteToAdmin(c.Request.Context(), userID)
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
	userID := c.Param("id")

	err := uc.userUseCase.DemoteFromAdmin(c.Request.Context(), userID)
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
