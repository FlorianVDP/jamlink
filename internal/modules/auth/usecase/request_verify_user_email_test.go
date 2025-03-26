package useCase

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	userDomain "jamlink-backend/internal/modules/auth/domain/user"
	"jamlink-backend/internal/modules/auth/mocks"
	"jamlink-backend/internal/shared/email"
	"os"
	"testing"
	"time"
)

func TestGetVerificationEmail_Success(t *testing.T) {
	// Arrange
	mockSecurity := new(mocks.MockSecurityService)
	mockUserRepo := new(mocks.MockUserRepository)
	mockEmailService := new(mocks.MockEmailService)

	userEmail := "user@example.com"
	jwtToken := "verification.jwt.token"
	verifyURL := "https://example.com/verify"

	os.Setenv("FRONTEND_VERIFY_URL", verifyURL)
	defer os.Unsetenv("FRONTEND_VERIFY_URL")

	// Create a user with unverified status
	user := &userDomain.User{
		Email: userEmail,
		Verification: userDomain.UserVerification{
			IsVerified: false,
			VerifiedAt: nil,
		},
		PreferredLang: "en",
	}

	mockUserRepo.On("FindByEmail", userEmail).Return(user, nil)
	mockSecurity.On("GenerateJWT",
		mock.Anything,
		mock.MatchedBy(func(email *string) bool { return *email == userEmail }),
		mock.AnythingOfType("time.Duration"),
		"verify_email",
		false).Return(jwtToken, nil)

	mockEmailService.On("Send",
		userEmail,
		email.TemplateVerification,
		"en",
		mock.MatchedBy(func(data map[string]string) bool {
			expectedURL := verifyURL + "?token=" + jwtToken
			return data["URL"] == expectedURL
		})).Return(nil)

	// Create the use case
	useCase := NewRequestVerifyUserEmailUseCase(mockSecurity, mockUserRepo, mockEmailService)

	// Act
	err := useCase.Execute(RequestVerifyUserEmailInput{
		Email: userEmail,
	})

	// Assert
	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
	mockSecurity.AssertExpectations(t)
	mockEmailService.AssertExpectations(t)
}

func TestGetVerificationEmail_UserNotFound(t *testing.T) {
	// Arrange
	mockSecurity := new(mocks.MockSecurityService)
	mockUserRepo := new(mocks.MockUserRepository)
	mockEmailService := new(mocks.MockEmailService)

	userEmail := "nonexistent@example.com"

	mockUserRepo.On("FindByEmail", userEmail).Return(nil, errors.New("user not found"))

	// Create the use case
	useCase := NewRequestVerifyUserEmailUseCase(mockSecurity, mockUserRepo, mockEmailService)

	// Act
	err := useCase.Execute(RequestVerifyUserEmailInput{
		Email: userEmail,
	})

	// Assert
	assert.Error(t, err)
	assert.Equal(t, userDomain.ErrUserNotFound, err)
	mockUserRepo.AssertExpectations(t)
	mockSecurity.AssertNotCalled(t, "GenerateJWT")
	mockEmailService.AssertNotCalled(t, "Send")
}

func TestGetVerificationEmail_AlreadyVerified(t *testing.T) {
	// Arrange
	mockSecurity := new(mocks.MockSecurityService)
	mockUserRepo := new(mocks.MockUserRepository)
	mockEmailService := new(mocks.MockEmailService)

	userEmail := "verified@example.com"
	jwtToken := "verification.jwt.token"

	// Create a verified user
	verifiedTime := time.Now()
	user := &userDomain.User{
		Email: userEmail,
		Verification: userDomain.UserVerification{
			IsVerified: true,
			VerifiedAt: &verifiedTime,
		},
	}

	mockUserRepo.On("FindByEmail", userEmail).Return(user, nil)
	mockSecurity.On("GenerateJWT",
		mock.Anything,
		mock.MatchedBy(func(email *string) bool { return *email == userEmail }),
		mock.AnythingOfType("time.Duration"),
		"verify_email",
		true).Return(jwtToken, nil)

	// Create the use case
	useCase := NewRequestVerifyUserEmailUseCase(mockSecurity, mockUserRepo, mockEmailService)

	// Act
	err := useCase.Execute(RequestVerifyUserEmailInput{
		Email: userEmail,
	})

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "your account is already verified", err.Error())
	mockUserRepo.AssertExpectations(t)
	mockSecurity.AssertExpectations(t)
	mockEmailService.AssertNotCalled(t, "Send")
}

func TestGetVerificationEmail_JWTGenerationFailure(t *testing.T) {
	// Arrange
	mockSecurity := new(mocks.MockSecurityService)
	mockUserRepo := new(mocks.MockUserRepository)
	mockEmailService := new(mocks.MockEmailService)

	userEmail := "user@example.com"

	// Create an unverified user
	user := &userDomain.User{
		Email: userEmail,
		Verification: userDomain.UserVerification{
			IsVerified: false,
			VerifiedAt: nil,
		},
	}

	mockUserRepo.On("FindByEmail", userEmail).Return(user, nil)
	mockSecurity.On("GenerateJWT",
		mock.Anything,
		mock.MatchedBy(func(email *string) bool { return *email == userEmail }),
		mock.AnythingOfType("time.Duration"),
		"verify_email",
		false).Return("", errors.New("jwt generation error"))

	// Create the use case
	useCase := NewRequestVerifyUserEmailUseCase(mockSecurity, mockUserRepo, mockEmailService)

	// Act
	err := useCase.Execute(RequestVerifyUserEmailInput{
		Email: userEmail,
	})

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "jwt generation error", err.Error())
	mockUserRepo.AssertExpectations(t)
	mockSecurity.AssertExpectations(t)
	mockEmailService.AssertNotCalled(t, "Send")
}

func TestGetVerificationEmail_EmailSendingFailure(t *testing.T) {
	// Arrange
	mockSecurity := new(mocks.MockSecurityService)
	mockUserRepo := new(mocks.MockUserRepository)
	mockEmailService := new(mocks.MockEmailService)

	userEmail := "user@example.com"
	jwtToken := "verification.jwt.token"
	verifyURL := "https://example.com/verify"

	os.Setenv("FRONTEND_VERIFY_URL", verifyURL)
	defer os.Unsetenv("FRONTEND_VERIFY_URL")

	// Create an unverified user
	user := &userDomain.User{
		Email: userEmail,
		Verification: userDomain.UserVerification{
			IsVerified: false,
			VerifiedAt: nil,
		},
		PreferredLang: "fr",
	}

	mockUserRepo.On("FindByEmail", userEmail).Return(user, nil)
	mockSecurity.On("GenerateJWT",
		mock.Anything,
		mock.MatchedBy(func(email *string) bool { return *email == userEmail }),
		mock.AnythingOfType("time.Duration"),
		"verify_email",
		false).Return(jwtToken, nil)

	mockEmailService.On("Send",
		userEmail,
		email.TemplateVerification,
		"fr",
		mock.MatchedBy(func(data map[string]string) bool {
			expectedURL := verifyURL + "?token=" + jwtToken
			return data["URL"] == expectedURL
		})).Return(errors.New("email sending error"))

	// Create the use case
	useCase := NewRequestVerifyUserEmailUseCase(mockSecurity, mockUserRepo, mockEmailService)

	// Act
	err := useCase.Execute(RequestVerifyUserEmailInput{
		Email: userEmail,
	})

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "email sending error", err.Error())
	mockUserRepo.AssertExpectations(t)
	mockSecurity.AssertExpectations(t)
	mockEmailService.AssertExpectations(t)
}
