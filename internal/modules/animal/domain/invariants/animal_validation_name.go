package animalInvariants

import "errors"

var (
	ErrInvalidAnimalName = errors.New("invalid animal name")
)

func ValidateName(name string) error {
	if name == "" {
		return ErrInvalidAnimalName
	}

	return nil
}
