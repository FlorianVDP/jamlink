package userUseCase

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	tokenDomain "jamlink-backend/internal/modules/auth/domain/token"
	userDomain "jamlink-backend/internal/modules/auth/domain/user"
	"jamlink-backend/internal/modules/auth/mocks"
	"jamlink-backend/internal/shared/security"
	"testing"
	"time"
)

func TestResetPassword_SuccessfulReset(t *testing.T) {
	// Arrange
	mockSecurity := new(mocks.MockSecurityService)
	mockUserRepo := new(mocks.MockUserRepository)
	mockTokenRepo := new(mocks.MockTokenRepository)

	validToken := "valid.jwt.token"
	tokenID := uuid.New()
	userID := uuid.New()
	email := "user@example.com"

	// Mock claims for token validation
	claims := jwt.MapClaims{
		"type":  "reset_password",
		"exp":   float64(time.Now().Add(time.Hour).Unix()),
		"email": email,
	}

	// Create mock token
	token := &tokenDomain.Token{
		ID:        tokenID,
		UserID:    userID,
		Token:     validToken,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Hour),
	}

	// Create user
	user := &userDomain.User{
		ID:       userID,
		Email:    email,
		Password: "old-hashed-password",
	}

	// Setup expectations
	mockSecurity.On("ValidateJWT", validToken).Return(claims, nil)
	mockTokenRepo.On("FindByToken", validToken).Return(token, nil)
	mockUserRepo.On("FindByEmail", email).Return(user, nil)
	mockSecurity.On("HashPassword", "NewSecurePassword123!").Return("new-hashed-password", nil)
	mockUserRepo.On("Update", mock.MatchedBy(func(u *userDomain.User) bool {
		return u.Password == "new-hashed-password" && u.Email == email
	})).Return(nil)
	mockTokenRepo.On("DeleteByID", tokenID).Return(nil)

	useCase := NewResetPasswordUseCase(mockTokenRepo, mockUserRepo, mockSecurity)

	// Act
	err := useCase.Execute(ResetPasswordInput{
		Token:                 validToken,
		NewPassword:           "NewSecurePassword123!",
		NewPasswordValidation: "NewSecurePassword123!",
	})

	// Assert
	assert.NoError(t, err)

	mockSecurity.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockTokenRepo.AssertExpectations(t)
}

func TestResetPassword_PasswordMismatch(t *testing.T) {
	// Arrange
	mockSecurity := new(mocks.MockSecurityService)
	mockUserRepo := new(mocks.MockUserRepository)
	mockTokenRepo := new(mocks.MockTokenRepository)

	useCase := NewResetPasswordUseCase(mockTokenRepo, mockUserRepo, mockSecurity)

	// Act
	err := useCase.Execute(ResetPasswordInput{
		Token:                 "valid.jwt.token",
		NewPassword:           "Password123!",
		NewPasswordValidation: "DifferentPassword123!",
	})

	// Assert
	assert.ErrorIs(t, err, tokenDomain.ErrPasswordDoesntMatch)

	mockSecurity.AssertNotCalled(t, "ValidateJWT")
	mockUserRepo.AssertNotCalled(t, "FindByEmail")
	mockTokenRepo.AssertNotCalled(t, "FindByToken")
}

func TestResetPassword_InvalidPassword(t *testing.T) {
	// Arrange
	mockSecurity := new(mocks.MockSecurityService)
	mockUserRepo := new(mocks.MockUserRepository)
	mockTokenRepo := new(mocks.MockTokenRepository)

	useCase := NewResetPasswordUseCase(mockTokenRepo, mockUserRepo, mockSecurity)

	// Act
	err := useCase.Execute(ResetPasswordInput{
		Token:                 "valid.jwt.token",
		NewPassword:           "weak",
		NewPasswordValidation: "weak",
	})

	// Assert
	assert.Error(t, err)
	mockSecurity.AssertNotCalled(t, "ValidateJWT")
	mockTokenRepo.AssertNotCalled(t, "FindByToken")
	mockUserRepo.AssertNotCalled(t, "FindByEmail")
}

func TestResetPassword_InvalidToken(t *testing.T) {
	// Arrange
	mockSecurity := new(mocks.MockSecurityService)
	mockUserRepo := new(mocks.MockUserRepository)
	mockTokenRepo := new(mocks.MockTokenRepository)

	invalidToken := "invalid.jwt.token"

	emptyClaims := jwt.MapClaims{}
	mockSecurity.On("ValidateJWT", invalidToken).Return(emptyClaims, errors.New("token invalide"))

	useCase := NewResetPasswordUseCase(mockTokenRepo, mockUserRepo, mockSecurity)

	// Act
	err := useCase.Execute(ResetPasswordInput{
		Token:                 invalidToken,
		NewPassword:           "NewSecurePassword123!",
		NewPasswordValidation: "NewSecurePassword123!",
	})

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token invalide")

	mockSecurity.AssertExpectations(t)
	mockTokenRepo.AssertNotCalled(t, "FindByToken")
	mockUserRepo.AssertNotCalled(t, "FindByEmail")
}

func TestResetPassword_InvalidTokenType(t *testing.T) {
	// Arrange
	mockSecurity := new(mocks.MockSecurityService)
	mockUserRepo := new(mocks.MockUserRepository)
	mockTokenRepo := new(mocks.MockTokenRepository)

	validToken := "valid.jwt.token"

	// Mock claims for token validation with wrong token type
	claims := jwt.MapClaims{
		"type":  "wrong_token_type",
		"exp":   float64(time.Now().Add(time.Hour).Unix()),
		"email": "user@example.com",
	}

	mockSecurity.On("ValidateJWT", validToken).Return(claims, nil)

	useCase := NewResetPasswordUseCase(mockTokenRepo, mockUserRepo, mockSecurity)

	// Act
	err := useCase.Execute(ResetPasswordInput{
		Token:                 validToken,
		NewPassword:           "NewSecurePassword123!",
		NewPasswordValidation: "NewSecurePassword123!",
	})

	// Assert
	assert.ErrorIs(t, err, tokenDomain.ErrTokenType)

	mockSecurity.AssertExpectations(t)
	mockTokenRepo.AssertNotCalled(t, "FindByToken")
	mockUserRepo.AssertNotCalled(t, "FindByEmail")
}

func TestResetPassword_ExpiredToken(t *testing.T) {
	// Arrange
	mockSecurity := new(mocks.MockSecurityService)
	mockUserRepo := new(mocks.MockUserRepository)
	mockTokenRepo := new(mocks.MockTokenRepository)

	validToken := "valid.jwt.token"

	// Mock claims for token validation with expired token
	claims := jwt.MapClaims{
		"type":  "reset_password",
		"exp":   float64(time.Now().Add(-1 * time.Hour).Unix()),
		"email": "user@example.com",
	}

	mockSecurity.On("ValidateJWT", validToken).Return(claims, nil)

	useCase := NewResetPasswordUseCase(mockTokenRepo, mockUserRepo, mockSecurity)

	// Act
	err := useCase.Execute(ResetPasswordInput{
		Token:                 validToken,
		NewPassword:           "NewSecurePassword123!",
		NewPasswordValidation: "NewSecurePassword123!",
	})

	// Assert
	assert.ErrorIs(t, err, tokenDomain.ErrTokenExpired)

	mockSecurity.AssertExpectations(t)
	mockTokenRepo.AssertNotCalled(t, "FindByToken")
	mockUserRepo.AssertNotCalled(t, "FindByEmail")
}

func TestResetPassword_MissingEmail(t *testing.T) {
	// Arrange
	mockSecurity := new(mocks.MockSecurityService)
	mockUserRepo := new(mocks.MockUserRepository)
	mockTokenRepo := new(mocks.MockTokenRepository)

	validToken := "valid.jwt.token"

	// Mock claims for token validation with missing email
	claims := jwt.MapClaims{
		"type": "reset_password",
		"exp":  float64(time.Now().Add(time.Hour).Unix()),
		//"email": "user@example.com",
	}

	mockSecurity.On("ValidateJWT", validToken).Return(claims, nil)

	useCase := NewResetPasswordUseCase(mockTokenRepo, mockUserRepo, mockSecurity)

	// Act
	err := useCase.Execute(ResetPasswordInput{
		Token:                 validToken,
		NewPassword:           "NewSecurePassword123!",
		NewPasswordValidation: "NewSecurePassword123!",
	})

	// Assert
	assert.ErrorIs(t, err, security.ErrInvalidUserEmail)

	mockSecurity.AssertExpectations(t)
	mockTokenRepo.AssertNotCalled(t, "FindByToken")
	mockUserRepo.AssertNotCalled(t, "FindByEmail")
}

func TestResetPassword_TokenNotFound(t *testing.T) {
	// Arrange
	mockSecurity := new(mocks.MockSecurityService)
	mockUserRepo := new(mocks.MockUserRepository)
	mockTokenRepo := new(mocks.MockTokenRepository)

	validToken := "valid.jwt.token"

	// Mock claims for token validation
	claims := jwt.MapClaims{
		"type":  "reset_password",
		"exp":   float64(time.Now().Add(time.Hour).Unix()),
		"email": "user@example.com",
	}

	mockSecurity.On("ValidateJWT", validToken).Return(claims, nil)
	mockTokenRepo.On("FindByToken", validToken).Return(nil, tokenDomain.ErrTokenNotFound)

	useCase := NewResetPasswordUseCase(mockTokenRepo, mockUserRepo, mockSecurity)

	// Act
	err := useCase.Execute(ResetPasswordInput{
		Token:                 validToken,
		NewPassword:           "NewSecurePassword123!",
		NewPasswordValidation: "NewSecurePassword123!",
	})

	// Assert
	assert.ErrorIs(t, err, tokenDomain.ErrTokenNotFound)

	mockSecurity.AssertExpectations(t)
	mockTokenRepo.AssertExpectations(t)
	mockUserRepo.AssertNotCalled(t, "FindByEmail")
}

func TestResetPassword_UserNotFound(t *testing.T) {
	// Arrange
	mockSecurity := new(mocks.MockSecurityService)
	mockUserRepo := new(mocks.MockUserRepository)
	mockTokenRepo := new(mocks.MockTokenRepository)

	validToken := "valid.jwt.token"
	tokenID := uuid.New()
	userID := uuid.New()
	email := "user@example.com"

	// Mock claims for token validation
	claims := jwt.MapClaims{
		"type":  "reset_password",
		"exp":   float64(time.Now().Add(time.Hour).Unix()),
		"email": "user@example.com",
	}

	// Create mock token
	token := &tokenDomain.Token{
		ID:        tokenID,
		UserID:    userID,
		Token:     validToken,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Hour),
	}

	mockSecurity.On("ValidateJWT", validToken).Return(claims, nil)
	mockTokenRepo.On("FindByToken", validToken).Return(token, nil)
	mockUserRepo.On("FindByEmail", email).Return(nil, userDomain.ErrUserNotFound)

	useCase := NewResetPasswordUseCase(mockTokenRepo, mockUserRepo, mockSecurity)

	// Act
	err := useCase.Execute(ResetPasswordInput{
		Token:                 validToken,
		NewPassword:           "NewSecurePassword123!",
		NewPasswordValidation: "NewSecurePassword123!",
	})

	// Assert
	assert.ErrorIs(t, err, userDomain.ErrUserNotFound)

	mockSecurity.AssertExpectations(t)
	mockTokenRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestResetPassword_HashPasswordError(t *testing.T) {
	// Arrange
	mockSecurity := new(mocks.MockSecurityService)
	mockUserRepo := new(mocks.MockUserRepository)
	mockTokenRepo := new(mocks.MockTokenRepository)

	validToken := "valid.jwt.token"
	tokenID := uuid.New()
	userID := uuid.New()
	email := "user@example.com"

	// Mock claims for token validation
	claims := jwt.MapClaims{
		"type":  "reset_password",
		"exp":   float64(time.Now().Add(time.Hour).Unix()),
		"email": "user@example.com",
	}

	// Create mock token
	token := &tokenDomain.Token{
		ID:        tokenID,
		UserID:    userID,
		Token:     validToken,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Hour),
	}

	// Create user
	user := &userDomain.User{
		ID:       userID,
		Email:    email,
		Password: "old-hashed-password",
	}

	mockSecurity.On("ValidateJWT", validToken).Return(claims, nil)
	mockTokenRepo.On("FindByToken", validToken).Return(token, nil)
	mockUserRepo.On("FindByEmail", email).Return(user, nil)
	mockSecurity.On("HashPassword", "NewSecurePassword123!").Return("", errors.New("erreur de hashage"))

	useCase := NewResetPasswordUseCase(mockTokenRepo, mockUserRepo, mockSecurity)

	// Act
	err := useCase.Execute(ResetPasswordInput{
		Token:                 validToken,
		NewPassword:           "NewSecurePassword123!",
		NewPasswordValidation: "NewSecurePassword123!",
	})

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "erreur de hashage")

	mockSecurity.AssertExpectations(t)
	mockTokenRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockUserRepo.AssertNotCalled(t, "Update")
}
