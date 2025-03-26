package userUseCase

import (
	"errors"
	"github.com/stretchr/testify/mock"
	tokenDomain "jamlink-backend/internal/modules/auth/domain/token"
	"jamlink-backend/internal/modules/auth/domain/user"
	"jamlink-backend/internal/shared/security"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"jamlink-backend/internal/modules/auth/mocks"
)

func TestLoginUser_Success(t *testing.T) {
	mockSecurity := new(mocks.MockSecurityService)
	userRepo := new(mocks.MockUserRepository)
	tokenRepo := new(mocks.MockTokenRepository)
	createdUser := &user.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		Password: "hashedpassword",
		Verification: user.UserVerification{
			IsVerified: false,
		},
	}

	input := LoginUserInput{
		Email:    "test@example.com",
		Password: "password123",
	}

	const expiringTimeForRefreshToken = time.Hour * 24 * 7

	refreshToken := "valid_refresh_token"
	accessToken := "valid_access_token"

	userRepo.On("FindByEmail", input.Email).Return(createdUser, nil)
	mockSecurity.On("CheckPassword", input.Password, createdUser.Password).Return(true)

	mockSecurity.On("GenerateJWT", &createdUser.ID, (*string)(nil), time.Minute*15, "login", createdUser.Verification.IsVerified).Return(accessToken, nil)
	mockSecurity.On("GenerateJWT", &createdUser.ID, (*string)(nil), expiringTimeForRefreshToken, "refresh_token", createdUser.Verification.IsVerified).Return(refreshToken, nil)

	tokenRepo.On("Create", mock.MatchedBy(func(token *tokenDomain.Token) bool {
		return token.UserID == createdUser.ID && token.Token == refreshToken
	})).Return(nil)

	usecase := NewLoginUserUseCase(userRepo, mockSecurity, tokenRepo)
	output, err := usecase.Execute(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, refreshToken, output.RefreshToken)
	assert.Equal(t, accessToken, output.Token)

	userRepo.AssertExpectations(t)
	mockSecurity.AssertExpectations(t)
	tokenRepo.AssertExpectations(t)
}

func TestLoginUser_InvalidPassword(t *testing.T) {
	mockSecurity := new(mocks.MockSecurityService)
	userRepo := new(mocks.MockUserRepository)
	tokenRepo := new(mocks.MockTokenRepository)

	user := &user.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	input := LoginUserInput{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	userRepo.On("FindByEmail", input.Email).Return(user, nil)
	mockSecurity.On("CheckPassword", input.Password, user.Password).Return(false)

	usecase := NewLoginUserUseCase(userRepo, mockSecurity, tokenRepo)
	output, err := usecase.Execute(input)

	assert.Error(t, err)
	assert.Nil(t, output)
	assert.ErrorIs(t, err, security.ErrPasswordComparison)

}

func TestLoginUser_UserNotFound(t *testing.T) {
	mockSecurity := new(mocks.MockSecurityService)
	userRepo := new(mocks.MockUserRepository)
	tokenRepo := new(mocks.MockTokenRepository)

	input := LoginUserInput{
		Email:    "notfound@example.com",
		Password: "password123",
	}

	userRepo.On("FindByEmail", input.Email).Return(nil, errors.New("not found"))

	usecase := NewLoginUserUseCase(userRepo, mockSecurity, tokenRepo)
	output, err := usecase.Execute(input)

	assert.Error(t, err)
	assert.Nil(t, output)
}
