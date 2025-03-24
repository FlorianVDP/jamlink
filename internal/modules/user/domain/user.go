package userDomain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID        `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Email         string           `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password      string           `gorm:"type:varchar(255);not null" json:"-"`
	CreatedAt     time.Time        `gorm:"autoCreateTime" json:"-"`
	UpdatedAt     time.Time        `gorm:"autoUpdateTime" json:"-"`
	PreferredLang string           `gorm:"type:varchar(5);default:'fr'" json:"preferredLang"`
	Verification  UserVerification `gorm:"embedded" json:"-"`
	Provider      string           `gorm:"default:'local'" json:"-"`
}

type UserVerification struct {
	IsVerified bool       `gorm:"boolean;default:false" json:"-"`
	VerifiedAt *time.Time `gorm:"autoUpdateTime;default:null" json:"-"`
}

func CreateUser(email string, password string, preferredLang string, provider string) (*User, error) {
	return &User{
		ID:            uuid.New(),
		Email:         email,
		Password:      password,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		PreferredLang: preferredLang,
		Verification: UserVerification{
			IsVerified: false,
			VerifiedAt: nil,
		},
		Provider: provider,
	}, nil
}
