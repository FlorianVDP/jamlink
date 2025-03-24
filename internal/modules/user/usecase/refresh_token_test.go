package userUseCase

import (
	"jamlink-backend/internal/modules/user/mocks"
	"jamlink-backend/internal/shared/security"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRefreshToken_Success(t *testing.T) {
	mockSecurity := new(mocks.MockSecurityService)

	refreshToken := "valid_refresh_token"
	userID := uuid.New()
	newAccessToken := "new.jwt.token"
	newRefreshToken := "new.refresh.token"

	mockSecurity.On("GetJWTInfo", refreshToken).Return(userID, nil)
	mockSecurity.On("GenerateJWT", userID).Return(newAccessToken, nil)
	mockSecurity.On("GenerateRefreshJWT", userID).Return(newRefreshToken, nil)

	usecase := NewRefreshTokenUseCase(mockSecurity)

	output, err := usecase.Execute(RefreshTokenInput{RefreshToken: refreshToken})

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, newAccessToken, output.Token)
	assert.Equal(t, newRefreshToken, output.RefreshToken)
}

func TestRefreshToken_InvalidToken(t *testing.T) {
	mockSecurity := new(mocks.MockSecurityService)

	refreshToken := "invalid_token"

	mockSecurity.On("GetJWTInfo", refreshToken).Return(uuid.Nil, security.ErrInvalidToken)

	usecase := NewRefreshTokenUseCase(mockSecurity)

	output, err := usecase.Execute(RefreshTokenInput{RefreshToken: refreshToken})

	assert.ErrorIs(t, err, security.ErrInvalidToken)
	assert.Nil(t, output)
}

func TestRefreshToken_JWTGenerationFails(t *testing.T) {
	mockSecurity := new(mocks.MockSecurityService)

	refreshToken := "valid_token"
	userID := uuid.New()

	mockSecurity.On("GetJWTInfo", refreshToken).Return(userID, nil)
	mockSecurity.On("GenerateJWT", userID).Return("", security.ErrJWTGeneration)

	usecase := NewRefreshTokenUseCase(mockSecurity)

	output, err := usecase.Execute(RefreshTokenInput{RefreshToken: refreshToken})

	assert.ErrorIs(t, err, security.ErrJWTGeneration)
	assert.Nil(t, output)
}

func TestRefreshToken_RefreshJWTGenerationFails(t *testing.T) {
	mockSecurity := new(mocks.MockSecurityService)

	refreshToken := "valid_token"
	userID := uuid.New()
	accessToken := "new.jwt.token"

	mockSecurity.On("GetJWTInfo", refreshToken).Return(userID, nil)
	mockSecurity.On("GenerateJWT", userID).Return(accessToken, nil)
	mockSecurity.On("GenerateRefreshJWT", userID).Return("", security.ErrRefreshJWTGeneration)

	usecase := NewRefreshTokenUseCase(mockSecurity)

	output, err := usecase.Execute(RefreshTokenInput{RefreshToken: refreshToken})

	assert.ErrorIs(t, err, security.ErrRefreshJWTGeneration)
	assert.Nil(t, output)
}
