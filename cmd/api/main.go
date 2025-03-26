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
	userRepository "jamlink-backend/internal/modules/auth/repository"
	userUsecase "jamlink-backend/internal/modules/auth/usecase"
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
	tokenRepo := userRepository.NewPostgresTokenRepository(database)

	// Services
	securityService := security.NewSecurityService()
	emailService := emailinfra.NewBrevoEmailService()
	langService := lang.NewLangNormalizer()

	// Use Cases
	createUserUseCase := userUsecase.NewCreateUserUseCase(userRepo, securityService)
	loginUserUseCase := userUsecase.NewLoginUserUseCase(userRepo, securityService, tokenRepo)
	loginUserWithGoogleUseCase := userUsecase.NewLoginUserWithGoogleUseCase(userRepo, securityService)
	refreshTokenUseCase := userUsecase.NewRefreshTokenUseCase(securityService, userRepo, tokenRepo)
	requestVerifyUserEmailUseCase := userUsecase.NewRequestVerifyUserEmailUseCase(securityService, userRepo, emailService)
	verifyUserUseCase := userUsecase.NewVerifyUserUseCase(userRepo, securityService)
	requestResetPasswordUseCase := userUsecase.NewRequestResetPasswordUseCase(tokenRepo, userRepo, securityService, emailService)
	resetPasswordUseCase := userUsecase.NewResetPasswordUseCase(tokenRepo, userRepo, securityService)
	disconnectUserUseCase := userUsecase.NewDisconnectUserUseCase(tokenRepo)

	// Setup router
	r := gin.Default()

	http.NewAuthHandler(r, securityService, langService, createUserUseCase, loginUserUseCase, loginUserWithGoogleUseCase, refreshTokenUseCase, verifyUserUseCase, requestVerifyUserEmailUseCase, requestResetPasswordUseCase, resetPasswordUseCase, disconnectUserUseCase)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler))

	// Run server
	if err := r.Run(":8080"); err != nil {
		return
	}
}
