package test

import (
	"testing"
	"tindermals-backend/internal/infra"
	"tindermals-backend/internal/usecase"
)

func TestCreateAnimal(t *testing.T) {
	repo := infra.NewMemoryAnimalRepository()
	uc := usecase.NewCreateAnimalUseCase(repo)

	input := usecase.CreateAnimalInput{
		Name:        "Muxu",
		Age:         2,
		Sexe:        "Male",
		Description: "Sweety cat",
		Image:       "test",
	}

	animal, err := uc.Execute(input)

	if err != nil {
		t.Fatalf("Erreur : %v", err)
	}

	if animal.Name != "Muxu" {
		t.Errorf("Incorrect name: expected 'Muxu', received '%s'", animal.Name)
	}
}
