package routes

import (

	"github.com/InkForge/Blog_Website/delivery/controllers"
	auth "github.com/InkForge/Blog_Website/infrastructures/auth"
	"github.com/gin-gonic/gin"
)

func newAuthRouter(authController controllers.AuthController, authService auth.AuthService, group gin.RouterGroup) {

	group.POST("/register", authController.Register)
	group.POST("/login", authController.Login)
	group.GET("/verify", authController.VerifyEmail)
	group.POST("/resend", authController.ResendVerification)
	group.POST("/forget", authController.RequestPasswordReset)
	group.POST("/reset", authController.ResetPassword)
	group.POST("/logout", authService.AuthWithRole("USER", "ADMIN"),authController.Logout)
	group.GET("/refresh", authService.AuthWithRole("USER", "ADMIN"), authController.RefreshToken)
}
