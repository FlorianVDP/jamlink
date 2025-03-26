package token

import "github.com/google/uuid"

type TokenRepository interface {
	Create(token *Token) error
	FindByToken(token string) (*Token, error)
	DeleteByID(userID uuid.UUID) error
	DeleteUserTokens(userID uuid.UUID) error
}
