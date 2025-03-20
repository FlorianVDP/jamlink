package test

import (
	"testing"
	"tindermals-backend/internal/domain"
	"tindermals-backend/internal/usecase"
	"tindermals-backend/test/mocks"
)

func TestGetAnimalList(t *testing.T) {
	mockRepo := &mocks.MockAnimalRepository{
		Animals: []*domain.Animal{
			{ID: "1", Name: "Muxu", Age: 2, Sexe: "ale", Description: "Sweety cat", Image: "test"},
			{ID: "2", Name: "Léa", Age: 3, Sexe: "Female", Description: "Sweety girl", Image: "test"},
		},
	}

	uc := usecase.NewGetAnimalListUseCase(mockRepo)

	animalList, err := uc.Execute()

	if err != nil {
		t.Fatalf("Erreur lors de l'exécution du cas d'utilisation: %v", err)
	}

	if len(animalList) != 2 {
		t.Fatalf("Attendu 2 animaux, mais obtenu %d", len(animalList))
	}

	if animalList[0].ID != "1" || animalList[0].Name != "Muxu" {
		t.Fatalf("Le premier animal devrait être Muxu, mais a été %v", animalList[0])
	}

	if animalList[1].ID != "2" || animalList[1].Name != "Léa" {
		t.Fatalf("Le deuxième animal devrait être Léa, mais a été %v", animalList[1])
	}
}
