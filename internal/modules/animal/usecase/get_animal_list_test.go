package animalUseCase

import (
	"github.com/google/uuid"
	"jamlink-backend/internal/modules/animal/domain"
	"jamlink-backend/internal/modules/animal/mocks"
	"testing"
)

func TestGetAnimalList(t *testing.T) {
	mockRepo := &animalMocks.MockAnimalRepository{
		Animals: []*animalDomain.Animal{
			{ID: uuid.New(), Name: "Muxu", Age: 2, Sexe: "Male", Description: "Sweety cat", Image: "test"},
			{ID: uuid.New(), Name: "Léa", Age: 3, Sexe: "Female", Description: "Sweety girl", Image: "test"},
		},
	}

	uc := NewGetAnimalListUseCase(mockRepo)

	animalList, err := uc.Execute()

	if err != nil {
		t.Fatalf("Error while executing the use case: %v", err)
	}

	if len(animalList) != 2 {
		t.Fatalf("Expected 2 animals, but got %d", len(animalList))
	}

	if animalList[0].Name != "Muxu" {
		t.Fatalf("The first animal should be Muxu, but got %v", animalList[0])
	}

	if animalList[1].Name != "Léa" {
		t.Fatalf("The second animal should be Léa, but got %v", animalList[1])
	}
}
