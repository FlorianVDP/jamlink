package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "jamlink-backend/docs"
	"jamlink-backend/internal/adapter/http"
	"jamlink-backend/internal/infra/db"
	emailinfra "jamlink-backend/internal/infra/email"
	userRepository "jamlink-backend/internal/modules/user/repository"
	userUsecase "jamlink-backend/internal/modules/user/usecase"
	"jamlink-backend/internal/shared/lang"
	"jamlink-backend/internal/shared/security"
)

// @title Jamlink API
// @version 1.0
// @description This is an API with Swagger and Gin.
// @host localhost:8080
// @BasePath /
func main() {
	_ = godotenv.Load()
	database := db.ConnectDB()
	db.MigrateDB(database)

	// Repositories
	userRepo := userRepository.NewPostgresUserRepository(database)

	// Services
	securityService := security.NewSecurityService()
	emailService := emailinfra.NewBrevoEmailService()
	langService := lang.NewLangNormalizer()

	// Use Cases
	createUserUseCase := userUsecase.NewCreateUserUseCase(userRepo, securityService, emailService)
	loginUserUseCase := userUsecase.NewLoginUserUseCase(userRepo, securityService)
	loginUserWithGoogleUseCase := userUsecase.NewLoginUserWithGoogleUseCase(userRepo, securityService)
	refreshTokenUseCase := userUsecase.NewRefreshTokenUseCase(securityService)
	verifyUserUseCase := userUsecase.NewVerifyUserUseCase(userRepo, securityService)

	// Setup router
	r := gin.Default()

	http.NewAuthHandler(r, langService, createUserUseCase, loginUserUseCase, loginUserWithGoogleUseCase, refreshTokenUseCase, verifyUserUseCase)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler))

	// Run server
	if err := r.Run(":8080"); err != nil {
		return
	}
}
