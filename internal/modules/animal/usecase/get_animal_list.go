package animalUseCase

import (
	"jamlink-backend/internal/modules/animal/domain"
)

type GetAnimalListUseCase struct {
	repo animalDomain.AnimalRepository
}

func NewGetAnimalListUseCase(repo animalDomain.AnimalRepository) *GetAnimalListUseCase {
	return &GetAnimalListUseCase{repo: repo}
}

func (uc *GetAnimalListUseCase) Execute() ([]*animalDomain.Animal, error) {
	animalList, err := uc.repo.FindAll()

	if err != nil {
		return nil, err
	}

	return animalList, nil
}
