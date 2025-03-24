package security

import "errors"

var (
	// Password
	ErrPasswordHashing    = errors.New("failed to hash password")
	ErrPasswordComparison = errors.New("password does not match")

	// JWT Generation
	ErrJWTGeneration        = errors.New("failed to generate JWT")
	ErrRefreshJWTGeneration = errors.New("failed to generate refresh JWT")

	// JWT Validation
	ErrInvalidJWTSigningMethod = errors.New("unexpected JWT signing method")
	ErrInvalidToken            = errors.New("invalid JWT token")
	ErrCannotExtractClaims     = errors.New("unable to extract claims from token")

	// JWT Claims parsing
	ErrInvalidUserID = errors.New("invalid or missing user ID in JWT claims")
)
