package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tindermals-backend/internal/modules/user/usecase"
)

type AuthHandler struct {
	CreateUserUseCase   *userUseCase.CreateUserUseCase
	LoginUserUseCase    *userUseCase.LoginUserUseCase
	RefreshTokenUseCase *userUseCase.RefreshTokenUseCase
}

func NewAuthHandler(router *gin.Engine, createUserUC *userUseCase.CreateUserUseCase, loginUserUC *userUseCase.LoginUserUseCase, refreshTokenUC *userUseCase.RefreshTokenUseCase) {
	handler := &AuthHandler{
		CreateUserUseCase:   createUserUC,
		LoginUserUseCase:    loginUserUC,
		RefreshTokenUseCase: refreshTokenUC,
	}

	router.POST("/auth/register", handler.RegisterUser)
	router.GET("/auth/login", handler.LoginUser)
	router.GET("/auth/refresh-token", handler.RefreshToken)
}

// RegisterUser register a new user
// @Summary Post a new user
// @Description Register a new user with email and password
// @Tags Auth
// @Produce json
// @Success 201 {object} userDomain.User
// @Failure 404 {object} map[string]string
// @Router /auth/register [post]
func (h *AuthHandler) RegisterUser(c *gin.Context) {
	var input userUseCase.CreateUserInput
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

// LoginUser login a user
// @Summary Login a user
// @Description Login a user with email and password
// @Tags Auth
// @Produce json
// @Success 200 {object} userUseCase.LoginUserOutput
// @Failure 404 {object} map[string]string
// @Router /auth/login [get]
func (h *AuthHandler) LoginUser(c *gin.Context) {
	var input userUseCase.LoginUserInput

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

// RefreshToken refresh a token for a user
// @Summary Refresh a token
// @Description Refresh a token for a user
// @Tags Auth
// @Produce json
// @Success 200 {object} userUseCase.RefreshTokenOutput
// @Failure 404 {object} map[string]string
// @Router /auth/refresh-token [get]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var input userUseCase.RefreshTokenInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := h.RefreshTokenUseCase.Execute(input)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, output)
}
