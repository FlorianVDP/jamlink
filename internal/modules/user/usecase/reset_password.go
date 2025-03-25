package userUseCase

import (
	"errors"
	userDomain "jamlink-backend/internal/modules/user/domain"
	userInvariants "jamlink-backend/internal/modules/user/domain/invariants"
	"jamlink-backend/internal/shared/security"
	"time"
)

var (
	ErrPasswordDoesntMatch = errors.New("password does not match")
	ErrTokenType           = errors.New("unexpected token type")
	ErrTokenExpired        = errors.New("token expired")
)

type ResetPasswordUseCase struct {
	tokenRepo userDomain.TokenRepository
	userRepo  userDomain.UserRepository
	security  security.SecurityService
}

type ResetPasswordInput struct {
	Token                 string `json:"token" binding:"required"`
	NewPassword           string `json:"new_password" binding:"required"`
	NewPasswordValidation string `json:"new_password_validation" binding:"required"`
}

func NewResetPasswordUseCase(tokenRepo userDomain.TokenRepository, userRepo userDomain.UserRepository, security security.SecurityService) *ResetPasswordUseCase {
	return &ResetPasswordUseCase{tokenRepo, userRepo, security}
}

func (uc *ResetPasswordUseCase) Execute(input ResetPasswordInput) error {
	if input.NewPasswordValidation != input.NewPassword {
		return ErrPasswordDoesntMatch
	}

	if err := userInvariants.ValidatePassword(input.NewPassword); err != nil {
		return err
	}

	claims, err := uc.security.ValidateJWT(input.Token)
	if err != nil {
		return err
	}

	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "reset_password" {
		return ErrTokenType
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return ErrTokenExpired
	}
	if time.Now().After(time.Unix(int64(exp), 0)) {
		return ErrTokenExpired
	}

	email, ok := claims["email"].(string)
	if !ok {
		return security.ErrInvalidUserEmail
	}

	token, err := uc.tokenRepo.FindByToken(input.Token)
	if err != nil {
		return err
	}

	user, err := uc.userRepo.FindByEmail(email)
	if err != nil {
		return err
	}

	hashedPassword, err := uc.security.HashPassword(input.NewPassword)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	if err := uc.userRepo.Update(user); err != nil {
		return err
	}

	if err := uc.tokenRepo.DeleteByID(token.ID); err != nil {
		return err
	}

	return nil
}
