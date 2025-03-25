package token

import (
	"github.com/google/uuid"
	tokenInvariants "jamlink-backend/internal/modules/auth/domain/token/invariants"
	"time"
)

type Token struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Token     string    `gorm:"type:text;not null;uniqueIndex"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func CreateToken(userID uuid.UUID, token string, expiresAt time.Time) (*Token, error) {
	if err := tokenInvariants.ValidateToken(expiresAt); err != nil {
		return nil, err
	}

	return &Token{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}, nil
}
