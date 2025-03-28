package useCase

import (
	"errors"
	tokenDomain "jamlink-backend/internal/modules/auth/domain/token"
	userDomain "jamlink-backend/internal/modules/auth/domain/user"
	"jamlink-backend/internal/shared/security"
	"time"
)

var (
	ErrInvalidEmailOrPassword = errors.New("invalid email or password")
)

type LoginUserUseCase struct {
	userRepo  userDomain.UserRepository
	security  security.SecurityService
	tokenRepo tokenDomain.TokenRepository
}

func NewLoginUserUseCase(userRepo userDomain.UserRepository, security security.SecurityService, tokenRepo tokenDomain.TokenRepository) *LoginUserUseCase {
	return &LoginUserUseCase{
		userRepo,
		security,
		tokenRepo,
	}
}

type LoginUserInput struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"Abcd1234!"`
	//DeviceInfo   string `json:"device_info" binding:"required"`
}

type LoginUserOutput struct {
	Token        string `json:"token"`
	RefreshToken string `json:"-"`
}

func (uc *LoginUserUseCase) Execute(input LoginUserInput) (*LoginUserOutput, error) {
	user, err := uc.userRepo.FindByEmail(input.Email)
	if err != nil {
		return nil, ErrInvalidEmailOrPassword
	}

	if !uc.security.CheckPassword(input.Password, user.Password) {
		return nil, security.ErrPasswordComparison
	}

	const tokenExpiringTime = time.Minute * 15

	token, err := uc.security.GenerateJWT(&user.ID, nil, tokenExpiringTime, "login", user.Verification.IsVerified)

	if err != nil {
		return nil, err
	}
	inDBToken, err := tokenDomain.CreateToken(user.ID, token, time.Now().Add(tokenExpiringTime))
	if err != nil {
		return nil, err
	}
	err = uc.tokenRepo.Create(inDBToken)
	if err != nil {
		return nil, err
	}

	const refreshTokenExpiringTime = time.Hour * 24 * 7
	refreshToken, err := uc.security.GenerateJWT(&user.ID, nil, refreshTokenExpiringTime, "refresh_token", user.Verification.IsVerified)
	if err != nil {
		return nil, err
	}

	inDBREfreshToken, err := tokenDomain.CreateToken(user.ID, refreshToken, time.Now().Add(refreshTokenExpiringTime))
	if err != nil {
		return nil, err
	}
	err = uc.tokenRepo.Create(inDBREfreshToken)
	if err != nil {
		return nil, err
	}

	return &LoginUserOutput{Token: token, RefreshToken: refreshToken}, nil
}
