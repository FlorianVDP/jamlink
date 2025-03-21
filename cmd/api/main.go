package main

import (
	"tindermals-backend/internal/adapter/http"
	"tindermals-backend/internal/infra/db"
	"tindermals-backend/internal/repository"
	"tindermals-backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	database := db.ConnectDB()
	db.MigrateDB(database)
	repo := repository.NewPostgresAnimalRepository(database)
	createAnimalUseCase := usecase.NewCreateAnimalUseCase(repo)
	getAnimalListUseCase := usecase.NewGetAnimalListUseCase(repo)
	getAnimalByIdUseCase := usecase.NewGetAnimalByIdUseCase(repo)

	r := gin.Default()

	http.NewAnimalHandler(r, createAnimalUseCase, getAnimalListUseCase, getAnimalByIdUseCase)

	err := r.Run(":8080")

	if err != nil {
		return
	}
}
