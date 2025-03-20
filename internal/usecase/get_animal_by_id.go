package usecase

import "tindermals-backend/internal/domain"

type GetAnimalByIdUseCase struct {
	repo domain.AnimalRepository
}

func NewGetAnimalByIdUseCase(repo domain.AnimalRepository) *GetAnimalByIdUseCase {
	return &GetAnimalByIdUseCase{repo: repo}
}

func (uc *GetAnimalByIdUseCase) Execute(id string) (*domain.Animal, error) {
	animal, err := uc.repo.FindByID(id)

	if err != nil {
		return nil, err
	}

	return animal, nil
}
