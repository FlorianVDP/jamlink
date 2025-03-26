package http

import (
	"github.com/gin-gonic/gin"
	"jamlink-backend/internal/adapter/http/middleware"
	"jamlink-backend/internal/modules/auth/usecase"
	"jamlink-backend/internal/shared/lang"
	"jamlink-backend/internal/shared/security"
	"net/http"
	"time"
)

type AuthHandler struct {
	securitySvc                   security.SecurityService
	LangNormalizer                lang.LangNormalizer
	CreateUserUseCase             *useCase.CreateUserUseCase
	LoginUserUseCase              *useCase.LoginUserUseCase
	LoginUserWithGoogleUseCase    *useCase.LoginUserWithGoogleUseCase
	RefreshTokenUseCase           *useCase.RefreshTokenUseCase
	VerifyUserUseCase             *useCase.VerifyUserUseCase
	RequestVerifyUserEmailUseCase *useCase.RequestVerifyUserEmailUseCase
	RequestResetPasswordUseCase   *useCase.RequestResetPasswordUseCase
	ResetPasswordUseCase          *useCase.ResetPasswordUseCase
	DisconnectUserUseCase         *useCase.DisconnectUserUseCase
}

func NewAuthHandler(router *gin.Engine, securitySvc security.SecurityService, langNormalizer lang.LangNormalizer, createUserUC *useCase.CreateUserUseCase, loginUserUC *useCase.LoginUserUseCase, loginWithGoogleUserUC *useCase.LoginUserWithGoogleUseCase, refreshTokenUC *useCase.RefreshTokenUseCase, verifyUserUC *useCase.VerifyUserUseCase, getVerificationTokenUC *useCase.RequestVerifyUserEmailUseCase, requestResetPasswordUC *useCase.RequestResetPasswordUseCase, resetPasswordUseCase *useCase.ResetPasswordUseCase, disconnectUserUseCase *useCase.DisconnectUserUseCase) {
	handler := &AuthHandler{
		securitySvc:                   securitySvc,
		LangNormalizer:                langNormalizer,
		CreateUserUseCase:             createUserUC,
		LoginUserUseCase:              loginUserUC,
		RefreshTokenUseCase:           refreshTokenUC,
		LoginUserWithGoogleUseCase:    loginWithGoogleUserUC,
		VerifyUserUseCase:             verifyUserUC,
		RequestVerifyUserEmailUseCase: getVerificationTokenUC,
		RequestResetPasswordUseCase:   requestResetPasswordUC,
		ResetPasswordUseCase:          resetPasswordUseCase,
		DisconnectUserUseCase:         disconnectUserUseCase,
	}

	router.POST("/auth/register", handler.RegisterUser)
	router.POST("/auth/login", handler.LoginUser)
	router.POST("/auth/login/google", handler.LoginUserWithGoogle)
	router.POST("/auth/refresh-token", handler.RefreshToken)
	router.POST("/auth/verify", handler.VerifyUser)
	router.POST("/auth/request-verify-user", handler.RequestVerifyUserEmail)
	router.POST("/auth/request-reset-password", handler.RequestResetPassword)
	router.POST("/auth/reset-password", handler.ResetPassword)
	router.POST("/auth/logout", handler.LogoutUser)

	// Protected routes
	protected := router.Group("/")
	protected.Use(middleware.JWTAuthMiddleware(securitySvc))

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
// @Param input body useCase.CreateUserInput true "User credentials"
// @Success 201 {object} user.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/register [post]
func (h *AuthHandler) RegisterUser(c *gin.Context) {
	var input useCase.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rawLang := c.GetHeader("Accept-Language")
	normalizedLang := h.LangNormalizer.Normalize(rawLang)

	input.PreferredLang = normalizedLang
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
// @Param credentials body useCase.LoginUserInput true "Login credentials"
// @Success 200 {object} useCase.LoginUserOutput
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func (h *AuthHandler) LoginUser(c *gin.Context) {
	var input useCase.LoginUserInput

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
// @Param credentials body useCase.LoginUserWithGoogleInput true "Login credentials"
// @Success 200 {object} useCase.LoginUserWithGoogleOutput
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/login/google [post]
func (h *AuthHandler) LoginUserWithGoogle(c *gin.Context) {
	var input useCase.LoginUserWithGoogleInput

	rawLang := c.GetHeader("Accept-Language")
	normalizedLang := h.LangNormalizer.Normalize(rawLang)

	input.PreferredLang = normalizedLang

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
// @Success 200 {object} useCase.RefreshTokenOutput
// @Failure 401 {object} map[string]string
// @Router /auth/refresh-token [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	cookie, err := c.Request.Cookie("refresh_token")

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No refresh token"})
		return
	}

	input := useCase.RefreshTokenInput{RefreshToken: cookie.Value}

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

// VerifyUser verify a user
// @Summary Verify a user
// @Description Verify a user account using the token received in the email
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body useCase.VerifyUserInput true "Verification token"
// @Success 200
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/verify [post]
func (h *AuthHandler) VerifyUser(c *gin.Context) {

	var input useCase.VerifyUserInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.VerifyUserUseCase.Execute(input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)

}

// RequestVerifyUserEmail get a verification token
// @Summary Get a verification token
// @Description Get a verification token for a user
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body useCase.RequestVerifyUserEmailInput true "User email"
// @Success 200
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/get-verification-token [post]
func (h *AuthHandler) RequestVerifyUserEmail(c *gin.Context) {
	var input useCase.RequestVerifyUserEmailInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.RequestVerifyUserEmailUseCase.Execute(input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// RequestResetPassword request a password reset email
// @Summary Request a password reset email
// @Description Request a password reset email for a user
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body useCase.RequestResetPasswordInput true "User email"
// @Success 200
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/request-reset-password [post]
func (h *AuthHandler) RequestResetPassword(c *gin.Context) {
	var input useCase.RequestResetPasswordInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.RequestResetPasswordUseCase.Execute(input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// ResetPassword reset a user password
// @Summary Reset a user password
// @Description Reset a user password using the token received in the email
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body useCase.ResetPasswordInput true "Reset password credentials"
// @Success 200
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var input useCase.ResetPasswordInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.ResetPasswordUseCase.Execute(input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// LogoutUser logout a user
// @Summary Logout a user
// @Description Logout a user and delete the refresh token
// @Tags Auth
// @Produce json
// @Success 200
// @Failure 401 {object} map[string]string
// @Router /auth/logout [post]
func (h *AuthHandler) LogoutUser(c *gin.Context) {
	cookie, err := c.Request.Cookie("refresh_token")

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No refresh token"})
		return
	}
	input := &useCase.DisconnectUserInput{RefreshToken: cookie.Value}

	err = h.DisconnectUserUseCase.Execute(input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	c.Status(http.StatusOK)
}
