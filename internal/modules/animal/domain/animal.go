package animalDomain

import (
	"time"

	"github.com/google/uuid"
	animalInvariants "tindermals-backend/internal/modules/animal/domain/invariants"
)

type Animal struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Age         int       `gorm:"not null" json:"age"`
	Sexe        string    `gorm:"type:varchar(10);not null" json:"sexe"`
	Description string    `gorm:"type:text" json:"description,omitempty"`
	Image       string    `gorm:"type:varchar(255)" json:"image,omitempty"`
	Adopted     bool      `gorm:"default:false" json:"adopted"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func CreateAnimal(name string, age int, sexe string, description string, image string) (*Animal, error) {
	if err := animalInvariants.ValidateAnimal(age, name); err != nil {
		return nil, err
	}

	return &Animal{
		ID:          uuid.New(),
		Name:        name,
		Age:         age,
		Sexe:        sexe,
		Description: description,
		Image:       image,
		Adopted:     false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}
