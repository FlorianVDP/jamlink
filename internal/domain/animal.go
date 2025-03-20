package domain

import (
	"errors"

	"github.com/google/uuid"
)

type Animal struct {
	ID          string
	Name        string
	Age         int
	Sexe        string
	Description string
	Image       string
	Adopted     bool
}

var (
	ErrInvalidAnimalName = errors.New("invalid animal name")
	ErrInvalidAnimalAge  = errors.New("invalid animal age")
)

func NewAnimal(name string, age int, sexe string, description string, image string) (*Animal, error) {
	if err := IsValidAnimal(name, age); err != nil {
		return nil, err
	}

	id := uuid.New().String()

	return &Animal{
		ID:          id,
		Name:        name,
		Age:         age,
		Sexe:        sexe,
		Description: description,
		Image:       image,
		Adopted:     false,
	}, nil
}

func IsValidAnimal(name string, age int) error {
	if name == "" {
		return ErrInvalidAnimalName
	}

	if age <= 0 {
		return ErrInvalidAnimalAge
	}

	return nil
}

func (t *Animal) IsAdopted() {
	t.Adopted = true
}
