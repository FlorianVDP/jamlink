package userRepository

import (
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
	var user user.User

	if err := r.db.Where("email  = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepository) Update(user *user.User) error {
	return r.db.Save(user).Error
}
