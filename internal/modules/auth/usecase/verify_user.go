package userUseCase

import (
	userDomain "jamlink-backend/internal/modules/auth/domain/user"
	"jamlink-backend/internal/shared/security"
	"time"
)

type VerifyUserUseCase struct {
	repo     userDomain.UserRepository
	security security.SecurityService
}

type VerifyUserInput struct {
	Token string `json:"token" binding:"required" example:"token"`
}

func NewVerifyUserUseCase(repo userDomain.UserRepository, security security.SecurityService) *VerifyUserUseCase {
	return &VerifyUserUseCase{repo: repo, security: security}
}

func (uc *VerifyUserUseCase) Execute(input VerifyUserInput) error {
	claims, err := uc.security.ValidateJWT(input.Token)

	if err != nil {
		return err
	}

	mail, ok := claims["email"].(string)

	if !ok {
		return security.ErrInvalidUserEmail
	}

	user, err := uc.repo.FindByEmail(mail)

	if err != nil {
		return err
	}

	user.Verification.IsVerified = true
	now := time.Now()
	user.Verification.VerifiedAt = &now
	err = uc.repo.Update(user)

	if err != nil {
		return err
	}

	return nil
}
