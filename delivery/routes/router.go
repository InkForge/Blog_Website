package routes

import (
	"github.com/InkForge/Blog_Website/delivery/controllers"
	auth "github.com/InkForge/Blog_Website/infrastructures/auth"
	infrastructures "github.com/InkForge/Blog_Website/infrastructures/auth"
	"github.com/gin-gonic/gin"
)

// RegisterBlogRoutes registers blog-related routes.
func RegisterBlogRoutes(router *gin.Engine, blogController *controllers.BlogController, authService *infrastructures.AuthService) {
	// Public routes
	router.GET("/blogs", blogController.GetAllBlogs)
	router.GET("/blogs/:id", authService.AuthWithRole("USER", "ADMIN"), blogController.GetBlogByID)

	authGroup := router.Group("/")
	authGroup.Use(authService.AuthWithRole("USER", "ADMIN"))
	{
		authGroup.POST("/blogs", blogController.CreateBlog)
		authGroup.PUT("/blogs/:id", blogController.UpdateBlog)
		authGroup.DELETE("/blogs/:id", blogController.DeleteBlog)
	}

	// Search and filter endpoints (public)
	router.GET("/blogs/search", blogController.Search)
	router.GET("/blogs/filter", blogController.FilterBlogs)
}

// RegisterBlogReactionRoutes registers blog reaction routes.
func RegisterBlogReactionRoutes(router *gin.Engine, blogReactionController *controllers.BlogReactionController, authService *infrastructures.AuthService) {
	authGroup := router.Group("/")
	authGroup.Use(authService.AuthWithRole("USER", "ADMIN"))
	{
		authGroup.POST("/blogs/:id/like", blogReactionController.LikeBlog)
		authGroup.POST("/blogs/:id/dislike", blogReactionController.DislikeBlog)
		authGroup.POST("/blogs/:id/unlike", blogReactionController.UnlikeBlog)
		authGroup.POST("/blogs/:id/undislike", blogReactionController.UndislikeBlog)
	}
}
func NewUserRoutes(userController *controllers.UserController,group gin.RouterGroup,authService *auth.AuthService){
	group.GET("/",authService.AuthWithRole("ADMIN"),userController.GetUsers)
	group.GET("/:id",authService.AuthWithRole("ADMIN"),userController.GetUserByID)
	group.GET("/me",authService.AuthWithRole("USER","ADMIN"),userController.GetMyProfile)
	group.PUT("/me",authService.AuthWithRole("USER","ADMIN"),userController.UpdateProfile)
	group.DELETE("/:id",authService.AuthWithRole("ADMIN"),userController.DeleteUser)
	group.GET("/search",authService.AuthWithRole("USER","ADMIN"),userController.SearchUsers)

}

func NewAuthRouter(authController controllers.AuthController, authService auth.AuthService, group gin.RouterGroup) {

	group.POST("/register", authController.Register)
	group.POST("/login", authController.Login)
	group.GET("/verify", authController.VerifyEmail)
	group.POST("/resend", authController.ResendVerification)
	group.POST("/forget", authController.RequestPasswordReset)
	group.POST("/reset", authController.ResetPassword)
	group.POST("/logout", authService.AuthWithRole("USER", "ADMIN"), authController.Logout)
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
	router.GET("/blogs/:id/comments", commentController.GetBlogComments)

	authGroup := router.Group("/")
	authGroup.Use(authService.AuthWithRole("USER", "ADMIN"))
	{
		// Comment CRUD
		authGroup.POST("/blogs/:id/comments", commentController.AddComment)
		authGroup.PUT("/comments/:commentID", commentController.UpdateComment)
		authGroup.DELETE("/blogs/:id/comments/:commentID", commentController.RemoveComment)

		// Comment Reactions
		authGroup.POST("/comments/:commentID/react/:status", commentReactionController.ReactToComment)
		authGroup.GET("/comments/:commentID/reaction", commentReactionController.GetUserReaction)
	}
}

func RegisterOAuthRoutes(
	router *gin.Engine,
	oauthController *controllers.OAuth2Controller,
) {
	oauth := router.Group("/oauth")
	{
		oauth.GET("/:provider/login", oauthController.RedirectToProvider)
		oauth.GET("/:provider/callback", oauthController.HandleCallback)
	}
}

func NewAIRouter(aiController *controllers.AIController, authService *infrastructures.AuthService, group gin.RouterGroup) {

	// authenticated routes - require logged-in users
	groupAuth := group.Group("/")
	groupAuth.Use(authService.AuthWithRole("USER", "ADMIN"))
	{
		groupAuth.POST("/suggest-tags", aiController.SuggestTags)
		groupAuth.POST("/summarize", aiController.Summarize)
		groupAuth.POST("/generate-title", aiController.GenerateTitle)
		groupAuth.POST("/suggest-content", aiController.SuggestContent)
		groupAuth.POST("/improve-content", aiController.ImproveContent)
		groupAuth.POST("/chat", aiController.Chat)
	}
}

func SetupRouter(
	commentController *controllers.CommentController,
	commentReactionController *controllers.CommentReactionController,
	blogController *controllers.BlogController,
	blogReactionController *controllers.BlogReactionController,
	authService *infrastructures.AuthService,
	authController *controllers.AuthController,
	oauthController *controllers.OAuth2Controller,
	userController *controllers.UserController,
	aiController *controllers.AIController,
) *gin.Engine {
	router := gin.Default()

	// Register comment & reaction routes
	RegisterCommentAndReactionRoutes(router, commentController, commentReactionController, authService)

	// Register blog routes
	RegisterBlogRoutes(router, blogController, authService)

	// Register blog reaction routes
	RegisterBlogReactionRoutes(router, blogReactionController, authService)

	// Auth routes
	authGroup := router.Group("/auth")
	NewAuthRouter(*authController, *authService, *authGroup)

	RegisterOAuthRoutes(router, oauthController)

	//user routes

	userGroup :=router.Group("/users")
	NewUserRoutes(userController, *userGroup,authService)

	
	// ai integration routes
	aiGroup := router.Group("/ai")
	NewAIRouter(aiController, authService, *aiGroup)

	return router
}
