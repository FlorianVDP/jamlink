package userUseCase

import (
	"jamlink-backend/internal/shared/security"
	"time"
)

type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshTokenOutput struct {
	Token        string `json:"token" example:"eyJhbGciOi..."`
	RefreshToken string `json:"-"`
}

type RefreshTokenUseCase struct {
	security security.SecurityService
}

func NewRefreshTokenUseCase(security security.SecurityService) *RefreshTokenUseCase {
	return &RefreshTokenUseCase{
		security: security,
	}
}

func (uc *RefreshTokenUseCase) Execute(input RefreshTokenInput) (*RefreshTokenOutput, error) {
	userId, err := uc.security.GetJWTInfo(input.RefreshToken)

	if err != nil {
		return nil, err
	}

	token, err := uc.security.GenerateJWT(&userId, nil, time.Minute*15, "login")
	if err != nil {
		return nil, err
	}

	refreshToken, err := uc.security.GenerateJWT(&userId, nil, time.Hour*24*7, "refresh_token")
	if err != nil {
		return nil, err
	}

	return &RefreshTokenOutput{Token: token, RefreshToken: refreshToken}, nil
}
