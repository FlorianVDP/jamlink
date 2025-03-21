package userUseCase

import (
	"errors"
	"tindermals-backend/internal/modules/user/domain"
	"tindermals-backend/internal/shared/security"
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
		return nil, errors.New("email already exists") // TODO: Ajouter une erreur personnalisée dans un fichier d'erreurs
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
