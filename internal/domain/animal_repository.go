package domain

type AnimalRepository interface {
	Save(animal *Animal) error
	FindByID(id string) (*Animal, error)
	FindAll() ([]*Animal, error)
}
