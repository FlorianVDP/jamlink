package main

import (
	"tindermals-backend/internal/adapter/http"
	"tindermals-backend/internal/infra"
	"tindermals-backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	repo := infra.NewMemoryAnimalRepository()
	createAnimalUseCase := usecase.NewCreateAnimalUseCase(repo)
	getAnimalListUseCase := usecase.NewGetAnimalListUseCase(repo)
	getAnimalByIdUseCase := usecase.NewGetAnimalByIdUseCase(repo)

	r := gin.Default()

	http.NewAnimalHandler(r, createAnimalUseCase, getAnimalListUseCase, getAnimalByIdUseCase)

	r.Run(":8080")
}
