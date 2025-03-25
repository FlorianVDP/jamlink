package userRepository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	userDomain "jamlink-backend/internal/modules/user/domain"
)

type PostgresTokenRepository struct {
	db *gorm.DB
}

func NewPostgresTokenRepository(db *gorm.DB) *PostgresTokenRepository {
	return &PostgresTokenRepository{db: db}
}

func (r *PostgresTokenRepository) Create(token *userDomain.Token) error {
	return r.db.Create(token).Error
}

func (r *PostgresTokenRepository) FindByToken(token string) (*userDomain.Token, error) {
	var t userDomain.Token

	if err := r.db.Where("token  = ?", token).First(&t).Error; err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *PostgresTokenRepository) DeleteByUserID(userID uuid.UUID) error {
	return r.db.Where("user_id = ?", userID).Delete(&userDomain.Token{}).Error
}
