package userRepository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"jamlink-backend/internal/modules/auth/domain/user"
)

type PostgresUserRepository struct {
	db *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Create(user *user.User) error {
	return r.db.Create(user).Error
}

func (r *PostgresUserRepository) FindByEmail(email string) (*user.User, error) {
	var foundUser user.User

	if err := r.db.Where("email  = ?", email).First(&foundUser).Error; err != nil {
		return nil, err
	}

	return &foundUser, nil
}

func (r *PostgresUserRepository) FindByID(id uuid.UUID) (*user.User, error) {
	var foundUser user.User

	if err := r.db.Where("id = ?", id.String()).First(&foundUser).Error; err != nil {
		return nil, err
	}

	return &foundUser, nil
}

func (r *PostgresUserRepository) Update(user *user.User) error {
	return r.db.Save(user).Error
}
