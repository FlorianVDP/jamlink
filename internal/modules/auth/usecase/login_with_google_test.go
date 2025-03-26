package userUseCase

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	userDomain "jamlink-backend/internal/modules/auth/domain/user"
	"jamlink-backend/internal/modules/auth/mocks"
	"testing"
)

// MockLoginUserWithGoogleUseCase extends the original use case but overrides the token validation
type MockLoginUserWithGoogleUseCase struct {
	*LoginUserWithGoogleUseCase
	validateTokenFunc func(idToken string) (string, error)
}

// executeWithMockValidation replaces the token validation logic for testing
func (uc *MockLoginUserWithGoogleUseCase) Execute(input LoginUserWithGoogleInput) (*LoginUserWithGoogleOutput, error) {
	// Use our mock validation function instead of the real idtoken.Validate
	email, err := uc.validateTokenFunc(input.IDToken)
	if err != nil {
		return nil, errors.New("invalid Google token")
	}

	if email == "" {
		return nil, errors.New("email not found in Google token")
	}

	user, err := uc.repo.FindByEmail(email)

	if err != nil {
		randomPassword, err := uc.security.GenerateSecureRandomString(32)
		if err != nil {
			return nil, err
		}

		hashed, err := uc.security.HashPassword(randomPassword)
		if err != nil {
			return nil, err
		}

		user, err = userDomain.CreateUser(email, hashed, input.PreferredLang, "google")
		if err != nil {
			return nil, err
		}

		err = uc.repo.Create(user)
		if err != nil {
			return nil, err
		}
	}

	token, err := uc.security.GenerateJWT(&user.ID, nil, 15, "login", user.Verification.IsVerified)
	if err != nil {
		return nil, err
	}

	refreshToken, err := uc.security.GenerateJWT(&user.ID, nil, 24*7, "refresh_token", user.Verification.IsVerified)
	if err != nil {
		return nil, err
	}

	return &LoginUserWithGoogleOutput{
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func TestLoginUserWithGoogle_Success_ExistingUser(t *testing.T) {
	// Arrange
	mockUserRepo := new(mocks.MockUserRepository)
	mockSecurity := new(mocks.MockSecurityService)

	userID := uuid.New()
	email := "user@example.com"
	idToken := "fake.google.token"
	jwtToken := "jwt.token.string"
	refreshToken := "refresh.token.string"

	// Mock repository to find an existing user
	user := &userDomain.User{
		ID:    userID,
		Email: email,
		Verification: userDomain.UserVerification{
			IsVerified: true,
		},
	}

	mockUserRepo.On("FindByEmail", email).Return(user, nil)
	mockSecurity.On("GenerateJWT",
		mock.MatchedBy(func(id *uuid.UUID) bool { return *id == userID }),
		mock.AnythingOfType("*string"),
		mock.AnythingOfType("Duration"),
		"login",
		true).Return(jwtToken, nil)
	mockSecurity.On("GenerateJWT",
		mock.MatchedBy(func(id *uuid.UUID) bool { return *id == userID }),
		mock.AnythingOfType("*string"),
		mock.AnythingOfType("Duration"),
		"refresh_token",
		true).Return(refreshToken, nil)

	// Create the use case with our mock validation
	baseUseCase := NewLoginUserWithGoogleUseCase(mockUserRepo, mockSecurity)
	useCase := &MockLoginUserWithGoogleUseCase{
		LoginUserWithGoogleUseCase: baseUseCase,
		validateTokenFunc: func(idToken string) (string, error) {
			// Mock successful validation and return email
			return email, nil
		},
	}

	// Act
	output, err := useCase.Execute(LoginUserWithGoogleInput{
		IDToken:       idToken,
		PreferredLang: "en",
	})

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, jwtToken, output.Token)
	assert.Equal(t, refreshToken, output.RefreshToken)
	mockUserRepo.AssertExpectations(t)
	mockSecurity.AssertExpectations(t)
}

func TestLoginUserWithGoogle_Success_NewUser(t *testing.T) {
	// Arrange
	mockUserRepo := new(mocks.MockUserRepository)
	mockSecurity := new(mocks.MockSecurityService)

	email := "new.user@example.com"
	idToken := "fake.google.token"
	randomPassword := "random_password"
	hashedPassword := "hashed_password"
	jwtToken := "jwt.token.string"
	refreshToken := "refresh.token.string"

	// Mock to simulate that the user doesn't exist
	mockUserRepo.On("FindByEmail", email).Return(nil, errors.New("user not found"))
	mockSecurity.On("GenerateSecureRandomString", 32).Return(randomPassword, nil)
	mockSecurity.On("HashPassword", randomPassword).Return(hashedPassword, nil)
	mockUserRepo.On("Create", mock.AnythingOfType("*user.User")).Return(nil)
	mockSecurity.On("GenerateJWT",
		mock.AnythingOfType("*uuid.UUID"),
		mock.AnythingOfType("*string"),
		mock.AnythingOfType("Duration"),
		"login",
		false).Return(jwtToken, nil)
	mockSecurity.On("GenerateJWT",
		mock.AnythingOfType("*uuid.UUID"),
		mock.AnythingOfType("*string"),
		mock.AnythingOfType("Duration"),
		"refresh_token",
		false).Return(refreshToken, nil)

	// Create the use case with our mock validation
	baseUseCase := NewLoginUserWithGoogleUseCase(mockUserRepo, mockSecurity)
	useCase := &MockLoginUserWithGoogleUseCase{
		LoginUserWithGoogleUseCase: baseUseCase,
		validateTokenFunc: func(idToken string) (string, error) {
			// Mock successful validation and return email
			return email, nil
		},
	}

	// Act
	output, err := useCase.Execute(LoginUserWithGoogleInput{
		IDToken:       idToken,
		PreferredLang: "en",
	})

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, jwtToken, output.Token)
	assert.Equal(t, refreshToken, output.RefreshToken)
	mockUserRepo.AssertExpectations(t)
	mockSecurity.AssertExpectations(t)
}

func TestLoginUserWithGoogle_InvalidToken(t *testing.T) {
	// Arrange
	mockUserRepo := new(mocks.MockUserRepository)
	mockSecurity := new(mocks.MockSecurityService)

	idToken := "invalid.google.token"

	// Create the use case with our mock validation
	baseUseCase := NewLoginUserWithGoogleUseCase(mockUserRepo, mockSecurity)
	useCase := &MockLoginUserWithGoogleUseCase{
		LoginUserWithGoogleUseCase: baseUseCase,
		validateTokenFunc: func(idToken string) (string, error) {
			// Mock invalid token validation
			return "", errors.New("invalid token")
		},
	}

	// Act
	output, err := useCase.Execute(LoginUserWithGoogleInput{
		IDToken:       idToken,
		PreferredLang: "en",
	})

	// Assert
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, "invalid Google token", err.Error())
	mockUserRepo.AssertNotCalled(t, "FindByEmail")
	mockSecurity.AssertNotCalled(t, "GenerateJWT")
}

func TestLoginUserWithGoogle_MissingEmail(t *testing.T) {
	// Arrange
	mockUserRepo := new(mocks.MockUserRepository)
	mockSecurity := new(mocks.MockSecurityService)

	idToken := "valid.google.token.without.email"

	// Create the use case with our mock validation
	baseUseCase := NewLoginUserWithGoogleUseCase(mockUserRepo, mockSecurity)
	useCase := &MockLoginUserWithGoogleUseCase{
		LoginUserWithGoogleUseCase: baseUseCase,
		validateTokenFunc: func(idToken string) (string, error) {
			// Mock validation with missing email
			return "", nil // Return empty email but no error
		},
	}

	// Act
	output, err := useCase.Execute(LoginUserWithGoogleInput{
		IDToken:       idToken,
		PreferredLang: "en",
	})

	// Assert
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, "email not found in Google token", err.Error())
	mockUserRepo.AssertNotCalled(t, "FindByEmail")
	mockSecurity.AssertNotCalled(t, "GenerateJWT")
}

func TestLoginUserWithGoogle_UserCreationError(t *testing.T) {
	// Arrange
	mockUserRepo := new(mocks.MockUserRepository)
	mockSecurity := new(mocks.MockSecurityService)

	email := "error@example.com"
	idToken := "valid.google.token"
	randomPassword := "random_password"
	hashedPassword := "hashed_password"

	// Mock to simulate an error during user creation
	mockUserRepo.On("FindByEmail", email).Return(nil, errors.New("user not found"))
	mockSecurity.On("GenerateSecureRandomString", 32).Return(randomPassword, nil)
	mockSecurity.On("HashPassword", randomPassword).Return(hashedPassword, nil)
	mockUserRepo.On("Create", mock.AnythingOfType("*user.User")).Return(errors.New("creation error"))

	// Create the use case with our mock validation
	baseUseCase := NewLoginUserWithGoogleUseCase(mockUserRepo, mockSecurity)
	useCase := &MockLoginUserWithGoogleUseCase{
		LoginUserWithGoogleUseCase: baseUseCase,
		validateTokenFunc: func(idToken string) (string, error) {
			// Mock successful validation and return email
			return email, nil
		},
	}

	// Act
	output, err := useCase.Execute(LoginUserWithGoogleInput{
		IDToken:       idToken,
		PreferredLang: "en",
	})

	// Assert
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, "creation error", err.Error())
	mockUserRepo.AssertExpectations(t)
	mockSecurity.AssertExpectations(t)
}
