package animalInvariants

import "errors"

var (
	ErrInvalidAnimalAge = errors.New("invalid animal age")
)

func ValidateAge(age int) error {
	if age <= 0 {
		return ErrInvalidAnimalAge
	}

	return nil
}
