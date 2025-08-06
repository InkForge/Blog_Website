package routes

import (
	"github.com/InkForge/Blog_Website/delivery/controllers"
	auth "github.com/InkForge/Blog_Website/infrastructures/auth"
  infrastructures "github.com/InkForge/Blog_Website/infrastructures/auth"
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

// RegisterCommentAndReactionRoutes registers both comment and comment reaction routes in one group.
func RegisterCommentAndReactionRoutes(
	router *gin.Engine,
	commentController *controllers.CommentController,
	commentReactionController *controllers.CommentReactionController,
	authService *infrastructures.AuthService,
) {
	// Public route: Anyone can view comments for a blog post
	router.GET("/blogs/:blogID/comments", commentController.GetBlogComments)

	// Authenticated group (users with "user" or "admin" roles)
	authGroup := router.Group("/")
	authGroup.Use(authService.AuthWithRole("user", "admin"))
	{
		// Comment CRUD
		authGroup.POST("/blogs/:blogID/comments", commentController.AddComment)
		authGroup.PUT("/comments/:commentID", commentController.UpdateComment)
		authGroup.DELETE("/blogs/:blogID/comments/:commentID", commentController.RemoveComment)

		// Comment Reactions
		authGroup.POST("/comments/:commentID/react/:status", commentReactionController.ReactToComment)
		authGroup.GET("/comments/:commentID/reaction", commentReactionController.GetUserReaction)
	}
}

func SetupRouter(
	commentController *controllers.CommentController,
	commentReactionController *controllers.CommentReactionController,
	authService *infrastructures.AuthService,
) *gin.Engine {
	
	router := gin.Default()

	// Register all comment & reaction routes 
	RegisterCommentAndReactionRoutes(router, commentController, commentReactionController, authService)

	return router
}