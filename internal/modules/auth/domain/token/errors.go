package token

import "errors"

var (
	ErrPasswordDoesntMatch = errors.New("password does not match")
	ErrTokenType           = errors.New("unexpected token type")
	ErrTokenExpired        = errors.New("token expired")
	ErrTokenCreationFailed = errors.New("token creation failed")
	ErrTokenDeletionFailed = errors.New("token deletion failed")
	ErrTokenNotFound       = errors.New("token not found")
)
