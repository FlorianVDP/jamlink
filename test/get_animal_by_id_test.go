package test

import (
	"testing"
	"tindermals-backend/internal/domain"
	"tindermals-backend/internal/usecase"
	"tindermals-backend/test/mocks"
)

func TestGetAnimalById(t *testing.T) {
	mockRepo := &mocks.MockAnimalRepository{
		Animals: []*domain.Animal{
			{ID: "1", Name: "Muxu", Age: 2, Sexe: "male", Description: "Sweety cat", Image: "test"},
			{ID: "2", Name: "Léa", Age: 3, Sexe: "Female", Description: "Sweety girl", Image: "test"},
		},
	}

	uc := usecase.NewGetAnimalByIdUseCase(mockRepo)

	animal, err := uc.Execute("1")

	if err != nil {
		t.Fatalf("Erreur lors de l'exécution du cas d'utilisation: %v", err)
	}

	if animal.ID != "1" || animal.Name != "Muxu" {
		t.Fatalf("L'animal devrait être Muxu, mais a été %v", animal)
	}
}
