package userUseCase

import (
	"jamlink-backend/internal/modules/user/domain"
	userInvariants "jamlink-backend/internal/modules/user/domain/invariants"
	"jamlink-backend/internal/shared/security"
)

type CreateUserUseCase struct {
	repo     userDomain.UserRepository
	security security.SecurityService
}

func NewCreateUserUseCase(repo userDomain.UserRepository, security security.SecurityService) *CreateUserUseCase {
	return &CreateUserUseCase{repo: repo, security: security}
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

	user, err := userDomain.CreateUser(input.Email, hashedPassword, input.PreferredLang, "local")
	if err != nil {
		return nil, err
	}

	err = uc.repo.Create(user)

	if err != nil {
		return nil, err
	}

	return user, err
}
