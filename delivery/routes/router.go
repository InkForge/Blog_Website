package routes

import (
	"fmt"
	"time"

	"github.com/InkForge/Blog_Website/delivery/controllers"
	infra "github.com/InkForge/Blog_Website/infrastructures"
	"github.com/InkForge/Blog_Website/infrastructures/auth"
	repo "github.com/InkForge/Blog_Website/repositories"
	"github.com/InkForge/Blog_Website/usecases"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func newAuthRouter(timeout time.Duration, db mongo.Database, group *gin.RouterGroup, config infra.Config) {
	userRepo := repo.NewUserRepository(&db)
	ps :=	infrastructures.NewPasswordService()
	idk := repo.NewRevocationRepository(&db)
	jwtService := infrastructures.NewJWTService(config.AccessTokenSecret, config.RefreshTokenSecret,  idk)
	notificationService := infra.NewSMTPService(config.SMTPHost, config.SMTPPort, config.SMTPUsername, config.SMTPPassword, config.EmailFrom)
	authUsecase := usecases.NewAuthUseCase(userRepo, ps, jwtService, notificationService, fmt.Sprintf("%s:%s", config.BaseURL, config.AppPort), timeout)
	authService := infrastructures.NewAuthService(jwtService , config.JWTSecretKey)
	authController := controllers.AuthController{AuthUsecase: authUsecase}
	group.POST("/register", authController.Register)
	group.POST("/login", authController.Login)
	group.POST("/logout", authService.AuthWithRole("USER", "ADMIN"),authController.Logout)
	group.GET("/verify", authController.VerifyEmail)
	group.POST("/resend", authController.ResendVerification)
	group.POST("/forget", authController.RequestPasswordReset)
	group.POST("/reset", authController.ResetPassword)
	group.GET("/refresh", authService.AuthWithRole("USER", "ADMIN"), authController.RefreshToken)
}

func newSomeThingHere() {
	// create your own repos here
	// setup the paths
}

func Setup (timeout time.Duration, db mongo.Database, router *gin.Engine, config infra.Config){
	Authentication := router.Group("/auth", )
	newAuthRouter(timeout, db, Authentication, config)
}
