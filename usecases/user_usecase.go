package usecases

import (
	"time"

	"github.com/InkForge/Blog_Website/domain"
)

type UserUseCase struct {
	UserRepo            domain.IUserRepository
	PasswordService     domain.IPasswordService
	JWTService          domain.IJWTService
	NotificationService domain.INotificationService
	BaseURL             string
	ContextTimeout      time.Duration
}

func NewUserUseCase(repo domain.IUserRepository, ps domain.IPasswordService, jw domain.IJWTService, ns domain.INotificationService, bs string, timeout time.Duration) domain.IUserUseCase {
	return &UserUseCase{
		UserRepo:            repo,
		PasswordService:     ps,
		JWTService:          jw,
		NotificationService: ns,
		BaseURL:             bs,
		ContextTimeout:      timeout,
	}
}