package userUseCase

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"jamlink-backend/internal/modules/user/domain"
	"jamlink-backend/internal/modules/user/mocks"
	"jamlink-backend/internal/shared/email"
	"testing"
)

func TestCreateUser_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockSecurity := new(mocks.MockSecurityService)
	mockEmail := new(mocks.MockEmailService)

	useCase := NewCreateUserUseCase(mockRepo, mockSecurity, mockEmail)

	input := CreateUserInput{
		Email:         "test@example.com",
		Password:      "Password123@",
		PreferredLang: "fr-FR",
	}

	mockRepo.On("FindByEmail", input.Email).Return(nil, errors.New("user not found"))
	mockSecurity.On("HashPassword", input.Password).Return("hashedpassword123", nil)
	mockRepo.On("Create", mock.Anything).Return(nil)
	mockSecurity.On("GenerateVerificationJWT", input.Email).Return("test@example.com", nil)
	mockEmail.On("Send", input.Email, email.TemplateVerification, "fr-FR", map[string]string{"URL": "?token=test@example.com"}).Return(nil)
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
	mockEmail := new(mocks.MockEmailService)

	useCase := NewCreateUserUseCase(mockRepo, mockSecurity, mockEmail)

	input := CreateUserInput{
		Email:         "test@example.com",
		Password:      "Password123@",
		PreferredLang: "fr-FR",
	}

	mockRepo.On("FindByEmail", input.Email).Return(&userDomain.User{Email: input.Email}, nil)

	user, err := useCase.Execute(input)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "email already exists", err.Error())
}

func TestCreateUser_FailOnHashing(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockSecurity := new(mocks.MockSecurityService)
	mockEmail := new(mocks.MockEmailService)

	useCase := NewCreateUserUseCase(mockRepo, mockSecurity, mockEmail)

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
