package userUseCase

import (
	"fmt"
	userDomain "jamlink-backend/internal/modules/user/domain"
	"jamlink-backend/internal/shared/email"
	"jamlink-backend/internal/shared/security"
	"os"
)

type GetVerificationEmailUseCase struct {
	security security.SecurityService
	email    email.EmailService
	repo     userDomain.UserRepository
}

type GetVerificationEmailInput struct {
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
}

func NewGetVerificationEmailUseCase(security security.SecurityService, repo userDomain.UserRepository, email email.EmailService) *GetVerificationEmailUseCase {
	return &GetVerificationEmailUseCase{security: security, repo: repo, email: email}
}

func (uc *GetVerificationEmailUseCase) Execute(input GetVerificationEmailInput) error {
	user, err := uc.repo.FindByEmail(input.Email)

	if err != nil {
		return userDomain.ErrUserNotFound
	}

	token, err := uc.security.GenerateVerificationJWT(input.Email)
	if err != nil {
		return err
	}

	err = uc.email.Send(user.Email, email.TemplateVerification, user.PreferredLang, map[string]string{
		"URL": fmt.Sprintf("%s?token=%s", os.Getenv("FRONTEND_VERIFY_URL"), token),
	})

	if err != nil {
		return err
	}

	return nil
}
