package animalRepository

import (
	"jamlink-backend/internal/modules/animal/domain"
	"jamlink-backend/internal/modules/animal/repository/errors"
)

type MemoryAnimalRepository struct {
	animals map[string]*animalDomain.Animal
}

func NewMemoryAnimalRepository() *MemoryAnimalRepository {

	return &MemoryAnimalRepository{animals: make(map[string]*animalDomain.Animal)}
}

func (r *MemoryAnimalRepository) Create(animal *animalDomain.Animal) error {
	r.animals[animal.ID.String()] = animal

	return nil
}

func (r *MemoryAnimalRepository) FindByID(id string) (*animalDomain.Animal, error) {
	animal, ok := r.animals[id]

	if !ok {
		return nil, errors.ErrAnimalNotFound
	}

	return animal, nil
}

func (r *MemoryAnimalRepository) FindAll() ([]*animalDomain.Animal, error) {
	var allAnimals []*animalDomain.Animal

	for _, animal := range r.animals {
		allAnimals = append(allAnimals, animal)
	}

	return allAnimals, nil
}
