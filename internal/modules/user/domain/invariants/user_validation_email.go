package userInvariants

import "errors"

var (
	ErrInvalidUserEmail = errors.New("invalid user email")
)

func ValidateEmail(email string) error {
	if email == "" {
		return ErrInvalidUserEmail
	}

	return nil
}
