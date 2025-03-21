package animalUsecase

import (
	"github.com/google/uuid"
	"testing"
	"tindermals-backend/internal/modules/animal/domain"
	"tindermals-backend/internal/modules/animal/mocks"
)

func TestGetAnimalById(t *testing.T) {
	mockRepo := &animalMocks.MockAnimalRepository{
		Animals: []*animalDomain.Animal{
			{ID: uuid.New(), Name: "Muxu", Age: 2, Sexe: "Male", Description: "Sweety cat", Image: "test"},
			{ID: uuid.New(), Name: "Léa", Age: 3, Sexe: "Female", Description: "Sweety girl", Image: "test"},
		},
	}

	uc := NewGetAnimalByIdUseCase(mockRepo)

	animal, err := uc.Execute(mockRepo.Animals[0].ID.String())

	if err != nil {
		t.Fatalf("Error while executing the use case: %v", err)
	}

	if animal.Name != "Muxu" {
		t.Fatalf("The expected animal was Muxu, but got %v", animal)
	}
}
