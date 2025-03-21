package animalUseCase

import (
	"testing"
	animal2 "tindermals-backend/internal/modules/animal/repository"
)

func TestCreateAnimal(t *testing.T) {
	repo := animal2.NewMemoryAnimalRepository()
	uc := NewCreateAnimalUseCase(repo)

	input := CreateAnimalInput{
		Name:        "Muxu",
		Age:         2,
		Sexe:        "Male",
		Description: "Sweety cat",
		Image:       "test",
	}

	animal, err := uc.Execute(input)

	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if animal.Name != "Muxu" {
		t.Errorf("Incorrect name: expected 'Muxu', received '%s'", animal.Name)
	}
}
