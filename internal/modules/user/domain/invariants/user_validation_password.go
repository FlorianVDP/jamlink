package userInvariants

import (
	"errors"
	"regexp"
)

var (
	ErrInvalidUserPassword = errors.New("invalid user password")
	ErrShortPassword       = errors.New("password must be at least 8 characters long")
	ErrLongPassword        = errors.New("password must be at most 64 characters long")
	ErrNoUpperCase         = errors.New("password must contain at least one uppercase letter")
	ErrNoLowerCase         = errors.New("password must contain at least one lowercase letter")
	ErrNoDigit             = errors.New("password must contain at least one digit")
	ErrNoSpecialChar       = errors.New("password must contain at least one special character")
)

var (
	upperCaseRegex   = regexp.MustCompile(`[A-Z]`)
	lowerCaseRegex   = regexp.MustCompile(`[a-z]`)
	digitRegex       = regexp.MustCompile(`[0-9]`)
	specialCharRegex = regexp.MustCompile(`[\W_]`)
)

func ValidatePassword(password string) error {
	if password == "" {
		return ErrInvalidUserPassword
	}

	if len(password) < 8 {
		return ErrShortPassword
	}

	if len(password) > 64 {
		return ErrLongPassword
	}

	if !upperCaseRegex.MatchString(password) {
		return ErrNoUpperCase
	}

	if !lowerCaseRegex.MatchString(password) {
		return ErrNoLowerCase
	}

	if !digitRegex.MatchString(password) {
		return ErrNoDigit
	}

	if !specialCharRegex.MatchString(password) {
		return ErrNoSpecialChar
	}

	return nil
}
