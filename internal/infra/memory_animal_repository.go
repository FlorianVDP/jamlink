package infra

import (
	"errors"
	"tindermals-backend/internal/domain"
)

type MemoryAnimalRepository struct {
	animals map[string]*domain.Animal
}

var (
	ErrAnimalNotFound = errors.New("animal not found with id")
)

func NewMemoryAnimalRepository() *MemoryAnimalRepository {

	return &MemoryAnimalRepository{animals: make(map[string]*domain.Animal)}
}

func (r *MemoryAnimalRepository) Save(animal *domain.Animal) error {
	r.animals[animal.ID] = animal

	return nil
}

func (r *MemoryAnimalRepository) FindByID(id string) (*domain.Animal, error) {
	animal, ok := r.animals[id]

	if !ok {
		return nil, ErrAnimalNotFound
	}

	return animal, nil
}

func (r *MemoryAnimalRepository) FindAll() ([]*domain.Animal, error) {
	var allAnimals []*domain.Animal

	for _, animal := range r.animals {
		allAnimals = append(allAnimals, animal)
	}

	return allAnimals, nil
}
