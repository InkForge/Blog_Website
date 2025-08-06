package routes

import (
	"github.com/InkForge/Blog_Website/delivery/controllers"
	"github.com/gin-gonic/gin"
)

// RegisterCommentRoutes registers comment-related routes.
func RegisterCommentRoutes(router *gin.Engine, commentController *controllers.CommentController) {
	router.POST("/blogs/:blogID/comments", commentController.AddComment)
	router.GET("/blogs/:blogID/comments", commentController.GetBlogComments)
	router.PUT("/comments/:commentID", commentController.UpdateComment)
	router.DELETE("/blogs/:blogID/comments/:commentID", commentController.RemoveComment)
}

// RegisterCommentReactionRoutes registers comment reaction-related routes.
func RegisterCommentReactionRoutes(router *gin.Engine, commentReactionController *controllers.CommentReactionController) {
	router.POST("/comments/:commentID/react/:status", commentReactionController.ReactToComment)
	router.GET("/comments/:commentID/reaction", commentReactionController.GetUserReaction)
}

// SetupRouter initializes the Gin engine and registers all routes.
func SetupRouter(commentController *controllers.CommentController, commentReactionController *controllers.CommentReactionController) *gin.Engine {
	router := gin.Default()

	RegisterCommentRoutes(router, commentController)
	RegisterCommentReactionRoutes(router, commentReactionController)

	return router
}