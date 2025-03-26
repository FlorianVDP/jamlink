package userUseCase

import (
	tokenDomain "jamlink-backend/internal/modules/auth/domain/token"
	userDomain "jamlink-backend/internal/modules/auth/domain/user"
	"jamlink-backend/internal/shared/security"
	"time"
)

type RefreshTokenInput struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
	//DeviceInfo   string `json:"device_info" binding:"required"`
}

type RefreshTokenOutput struct {
	Token        string `json:"token" example:"eyJhbGciOi..."`
	RefreshToken string `json:"-"`
}

type RefreshTokenUseCase struct {
	security  security.SecurityService
	userRepo  userDomain.UserRepository
	tokenRepo tokenDomain.TokenRepository
}

func NewRefreshTokenUseCase(security security.SecurityService, userRepo userDomain.UserRepository, tokenRepo tokenDomain.TokenRepository) *RefreshTokenUseCase {
	return &RefreshTokenUseCase{
		security, userRepo, tokenRepo,
	}
}

func (uc *RefreshTokenUseCase) Execute(input RefreshTokenInput) (*RefreshTokenOutput, error) {
	existingToken, err := uc.tokenRepo.FindByToken(input.RefreshToken)
	if err != nil || existingToken == nil || existingToken.ExpiresAt.Before(time.Now()) {
		return nil, tokenDomain.ErrTokenExpired
	}

	userId, err := uc.security.GetJWTInfo(input.RefreshToken)
	if err != nil {
		return nil, err
	}

	user, err := uc.userRepo.FindByID(userId)

	if err != nil {
		return nil, err
	}

	token, err := uc.security.GenerateJWT(&userId, nil, time.Minute*15, "login", user.Verification.IsVerified)
	if err != nil {
		return nil, err
	}

	const expiringTime = time.Hour * 24 * 7
	refreshToken, err := uc.security.GenerateJWT(&userId, nil, expiringTime, "refresh_token", user.Verification.IsVerified)
	if err != nil {
		return nil, err
	}

	inDBToken, err := tokenDomain.CreateToken(userId, refreshToken, time.Now().Add(expiringTime))
	if err != nil {
		return nil, err
	}
	err = uc.tokenRepo.Create(inDBToken)
	if err != nil {
		return nil, tokenDomain.ErrTokenCreationFailed
	}

	err = uc.tokenRepo.DeleteByID(existingToken.ID)
	if err != nil {
		return nil, tokenDomain.ErrTokenDeletionFailed
	}

	return &RefreshTokenOutput{Token: token, RefreshToken: refreshToken}, nil
}
