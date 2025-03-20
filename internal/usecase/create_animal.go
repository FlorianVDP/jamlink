package usecase

import "tindermals-backend/internal/domain"

type CreateAnimalUseCase struct {
	repo domain.AnimalRepository
}

func NewCreateAnimalUseCase(repo domain.AnimalRepository) *CreateAnimalUseCase {
	return &CreateAnimalUseCase{repo: repo}
}

type CreateAnimalInput struct {
	Name        string
	Age         int
	Sexe        string
	Description string
	Image       string
}

func (uc *CreateAnimalUseCase) Execute(input CreateAnimalInput) (*domain.Animal, error) {
	animal, err := domain.NewAnimal(input.Name, input.Age, input.Sexe, input.Description, input.Image)

	if err != nil {
		return nil, err
	}

	err = uc.repo.Save(animal)
	return animal, err
}
