package userDomain

import (
	"time"

	"github.com/google/uuid"
	"tindermals-backend/internal/modules/user/domain/invariants"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func CreateUser(email string, password string) (*User, error) {
	if err := userInvariants.ValidateUser(email, password); err != nil {
		return nil, err
	}

	return &User{
		ID:        uuid.New(),
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
