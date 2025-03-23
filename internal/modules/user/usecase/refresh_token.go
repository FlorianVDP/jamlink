package userUseCase

import "jamlink-backend/internal/shared/security"

type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshTokenOutput struct {
	Token string `json:"token"`
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

	token, err := uc.security.GenerateJWT(userId)
	if err != nil {
		return nil, err
	}

	return &RefreshTokenOutput{Token: token}, nil
}
