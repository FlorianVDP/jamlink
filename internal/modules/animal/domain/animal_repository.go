package animalDomain

type AnimalRepository interface {
	Create(animal *Animal) error
	FindByID(id string) (*Animal, error)
	FindAll() ([]*Animal, error)
}
