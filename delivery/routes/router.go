package routes

import (
	"github.com/InkForge/Blog_Website/delivery/controllers"
	infrastructures "github.com/InkForge/Blog_Website/infrastructures/auth"
	"github.com/gin-gonic/gin"
)

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