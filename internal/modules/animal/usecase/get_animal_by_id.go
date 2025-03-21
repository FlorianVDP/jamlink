package animalUsecase

import (
	"tindermals-backend/internal/modules/animal/domain"
)

type GetAnimalByIdUseCase struct {
	repo animalDomain.AnimalRepository
}

func NewGetAnimalByIdUseCase(repo animalDomain.AnimalRepository) *GetAnimalByIdUseCase {
	return &GetAnimalByIdUseCase{repo: repo}
}

func (uc *GetAnimalByIdUseCase) Execute(id string) (*animalDomain.Animal, error) {
	animal, err := uc.repo.FindByID(id)

	if err != nil {
		return nil, err
	}

	return animal, nil
}
