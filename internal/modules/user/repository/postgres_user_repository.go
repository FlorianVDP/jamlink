package userRepository

import (
	"gorm.io/gorm"
	"jamlink-backend/internal/modules/user/domain"
)

type PostgresUserRepository struct {
	db *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Create(user *userDomain.User) error {
	return r.db.Create(user).Error
}

func (r *PostgresUserRepository) FindByEmail(email string) (*userDomain.User, error) {
	var user userDomain.User

	if err := r.db.Where("email  = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepository) Update(user *userDomain.User) error {
	return r.db.Save(user).Error
}
