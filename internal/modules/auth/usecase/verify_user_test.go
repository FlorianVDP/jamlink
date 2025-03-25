package userUseCase

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	userDomain "jamlink-backend/internal/modules/auth/domain/user"
	"jamlink-backend/internal/modules/auth/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"jamlink-backend/internal/shared/security"
)

func TestVerifyUserUseCase_Execute_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockSec := new(mocks.MockSecurityService)

	verifyUC := NewVerifyUserUseCase(mockRepo, mockSec)

	input := VerifyUserInput{
		Token: "test-token",
	}

	claims := jwt.MapClaims{
		"email": "test@example.com",
	}

	user := &userDomain.User{
		Email: "test@example.com",
		Verification: userDomain.UserVerification{
			IsVerified: false,
			VerifiedAt: nil,
		},
	}

	mockSec.
		On("ValidateJWT", "test-token").
		Return(claims, nil)

	mockRepo.
		On("FindByEmail", "test@example.com").
		Return(user, nil)

	mockRepo.
		On("Update", mock.AnythingOfType("*user.User")).
		Return(nil)

	err := verifyUC.Execute(input)

	assert.NoError(t, err)
	assert.True(t, user.Verification.IsVerified)
	assert.NotNil(t, user.Verification.VerifiedAt)
	assert.WithinDuration(t, time.Now(), *user.Verification.VerifiedAt, time.Second)

	mockSec.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestVerifyUserUseCase_Execute_InvalidToken(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockSec := new(mocks.MockSecurityService)
	verifyUC := NewVerifyUserUseCase(mockRepo, mockSec)

	input := VerifyUserInput{
		Token: "invalid-token",
	}

	claims := jwt.MapClaims{
		"mail": "notfound@example.com",
	}

	mockSec.
		On("ValidateJWT", "invalid-token").
		Return(claims, errors.New("invalid token"))

	err := verifyUC.Execute(input)

	assert.Error(t, err, security.ErrInvalidToken)

	mockSec.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestVerifyUserUseCase_Execute_NoMailInClaims(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockSec := new(mocks.MockSecurityService)
	verifyUC := NewVerifyUserUseCase(mockRepo, mockSec)

	input := VerifyUserInput{
		Token: "some-token",
	}

	claims := jwt.MapClaims{
		"foo": "bar",
	}

	mockSec.
		On("ValidateJWT", "some-token").
		Return(claims, nil)

	err := verifyUC.Execute(input)

	assert.ErrorIs(t, err, security.ErrInvalidUserEmail)
	mockSec.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "FindByEmail", mock.Anything)
}
