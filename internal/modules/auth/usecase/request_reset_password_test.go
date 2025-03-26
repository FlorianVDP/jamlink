package useCase

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	userDomain "jamlink-backend/internal/modules/auth/domain/user"
	"jamlink-backend/internal/modules/auth/mocks"
	"jamlink-backend/internal/shared/email"
	"testing"
	"time"
)

func TestRequestResetPassword_Success(t *testing.T) {
	// Arrange
	mockTokenRepo := new(mocks.MockTokenRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockSecurity := new(mocks.MockSecurityService)
	mockEmailService := new(mocks.MockEmailService)

	userID := uuid.New()
	userEmail := "user@example.com"
	jwtToken := "jwt.token.string"

	user := &userDomain.User{
		ID:            userID,
		Email:         userEmail,
		PreferredLang: "fr",
		Verification: userDomain.UserVerification{
			IsVerified: true,
		},
	}

	// Setup expectations
	mockUserRepo.On("FindByEmail", userEmail).Return(user, nil)
	mockSecurity.On("GenerateJWT",
		mock.MatchedBy(func(id *uuid.UUID) bool { return *id == userID }),
		mock.MatchedBy(func(email *string) bool { return *email == userEmail }),
		time.Minute*15,
		"reset_password",
		user.Verification.IsVerified).Return(jwtToken, nil)
	mockTokenRepo.On("Create", mock.AnythingOfType("*token.Token")).Return(nil)
	mockEmailService.On("Send",
		userEmail,
		email.TemplateResetPassword,
		"fr",
		mock.AnythingOfType("map[string]string")).Return(nil)

	useCase := NewRequestResetPasswordUseCase(mockTokenRepo, mockUserRepo, mockSecurity, mockEmailService)

	// Act
	err := useCase.Execute(RequestResetPasswordInput{
		Email: userEmail,
	})

	// Assert
	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
	mockSecurity.AssertExpectations(t)
	mockTokenRepo.AssertExpectations(t)
	mockEmailService.AssertExpectations(t)
}

func TestRequestResetPassword_UserNotFound(t *testing.T) {
	// Arrange
	mockTokenRepo := new(mocks.MockTokenRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockSecurity := new(mocks.MockSecurityService)
	mockEmailService := new(mocks.MockEmailService)

	userEmail := "nonexistent@example.com"

	// Setup expectations
	mockUserRepo.On("FindByEmail", userEmail).Return(nil, errors.New("utilisateur non trouvé"))

	useCase := NewRequestResetPasswordUseCase(mockTokenRepo, mockUserRepo, mockSecurity, mockEmailService)

	// Act
	err := useCase.Execute(RequestResetPasswordInput{
		Email: userEmail,
	})

	// Assert
	assert.Error(t, err)
	mockUserRepo.AssertExpectations(t)
	mockSecurity.AssertNotCalled(t, "GenerateJWT")
	mockTokenRepo.AssertNotCalled(t, "Create")
	mockEmailService.AssertNotCalled(t, "Send")
}

func TestRequestResetPassword_JWTGenerationFailed(t *testing.T) {
	// Arrange
	mockTokenRepo := new(mocks.MockTokenRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockSecurity := new(mocks.MockSecurityService)
	mockEmailService := new(mocks.MockEmailService)

	userID := uuid.New()
	userEmail := "user@example.com"

	user := &userDomain.User{
		ID:    userID,
		Email: userEmail,
		Verification: userDomain.UserVerification{
			IsVerified: true,
		},
	}

	// Setup expectations
	mockUserRepo.On("FindByEmail", userEmail).Return(user, nil)
	mockSecurity.On("GenerateJWT",
		mock.MatchedBy(func(id *uuid.UUID) bool { return *id == userID }),
		mock.MatchedBy(func(email *string) bool { return *email == userEmail }),
		time.Minute*15,
		"reset_password",
		true).Return("", errors.New("erreur lors de la génération du JWT"))

	useCase := NewRequestResetPasswordUseCase(mockTokenRepo, mockUserRepo, mockSecurity, mockEmailService)

	// Act
	err := useCase.Execute(RequestResetPasswordInput{
		Email: userEmail,
	})

	// Assert
	assert.Error(t, err)
	mockUserRepo.AssertExpectations(t)
	mockSecurity.AssertExpectations(t)
	mockTokenRepo.AssertNotCalled(t, "Create")
	mockEmailService.AssertNotCalled(t, "Send")
}

func TestRequestResetPassword_TokenCreationFailed(t *testing.T) {
	// Arrange
	mockTokenRepo := new(mocks.MockTokenRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockSecurity := new(mocks.MockSecurityService)
	mockEmailService := new(mocks.MockEmailService)

	userID := uuid.New()
	userEmail := "user@example.com"
	jwtToken := "jwt.token.string"

	user := &userDomain.User{
		ID:    userID,
		Email: userEmail,
		Verification: userDomain.UserVerification{
			IsVerified: true,
		},
	}

	// Setup expectations
	mockUserRepo.On("FindByEmail", userEmail).Return(user, nil)
	mockSecurity.On("GenerateJWT",
		mock.MatchedBy(func(id *uuid.UUID) bool { return *id == userID }),
		mock.MatchedBy(func(email *string) bool { return *email == userEmail }),
		time.Minute*15,
		"reset_password",
		true).Return(jwtToken, nil)
	mockTokenRepo.On("Create", mock.AnythingOfType("*token.Token")).Return(errors.New("erreur de création du token"))

	useCase := NewRequestResetPasswordUseCase(mockTokenRepo, mockUserRepo, mockSecurity, mockEmailService)

	// Act
	err := useCase.Execute(RequestResetPasswordInput{
		Email: userEmail,
	})

	// Assert
	assert.Error(t, err)
	mockUserRepo.AssertExpectations(t)
	mockSecurity.AssertExpectations(t)
	mockTokenRepo.AssertExpectations(t)
	mockEmailService.AssertNotCalled(t, "Send")
}

func TestRequestResetPassword_EmailSendingFailed(t *testing.T) {
	// Arrange
	mockTokenRepo := new(mocks.MockTokenRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockSecurity := new(mocks.MockSecurityService)
	mockEmailService := new(mocks.MockEmailService)

	userID := uuid.New()
	userEmail := "user@example.com"
	jwtToken := "jwt.token.string"

	user := &userDomain.User{
		ID:            userID,
		Email:         userEmail,
		PreferredLang: "fr",
		Verification: userDomain.UserVerification{
			IsVerified: true,
		},
	}

	// Setup expectations
	mockUserRepo.On("FindByEmail", userEmail).Return(user, nil)
	mockSecurity.On("GenerateJWT",
		mock.MatchedBy(func(id *uuid.UUID) bool { return *id == userID }),
		mock.MatchedBy(func(email *string) bool { return *email == userEmail }),
		time.Minute*15,
		"reset_password",
		true).Return(jwtToken, nil)
	mockTokenRepo.On("Create", mock.AnythingOfType("*token.Token")).Return(nil)
	mockEmailService.On("Send",
		userEmail,
		email.TemplateResetPassword,
		"fr",
		mock.AnythingOfType("map[string]string")).Return(errors.New("erreur d'envoi d'email"))

	useCase := NewRequestResetPasswordUseCase(mockTokenRepo, mockUserRepo, mockSecurity, mockEmailService)

	// Act
	err := useCase.Execute(RequestResetPasswordInput{
		Email: userEmail,
	})

	// Assert
	assert.Error(t, err)
	mockUserRepo.AssertExpectations(t)
	mockSecurity.AssertExpectations(t)
	mockTokenRepo.AssertExpectations(t)
	mockEmailService.AssertExpectations(t)
}
