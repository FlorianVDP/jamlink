package userUseCase

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"jamlink-backend/internal/modules/auth/domain/user"
	"jamlink-backend/internal/modules/auth/mocks"
	"testing"
)

func TestCreateUser_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockSecurity := new(mocks.MockSecurityService)

	useCase := NewCreateUserUseCase(mockRepo, mockSecurity)

	input := CreateUserInput{
		Email:         "test@example.com",
		Password:      "Password123@",
		PreferredLang: "fr-FR",
	}

	mockRepo.On("FindByEmail", input.Email).Return(nil, errors.New("user not found"))
	mockSecurity.On("HashPassword", input.Password).Return("hashedpassword123", nil)
	mockRepo.On("Create", mock.Anything).Return(nil)
	user, err := useCase.Execute(input)

	assert.NoError(t, err)
	if assert.NotNil(t, user) {
		assert.Equal(t, input.Email, user.Email, input.PreferredLang)
		assert.Equal(t, "hashedpassword123", user.Password)
	}
}

func TestCreateUser_EmailAlreadyExists(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockSecurity := new(mocks.MockSecurityService)

	useCase := NewCreateUserUseCase(mockRepo, mockSecurity)

	input := CreateUserInput{
		Email:         "test@example.com",
		Password:      "Password123@",
		PreferredLang: "fr-FR",
	}

	mockRepo.On("FindByEmail", input.Email).Return(&user.User{Email: input.Email}, nil)

	user, err := useCase.Execute(input)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "email already exists", err.Error())
}

func TestCreateUser_FailOnHashing(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockSecurity := new(mocks.MockSecurityService)

	useCase := NewCreateUserUseCase(mockRepo, mockSecurity)

	input := CreateUserInput{
		Email:         "test@example.com",
		Password:      "Password123@",
		PreferredLang: "fr-FR",
	}

	mockRepo.On("FindByEmail", input.Email).Return(nil, errors.New("user not found"))
	mockSecurity.On("HashPassword", input.Password).Return("", errors.New("hashing error"))

	user, err := useCase.Execute(input)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "hashing error", err.Error())
}
