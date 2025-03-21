package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tindermals-backend/internal/adapter/http/middleware"
	"tindermals-backend/internal/modules/animal/usecase"
	"tindermals-backend/internal/shared/security"
)

type AnimalHandler struct {
	CreateAnimalUseCase  *animalUsecase.CreateAnimalUseCase
	GetAnimalListUseCase *animalUsecase.GetAnimalListUseCase
	GetAnimalByIdUseCase *animalUsecase.GetAnimalByIdUseCase

	SecurityService security.SecurityService
}

func NewAnimalHandler(router *gin.Engine, createAnimalUC *animalUsecase.CreateAnimalUseCase, getAnimalListUC *animalUsecase.GetAnimalListUseCase, getAnimalByIdUC *animalUsecase.GetAnimalByIdUseCase, securitySvc security.SecurityService) {
	handler := &AnimalHandler{
		CreateAnimalUseCase:  createAnimalUC,
		GetAnimalListUseCase: getAnimalListUC,
		GetAnimalByIdUseCase: getAnimalByIdUC,

		SecurityService: securitySvc,
	}

	// 🔓 Public routes
	router.GET("/animals", handler.GetAnimalList)
	router.GET("/animals/:id", handler.GetAnimalById)

	// 🔒 Protected routes
	protected := router.Group("/")
	protected.Use(middleware.JWTAuthMiddleware(securitySvc))
	protected.POST("/animals", handler.CreateAnimal)
}

func (h *AnimalHandler) CreateAnimal(c *gin.Context) {
	var input animalUsecase.CreateAnimalInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	animal, err := h.CreateAnimalUseCase.Execute(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, animal)
}

func (h *AnimalHandler) GetAnimalList(c *gin.Context) {
	allAnimals, err := h.GetAnimalListUseCase.Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, allAnimals)
}

func (h *AnimalHandler) GetAnimalById(c *gin.Context) {
	id := c.Param("id")
	animal, err := h.GetAnimalByIdUseCase.Execute(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, animal)
}
