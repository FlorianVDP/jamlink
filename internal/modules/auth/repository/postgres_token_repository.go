package userRepository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	tokenDomain "jamlink-backend/internal/modules/auth/domain/token"
)

type PostgresTokenRepository struct {
	db *gorm.DB
}

func NewPostgresTokenRepository(db *gorm.DB) *PostgresTokenRepository {
	return &PostgresTokenRepository{db: db}
}

func (r *PostgresTokenRepository) Create(token *tokenDomain.Token) error {
	return r.db.Create(token).Error
}

func (r *PostgresTokenRepository) FindByToken(token string) (*tokenDomain.Token, error) {
	var t tokenDomain.Token

	if err := r.db.Where("token  = ?", token).First(&t).Error; err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *PostgresTokenRepository) DeleteByID(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&tokenDomain.Token{}).Error
}

func (r *PostgresTokenRepository) DeleteUserTokens(userID uuid.UUID) error {
	return r.db.Where("user_id = ?", userID).Delete(&tokenDomain.Token{}).Error
}
