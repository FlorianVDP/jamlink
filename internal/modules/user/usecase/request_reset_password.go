package userUseCase

import (
	"fmt"
	userDomain "jamlink-backend/internal/modules/user/domain"
	"jamlink-backend/internal/shared/email"
	"jamlink-backend/internal/shared/security"
	"os"
	"time"
)

type RequestResetPasswordUseCase struct {
	tokenRepo    userDomain.TokenRepository
	userRepo     userDomain.UserRepository
	security     security.SecurityService
	emailService email.EmailService
}

type RequestResetPasswordInput struct {
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
}

func NewRequestResetPasswordUseCase(tokenRepo userDomain.TokenRepository, userRepo userDomain.UserRepository, security security.SecurityService, emailService email.EmailService) *RequestResetPasswordUseCase {
	return &RequestResetPasswordUseCase{tokenRepo: tokenRepo, userRepo: userRepo, security: security, emailService: emailService}
}

func (uc *RequestResetPasswordUseCase) Execute(input RequestResetPasswordInput) error {
	user, err := uc.userRepo.FindByEmail(input.Email)
	if err != nil {
		return err
	}

	jwt, err := uc.security.GenerateJWT(&user.ID, &user.Email, time.Minute*15, "reset_password")
	if err != nil {
		return err
	}

	token := &userDomain.Token{
		UserID:    user.ID,
		Token:     jwt,
		ExpiresAt: time.Now().Add(time.Minute * 15),
	}

	err = uc.tokenRepo.Create(token)

	if err != nil {
		return err
	}

	return uc.emailService.Send(user.Email, email.TemplateResetPassword, user.PreferredLang, map[string]string{
		"URL": fmt.Sprintf("%s?token=%s", os.Getenv("FRONTEND_VERIFY_URL"), token),
	})
}
