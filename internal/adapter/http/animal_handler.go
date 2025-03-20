package http

import (
	"net/http"
	"tindermals-backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

type AnimalHandler struct {
	CreateAnimalUseCase  *usecase.CreateAnimalUseCase
	GetAnimalListUseCase *usecase.GetAnimalListUseCase
	GetAnimalByIdUseCase *usecase.GetAnimalByIdUseCase
}

func NewAnimalHandler(router *gin.Engine, createAnimalUC *usecase.CreateAnimalUseCase, getAnimalListUC *usecase.GetAnimalListUseCase, getAnimalByIdUC *usecase.GetAnimalByIdUseCase) {
	handler := &AnimalHandler{
		CreateAnimalUseCase:  createAnimalUC,
		GetAnimalListUseCase: getAnimalListUC,
		GetAnimalByIdUseCase: getAnimalByIdUC,
	}

	router.POST("/animals", handler.CreateAnimal)
	router.GET("/animals", handler.GetAnimalList)
	router.GET("/animals/:id", handler.GetAnimalById)
}

func (h *AnimalHandler) CreateAnimal(c *gin.Context) {
	var input usecase.CreateAnimalInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	animal, err := h.CreateAnimalUseCase.Execute(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, animal)
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
