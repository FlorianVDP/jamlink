package userInvariants

import "errors"

var (
	ErrInvalidUserPassword = errors.New("invalid user password")
	ErrShortPassword       = errors.New("password must be at least 8 characters long")
)

func ValidatePassword(password string) error {
	if password == "" {
		return ErrInvalidUserPassword
	}

	if len(password) < 8 {
		return ErrShortPassword
	}

	return nil
}
