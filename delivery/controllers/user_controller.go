package controllers

// imports
import (
	"github.com/InkForge/Blog_Website/domain"
)

// user controller
type UserController struct {
	userUseCase domain.IUserUseCase      // user usecase for user operations
}

