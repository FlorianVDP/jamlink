package animalMocks

import (
	"tindermals-backend/internal/modules/animal/domain"
)

type MockAnimalRepository struct {
	Animals []*animalDomain.Animal
}

func (m *MockAnimalRepository) FindAll() ([]*animalDomain.Animal, error) {
	return m.Animals, nil
}

func (m *MockAnimalRepository) Create(animal *animalDomain.Animal) error {
	m.Animals = append(m.Animals, animal)
	return nil
}

func (m *MockAnimalRepository) FindByID(id string) (*animalDomain.Animal, error) {
	for _, animal := range m.Animals {
		if animal.ID.String() == id {
			return animal, nil
		}
	}
	return nil, nil
}
