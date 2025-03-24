package userUseCase

import (
	"jamlink-backend/internal/modules/user/domain"
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
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (uc *CreateUserUseCase) Execute(input CreateUserInput) (*userDomain.User, error) {
	_, err := uc.repo.FindByEmail(input.Email)

	if err == nil {
		return nil, userDomain.ErrEmailAlreadyExists
	}

	hashedPassword, err := uc.security.HashPassword(input.Password)

	if err != nil {
		return nil, err
	}

	user, err := userDomain.CreateUser(input.Email, hashedPassword)

	if err != nil {
		return nil, err
	}

	err = uc.repo.Create(user)
	return user, err
}
