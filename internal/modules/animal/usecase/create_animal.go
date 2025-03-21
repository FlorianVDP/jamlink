package animalUseCase

import (
	"tindermals-backend/internal/modules/animal/domain"
)

type CreateAnimalUseCase struct {
	repo animalDomain.AnimalRepository
}

func NewCreateAnimalUseCase(repo animalDomain.AnimalRepository) *CreateAnimalUseCase {
	return &CreateAnimalUseCase{repo: repo}
}

type CreateAnimalInput struct {
	Name        string
	Age         int
	Sexe        string
	Description string
	Image       string
}

func (uc *CreateAnimalUseCase) Execute(input CreateAnimalInput) (*animalDomain.Animal, error) {
	animal, err := animalDomain.CreateAnimal(input.Name, input.Age, input.Sexe, input.Description, input.Image)

	if err != nil {
		return nil, err
	}

	err = uc.repo.Create(animal)
	return animal, err
}
