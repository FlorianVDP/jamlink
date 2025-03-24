package userUseCase

import (
	"context"
	"errors"
	"google.golang.org/api/idtoken"
	userDomain "jamlink-backend/internal/modules/user/domain"
	"jamlink-backend/internal/shared/security"
	"os"
)

type LoginUserWithGoogleInput struct {
	IDToken       string `json:"id_token" binding:"required"`
	PreferredLang string `gorm:"type:varchar(5);default:'en'" json:"-"`
}

type LoginUserWithGoogleOutput struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginUserWithGoogleUseCase struct {
	repo     userDomain.UserRepository
	security security.SecurityService
}

func NewLoginUserWithGoogleUseCase(repo userDomain.UserRepository, security security.SecurityService) *LoginUserWithGoogleUseCase {
	return &LoginUserWithGoogleUseCase{
		repo:     repo,
		security: security,
	}
}

func (uc *LoginUserWithGoogleUseCase) Execute(input LoginUserWithGoogleInput) (*LoginUserWithGoogleOutput, error) {
	payload, err := idtoken.Validate(context.Background(), input.IDToken, os.Getenv("GOOGLE_CLIENT_ID"))

	if err != nil {
		return nil, errors.New("invalid Google token")
	}

	email, ok := payload.Claims["email"].(string)

	if !ok {
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

	token, err := uc.security.GenerateJWT(user.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := uc.security.GenerateRefreshJWT(user.ID)
	if err != nil {
		return nil, err
	}

	return &LoginUserWithGoogleOutput{
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}
