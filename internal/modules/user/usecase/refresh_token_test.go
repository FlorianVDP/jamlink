package userUseCase

import (
	"errors"
	"testing"
	"tindermals-backend/internal/modules/user/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRefreshToken_Success(t *testing.T) {
	mockSecurity := new(mocks.MockSecurityService)

	refreshToken := "valid_refresh_token"
	userID := uuid.New()
	newAccessToken := "new.jwt.token"

	mockSecurity.On("GetJWTInfo", refreshToken).Return(userID, nil)
	mockSecurity.On("GenerateJWT", userID).Return(newAccessToken, nil)

	usecase := NewRefreshTokenUseCase(mockSecurity)

	output, err := usecase.Execute(RefreshTokenInput{RefreshToken: refreshToken})

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, newAccessToken, output.Token)
}

func TestRefreshToken_InvalidToken(t *testing.T) {
	mockSecurity := new(mocks.MockSecurityService)

	refreshToken := "invalid_token"

	mockSecurity.On("GetJWTInfo", refreshToken).Return(uuid.Nil, errors.New("invalid refresh token"))

	usecase := NewRefreshTokenUseCase(mockSecurity)

	output, err := usecase.Execute(RefreshTokenInput{RefreshToken: refreshToken})

	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, "invalid refresh token", err.Error())
}

func TestRefreshToken_JWTGenerationFails(t *testing.T) {
	mockSecurity := new(mocks.MockSecurityService)

	refreshToken := "valid_token_but_jwt_fails"
	userID := uuid.New()

	mockSecurity.On("GetJWTInfo", refreshToken).Return(userID, nil)
	mockSecurity.On("GenerateJWT", userID).Return("", errors.New("JWT generation failed"))

	usecase := NewRefreshTokenUseCase(mockSecurity)

	output, err := usecase.Execute(RefreshTokenInput{RefreshToken: refreshToken})

	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, "JWT generation failed", err.Error())
}
