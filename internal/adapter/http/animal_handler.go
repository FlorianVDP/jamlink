package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tindermals-backend/internal/adapter/http/middleware"
	"tindermals-backend/internal/modules/animal/usecase"
	"tindermals-backend/internal/shared/security"
)

type AnimalHandler struct {
	CreateAnimalUseCase  *animalUseCase.CreateAnimalUseCase
	GetAnimalListUseCase *animalUseCase.GetAnimalListUseCase
	GetAnimalByIdUseCase *animalUseCase.GetAnimalByIdUseCase

	SecurityService security.SecurityService
}

func NewAnimalHandler(router *gin.Engine, createAnimalUC *animalUseCase.CreateAnimalUseCase, getAnimalListUC *animalUseCase.GetAnimalListUseCase, getAnimalByIdUC *animalUseCase.GetAnimalByIdUseCase, securitySvc security.SecurityService) {
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

// CreateAnimal creates a new animal
// @Summary Create a new animal
// @Description Add a new animal with name, age, sexe, etc.
// @Tags Animals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param animal body animalUseCase.CreateAnimalInput true "Animal data"
// @Success 201 {object} animalDomain.Animal
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /animals [post]
func (h *AnimalHandler) CreateAnimal(c *gin.Context) {
	var input animalUseCase.CreateAnimalInput
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

// GetAnimalList retrieves all animals
// @Summary Get all animals
// @Description Retrieve the list of all registered animals
// @Tags Animals
// @Produce json
// @Success 200 {array} animalDomain.Animal
// @Failure 500 {object} map[string]string
// @Router /animals [get]
func (h *AnimalHandler) GetAnimalList(c *gin.Context) {
	allAnimals, err := h.GetAnimalListUseCase.Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, allAnimals)
}

// GetAnimalById retrieves an animal by its ID
// @Summary Get animal by ID
// @Description Retrieve a specific animal by its unique identifier
// @Tags Animals
// @Produce json
// @Param id path string true "Animal ID"
// @Success 200 {object} animalDomain.Animal
// @Failure 404 {object} map[string]string
// @Router /animals/{id} [get]
func (h *AnimalHandler) GetAnimalById(c *gin.Context) {
	id := c.Param("id")
	animal, err := h.GetAnimalByIdUseCase.Execute(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, animal)
}
