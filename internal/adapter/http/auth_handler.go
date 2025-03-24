package http

import (
	"github.com/gin-gonic/gin"
	"jamlink-backend/internal/modules/user/usecase"
	"net/http"
	"time"
)

type AuthHandler struct {
	CreateUserUseCase          *userUseCase.CreateUserUseCase
	LoginUserUseCase           *userUseCase.LoginUserUseCase
	LoginUserWithGoogleUseCase *userUseCase.LoginUserWithGoogleUseCase
	RefreshTokenUseCase        *userUseCase.RefreshTokenUseCase
}

func NewAuthHandler(router *gin.Engine, createUserUC *userUseCase.CreateUserUseCase, loginUserUC *userUseCase.LoginUserUseCase, loginWithGoogleUserUC *userUseCase.LoginUserWithGoogleUseCase, refreshTokenUC *userUseCase.RefreshTokenUseCase) {
	handler := &AuthHandler{
		CreateUserUseCase:          createUserUC,
		LoginUserUseCase:           loginUserUC,
		RefreshTokenUseCase:        refreshTokenUC,
		LoginUserWithGoogleUseCase: loginWithGoogleUserUC,
	}

	router.POST("/auth/register", handler.RegisterUser)
	router.POST("/auth/login", handler.LoginUser)
	router.POST("/auth/login/google", handler.LoginUserWithGoogle)
	router.POST("/auth/refresh-token", handler.RefreshToken)
}

// RegisterUser register a new user
// @Summary Register a new user
// @Description Create a new user account.
// @Description Password must:
// @Description - Be between 8 and 64 characters
// @Description - Contain at least one uppercase letter
// @Description - Contain at least one lowercase letter
// @Description - Contain at least one digit
// @Description - Contain at least one special character (e.g. !@#$%^&*)
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body userUseCase.CreateUserInput true "User credentials"
// @Success 201 {object} userDomain.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
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
// @Description Authenticate a user with email and password and store the refresh token (stored in HttpOnly cookie named 'refresh_token')
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body userUseCase.LoginUserInput true "Login credentials"
// @Success 200 {object} userUseCase.LoginUserOutput
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
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

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    output.RefreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	c.JSON(http.StatusOK, output.Token)
}

// LoginUserWithGoogle login a user with Google account
// @Summary Login a user with Google account
// @Description Authenticate a user with Google account and store the refresh token (stored in HttpOnly cookie named 'refresh_token')
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body userUseCase.LoginUserWithGoogleInput true "Login credentials"
// @Success 200 {object} userUseCase.LoginUserWithGoogleOutput
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/login/google [post]
func (h *AuthHandler) LoginUserWithGoogle(c *gin.Context) {
	var input userUseCase.LoginUserWithGoogleInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := h.LoginUserWithGoogleUseCase.Execute(input)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    output.RefreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	c.JSON(http.StatusOK, output.Token)
}

// RefreshToken refresh a token for a user
// @Summary Refresh a token
// @Description Refresh the JWT token using the refresh token (stored in HttpOnly cookie named 'refresh_token')
// @Tags Auth
// @Produce json
// @Success 200 {object} userUseCase.RefreshTokenOutput
// @Failure 401 {object} map[string]string
// @Router /auth/refresh-token [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	cookie, err := c.Request.Cookie("refresh_token")

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No refresh token"})
		return
	}

	input := userUseCase.RefreshTokenInput{RefreshToken: cookie.Value}

	output, err := h.RefreshTokenUseCase.Execute(input)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    output.RefreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	c.JSON(http.StatusOK, output)
}
