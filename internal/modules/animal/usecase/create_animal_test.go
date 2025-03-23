package animalUseCase

import (
	animal2 "jamlink-backend/internal/modules/animal/repository"
	"testing"
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
