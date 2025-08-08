package main

import (
	"log"
	"time"

	"github.com/InkForge/Blog_Website/delivery/controllers"
	"github.com/InkForge/Blog_Website/delivery/routes"
	infrastructures2 "github.com/InkForge/Blog_Website/infrastructures"
	infrastructures "github.com/InkForge/Blog_Website/infrastructures/auth"
	mongo "github.com/InkForge/Blog_Website/infrastructures/db/mongo"
	"github.com/InkForge/Blog_Website/repositories"
	mongo2 "github.com/InkForge/Blog_Website/repositories/mongo"
	"github.com/InkForge/Blog_Website/usecases"
)

func main() {
	configs, err := infrastructures2.LoadConfig()
	if err != nil {
		log.Fatal("error: ", err)
	}

	client := mongo.NewMongoClient()
	db := client.Database(configs.DBName)

	userRepo := repositories.NewUserRepository(db)
	commentRepo := repositories.NewCommentMongoRepository(db)
	commentReactionRepo := repositories.NewCommentReactionMongoRepository(db)
	blogRepo := repositories.NewBlogMongoRepository(db)
	blogReactionRepo := repositories.NewBlogReactionRepository(db)
	blogViewRepo := repositories.NewBlogViewRepository(db)
	tagRepo := repositories.NewTagMongoRepository(db)

	passwordService := infrastructures.NewPasswordService()
	jwtService := infrastructures.NewJWTService(configs.AccessTokenSecret, configs.RefreshTokenSecret, userRepo)
	notificationService := infrastructures2.NewSMTPService(configs.SMTPHost, configs.SMTPPort, configs.SMTPUsername, configs.SMTPPassword, configs.EmailFrom)
	txManager := mongo2.NewMongoTransactionManager(client)

	providersConfigs, err := infrastructures2.BuildProviderConfigs()
	if err != nil {
		log.Fatal("error: ", err)
	}

	authService := infrastructures.NewAuthService(jwtService, configs.JWTSecretKey)
	oauth2Service, err := infrastructures.NewOAuth2Service(providersConfigs)

	
	blogUsecase := usecases.NewBlogUsecase(blogRepo, blogViewRepo, tagRepo, userRepo, txManager)
	blogReactionUsecase := usecases.NewBlogReactionUseCase(blogRepo, blogReactionRepo, txManager)
	
	userUsecase:=usecases.NewUserUseCase(userRepo)

	commentUsecase := usecases.NewCommentUsecase(blogRepo, commentRepo, txManager)
	commentReactionUsecase := usecases.NewCommentReactionUsecase(commentRepo, commentReactionRepo, txManager)
	authUsecase := usecases.NewAuthUseCase(
		userRepo,
		passwordService,
		jwtService,
		notificationService,
		configs.BaseURL,
		time.Second*10,
	)
	blogController := controllers.NewBlogController(blogUsecase)
	blogReactionController := controllers.NewBlogReactionController(blogReactionUsecase)
	commentController := controllers.NewCommentController(commentUsecase)
	commentReactionController := controllers.NewCommentReactionController(commentReactionUsecase)
	authController := controllers.NewAuthController(authUsecase)
	oauthController := controllers.NewOAuth2Controller(oauth2Service, authUsecase)
	userControler:=controllers.NewUserController(userUsecase)

	r := routes.SetupRouter(commentController, commentReactionController, blogController, blogReactionController, authService, authController, oauthController,userControler)

	r.Run(":" + configs.AppPort)
}
