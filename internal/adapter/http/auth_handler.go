package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tindermals-backend/internal/modules/user/usecase"
)

type AuthHandler struct {
	CreateUserUseCase *userUsecase.CreateUserUseCase
	LoginUserUseCase  *userUsecase.LoginUserUseCase
}

func NewAuthHandler(router *gin.Engine, createUserUC *userUsecase.CreateUserUseCase, loginUserUC *userUsecase.LoginUserUseCase) {
	handler := &AuthHandler{
		CreateUserUseCase: createUserUC,
		LoginUserUseCase:  loginUserUC,
	}

	router.POST("/auth/register", handler.RegisterUser)
	router.GET("/auth/login", handler.LoginUser)
}

func (h *AuthHandler) RegisterUser(c *gin.Context) {
	var input userUsecase.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.CreateUserUseCase.Execute(input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *AuthHandler) LoginUser(c *gin.Context) {
	var input userUsecase.LoginUserInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := h.LoginUserUseCase.Execute(input)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, output)
}
