package useCase

import (
	"errors"
	"fmt"
	"jamlink-backend/internal/modules/auth/domain/user"
	"jamlink-backend/internal/shared/email"
	"jamlink-backend/internal/shared/security"
	"os"
	"time"
)

type RequestVerifyUserEmailUseCase struct {
	security security.SecurityService
	email    email.EmailService
	repo     user.UserRepository
}

type RequestVerifyUserEmailInput struct {
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
}

func NewRequestVerifyUserEmailUseCase(security security.SecurityService, repo user.UserRepository, email email.EmailService) *RequestVerifyUserEmailUseCase {
	return &RequestVerifyUserEmailUseCase{security: security, repo: repo, email: email}
}

func (uc *RequestVerifyUserEmailUseCase) Execute(input RequestVerifyUserEmailInput) error {
	foundUser, err := uc.repo.FindByEmail(input.Email)

	if err != nil {
		return user.ErrUserNotFound
	}

	token, err := uc.security.GenerateJWT(nil, &input.Email, time.Hour*24, "verify_email", foundUser.Verification.IsVerified)
	if err != nil {
		return err
	}

	if foundUser.Verification.IsVerified || foundUser.Verification.VerifiedAt != nil {
		return errors.New("your account is already verified")
	}

	err = uc.email.Send(foundUser.Email, email.TemplateVerification, foundUser.PreferredLang, map[string]string{
		"URL": fmt.Sprintf("%s?token=%s", os.Getenv("FRONTEND_VERIFY_URL"), token),
	})

	if err != nil {
		return err
	}

	return nil
}
