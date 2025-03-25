package userUseCase

import (
	"errors"
	userDomain "jamlink-backend/internal/modules/user/domain"
	"jamlink-backend/internal/shared/security"
	"time"
)

var (
	ErrInvalidEmailOrPassword = errors.New("invalid email or password")
)

type LoginUserUseCase struct {
	repo     userDomain.UserRepository
	security security.SecurityService
}

func NewLoginUserUseCase(repo userDomain.UserRepository, security security.SecurityService) *LoginUserUseCase {
	return &LoginUserUseCase{
		repo:     repo,
		security: security,
	}
}

type LoginUserInput struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"Abcd1234!"`
}

type LoginUserOutput struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func (uc *LoginUserUseCase) Execute(input LoginUserInput) (*LoginUserOutput, error) {
	user, err := uc.repo.FindByEmail(input.Email)
	if err != nil {
		return nil, ErrInvalidEmailOrPassword
	}

	if !uc.security.CheckPassword(input.Password, user.Password) {
		return nil, ErrInvalidEmailOrPassword
	}

	token, err := uc.security.GenerateJWT(&user.ID, nil, time.Minute*15, "login")

	if err != nil {
		return nil, err
	}

	refreshToken, err := uc.security.GenerateJWT(&user.ID, nil, time.Hour*24*7, "refresh_token")
	if err != nil {
		return nil, err
	}

	return &LoginUserOutput{Token: token, RefreshToken: refreshToken}, nil
}
