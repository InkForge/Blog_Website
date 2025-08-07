package main

import (
	"log"
	"time"

	"github.com/InkForge/Blog_Website/delivery/controllers"
	"github.com/InkForge/Blog_Website/delivery/routes"
	infrastructures "github.com/InkForge/Blog_Website/infrastructures/auth"
	infrastructures2 "github.com/InkForge/Blog_Website/infrastructures"
	mongo "github.com/InkForge/Blog_Website/infrastructures/db/mongo"
	mongo2 "github.com/InkForge/Blog_Website/repositories/mongo"
	"github.com/InkForge/Blog_Website/repositories"
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

	passwordService := infrastructures.NewPasswordService()
	jwtService := infrastructures.NewJWTService(configs.AccessTokenSecret, configs.RefreshTokenSecret, userRepo)
	notificationService := infrastructures2.NewSMTPService(configs.SMTPHost, configs.SMTPPort, configs.SMTPUsername, configs.SMTPPassword, configs.EmailFrom)
	txManager := mongo2.NewMongoTransactionManager(client)


	providersConfigs, err := infrastructures2.BuildProviderConfigs()
	if err != nil {
		log.Fatal("error: ", err)
	}
	
	authService := infrastructures.NewAuthService(jwtService, configs.JWTSecretKey)
	oauth2Service, err:= infrastructures.NewOAuth2Service(providersConfigs)



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

	commentController := controllers.NewCommentController(commentUsecase)
	commentReactionController := controllers.NewCommentReactionController(commentReactionUsecase)
	authController := controllers.NewAuthController(authUsecase)
	oauthController := controllers.NewOAuth2Controller(oauth2Service, authUsecase)

	r := routes.SetupRouter(commentController, commentReactionController, authService, authController, oauthController)

	r.Run(":" + configs.AppPort)
}
