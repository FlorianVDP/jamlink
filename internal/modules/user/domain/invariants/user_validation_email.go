package userInvariants

import (
	"errors"
	"regexp"
)

var (
	ErrInvalidUserEmail = errors.New("invalid user email")
	emailRegex          = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9\-]+(\.[a-zA-Z0-9\-]+)*\.[a-zA-Z]{2,}$`)
)

func ValidateEmail(email string) error {
	if email == "" || !emailRegex.MatchString(email) {
		return ErrInvalidUserEmail
	}

	return nil
}
