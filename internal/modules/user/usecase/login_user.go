package userUseCase

import (
	"errors"
	userDomain "jamlink-backend/internal/modules/user/domain"
	"jamlink-backend/internal/shared/security"
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
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
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

	token, err := uc.security.GenerateJWT(user.ID)

	if err != nil {
		return nil, err
	}

	refreshToken, err := uc.security.GenerateRefreshJWT(user.ID)
	if err != nil {
		return nil, err
	}

	return &LoginUserOutput{Token: token, RefreshToken: refreshToken}, nil
}
