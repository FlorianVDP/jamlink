package mocks

import "tindermals-backend/internal/domain"

type MockAnimalRepository struct {
	Animals []*domain.Animal
}

func (m *MockAnimalRepository) FindAll() ([]*domain.Animal, error) {
	return m.Animals, nil
}

func (m *MockAnimalRepository) Save(animal *domain.Animal) error {
	m.Animals = append(m.Animals, animal)
	return nil
}

func (m *MockAnimalRepository) FindByID(id string) (*domain.Animal, error) {
	for _, animal := range m.Animals {
		if animal.ID == id {
			return animal, nil
		}
	}
	return nil, nil
}
