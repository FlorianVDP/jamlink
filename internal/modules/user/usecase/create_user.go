package userUseCase

import (
	"fmt"
	"jamlink-backend/internal/modules/user/domain"
	userInvariants "jamlink-backend/internal/modules/user/domain/invariants"
	"jamlink-backend/internal/shared/email"
	"jamlink-backend/internal/shared/security"
	"os"
)

type CreateUserUseCase struct {
	repo     userDomain.UserRepository
	security security.SecurityService
	email    email.EmailService
}

func NewCreateUserUseCase(repo userDomain.UserRepository, security security.SecurityService, email email.EmailService) *CreateUserUseCase {
	return &CreateUserUseCase{repo: repo, security: security, email: email}
}

type CreateUserInput struct {
	Email         string `json:"email" binding:"required,email" example:"user@example.com"`
	Password      string `json:"password" binding:"required" example:"Abcd1234!"`
	PreferredLang string `gorm:"type:varchar(5);default:'en'" json:"-"`
}

func (uc *CreateUserUseCase) Execute(input CreateUserInput) (*userDomain.User, error) {
	_, err := uc.repo.FindByEmail(input.Email)

	if err == nil {
		return nil, userDomain.ErrEmailAlreadyExists
	}

	if err := userInvariants.ValidateUser(input.Email, input.Password); err != nil {
		return nil, err
	}

	hashedPassword, err := uc.security.HashPassword(input.Password)

	if err != nil {
		return nil, err
	}

	user, err := userDomain.CreateUser(input.Email, hashedPassword, input.PreferredLang)
	if err != nil {
		return nil, err
	}

	err = uc.repo.Create(user)

	if err != nil {
		return nil, err
	}

	token, err := uc.security.GenerateVerificationJWT(user.Email)
	if err != nil {
		return nil, err
	}

	err = uc.email.Send(user.Email, email.TemplateVerification, user.PreferredLang, map[string]string{
		"URL": fmt.Sprintf("%s?token=%s", os.Getenv("FRONTEND_VERIFY_URL"), token),
	})

	if err != nil {
		return nil, err
	}

	return user, err
}
