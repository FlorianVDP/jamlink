package userUseCase

import (
	"jamlink-backend/internal/modules/auth/mocks"
	"jamlink-backend/internal/shared/security"
	"testing"
	"time"

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

	mockSecurity.On("GenerateJWT", &userID, (*string)(nil), time.Minute*15, "login").Return(newAccessToken, nil)
	mockSecurity.On("GenerateJWT", &userID, (*string)(nil), time.Hour*24*7, "refresh_token").Return(newRefreshToken, nil)

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
	mockSecurity.On("GenerateJWT", &userID, (*string)(nil), time.Minute*15, "login").Return("", security.ErrJWTGeneration)

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

	mockSecurity.On("GenerateJWT", &userID, (*string)(nil), time.Minute*15, "login").Return(accessToken, nil)
	mockSecurity.On("GenerateJWT", &userID, (*string)(nil), time.Hour*24*7, "refresh_token").Return("", security.ErrJWTGeneration)

	usecase := NewRefreshTokenUseCase(mockSecurity)

	output, err := usecase.Execute(RefreshTokenInput{RefreshToken: refreshToken})

	assert.ErrorIs(t, err, security.ErrJWTGeneration)
	assert.Nil(t, output)
}
