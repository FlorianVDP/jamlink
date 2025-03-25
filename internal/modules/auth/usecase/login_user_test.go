package userUseCase

import (
	"errors"
	"jamlink-backend/internal/modules/auth/domain/user"
	"jamlink-backend/internal/shared/security"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"jamlink-backend/internal/modules/auth/mocks"
)

func TestLoginUser_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockSecurity := new(mocks.MockSecurityService)

	user := &user.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	input := LoginUserInput{
		Email:    "test@example.com",
		Password: "password123",
	}

	mockRepo.On("FindByEmail", input.Email).Return(user, nil)
	mockSecurity.On("CheckPassword", input.Password, user.Password).Return(true)

	mockSecurity.On("GenerateJWT", &user.ID, (*string)(nil), time.Minute*15, "login").Return("mocked.jwt.token", nil)
	mockSecurity.On("GenerateJWT", &user.ID, (*string)(nil), time.Hour*24*7, "refresh_token").Return("mocked.jwt.token", nil)

	usecase := NewLoginUserUseCase(mockRepo, mockSecurity)
	output, err := usecase.Execute(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, "mocked.jwt.token", output.Token)
}

func TestLoginUser_InvalidPassword(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockSecurity := new(mocks.MockSecurityService)

	user := &user.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	input := LoginUserInput{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	mockRepo.On("FindByEmail", input.Email).Return(user, nil)
	mockSecurity.On("CheckPassword", input.Password, user.Password).Return(false)

	usecase := NewLoginUserUseCase(mockRepo, mockSecurity)
	output, err := usecase.Execute(input)

	assert.Error(t, err)
	assert.Nil(t, output)
	assert.ErrorIs(t, err, security.ErrPasswordComparison)

}

func TestLoginUser_UserNotFound(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockSecurity := new(mocks.MockSecurityService)

	input := LoginUserInput{
		Email:    "notfound@example.com",
		Password: "password123",
	}

	mockRepo.On("FindByEmail", input.Email).Return(nil, errors.New("not found"))

	usecase := NewLoginUserUseCase(mockRepo, mockSecurity)
	output, err := usecase.Execute(input)

	assert.Error(t, err)
	assert.Nil(t, output)
}
