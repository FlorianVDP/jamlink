package usecase

import "tindermals-backend/internal/domain"

type GetAnimalListUseCase struct {
	repo domain.AnimalRepository
}

func NewGetAnimalListUseCase(repo domain.AnimalRepository) *GetAnimalListUseCase {
	return &GetAnimalListUseCase{repo: repo}
}

func (uc *GetAnimalListUseCase) Execute() ([]*domain.Animal, error) {
	animalList, err := uc.repo.FindAll()

	if err != nil {
		return nil, err
	}

	return animalList, nil
}
