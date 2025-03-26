package userUseCase

import (
	"fmt"
	"jamlink-backend/internal/modules/auth/domain/token"
	"jamlink-backend/internal/modules/auth/domain/user"
	"jamlink-backend/internal/shared/email"
	"jamlink-backend/internal/shared/security"
	"os"
	"time"
)

type RequestResetPasswordUseCase struct {
	tokenRepo    token.TokenRepository
	userRepo     user.UserRepository
	security     security.SecurityService
	emailService email.EmailService
}

type RequestResetPasswordInput struct {
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
}

func NewRequestResetPasswordUseCase(tokenRepo token.TokenRepository, userRepo user.UserRepository, security security.SecurityService, emailService email.EmailService) *RequestResetPasswordUseCase {
	return &RequestResetPasswordUseCase{tokenRepo: tokenRepo, userRepo: userRepo, security: security, emailService: emailService}
}

func (uc *RequestResetPasswordUseCase) Execute(input RequestResetPasswordInput) error {
	foundUser, err := uc.userRepo.FindByEmail(input.Email)
	if err != nil {
		return err
	}

	jwt, err := uc.security.GenerateJWT(&foundUser.ID, &foundUser.Email, time.Minute*15, "reset_password", foundUser.Verification.IsVerified)
	if err != nil {
		return err
	}

	createdToken, err := token.CreateToken(foundUser.ID, jwt, time.Now().Add(time.Minute*15))

	if err != nil {
		return err
	}

	err = uc.tokenRepo.Create(createdToken)

	if err != nil {
		return err
	}

	return uc.emailService.Send(foundUser.Email, email.TemplateResetPassword, foundUser.PreferredLang, map[string]string{
		"URL": fmt.Sprintf("%s?token=%s", os.Getenv("FRONTEND_VERIFY_URL"), createdToken),
	})
}
