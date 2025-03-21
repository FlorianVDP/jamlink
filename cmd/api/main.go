package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "tindermals-backend/docs"
	"tindermals-backend/internal/adapter/http"
	"tindermals-backend/internal/infra/db"
	animalRepository "tindermals-backend/internal/modules/animal/repository"
	animalUsecase "tindermals-backend/internal/modules/animal/usecase"
	userRepository "tindermals-backend/internal/modules/user/repository"
	userUsecase "tindermals-backend/internal/modules/user/usecase"
	"tindermals-backend/internal/shared/security"
)

// @title Tindermals API
// @version 1.0
// @description This is an API with Swagger and Gin.
// @host localhost:8080
// @BasePath /

func main() {
	_ = godotenv.Load()
	database := db.ConnectDB()
	db.MigrateDB(database)

	// Repositories
	animalRepo := animalRepository.NewPostgresAnimalRepository(database)
	userRepo := userRepository.NewPostgresUserRepository(database)

	// Services
	securityService := security.NewSecurityService()

	// Use Cases
	createAnimalUseCase := animalUsecase.NewCreateAnimalUseCase(animalRepo)
	getAnimalListUseCase := animalUsecase.NewGetAnimalListUseCase(animalRepo)
	getAnimalByIdUseCase := animalUsecase.NewGetAnimalByIdUseCase(animalRepo)

	createUserUseCase := userUsecase.NewCreateUserUseCase(userRepo, securityService)
	loginUserUseCase := userUsecase.NewLoginUserUseCase(userRepo, securityService)

	// Setup router
	r := gin.Default()

	http.NewAnimalHandler(r, createAnimalUseCase, getAnimalListUseCase, getAnimalByIdUseCase, securityService)
	http.NewAuthHandler(r, createUserUseCase, loginUserUseCase)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler))

	// Run server
	if err := r.Run(":8080"); err != nil {
		return
	}
}
