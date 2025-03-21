package animalRepository

import (
	"gorm.io/gorm"
	"tindermals-backend/internal/modules/animal/domain"
)

type PostgresAnimalRepository struct {
	db *gorm.DB
}

func NewPostgresAnimalRepository(db *gorm.DB) *PostgresAnimalRepository {
	return &PostgresAnimalRepository{db: db}
}

func (r *PostgresAnimalRepository) Create(animal *animalDomain.Animal) error {
	return r.db.Create(animal).Error
}

func (r *PostgresAnimalRepository) FindByID(id string) (*animalDomain.Animal, error) {
	var animal animalDomain.Animal

	if err := r.db.First(&animal, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &animal, nil
}

func (r *PostgresAnimalRepository) FindAll() ([]*animalDomain.Animal, error) {
	var animals []*animalDomain.Animal

	if err := r.db.Find(&animals).Error; err != nil {
		return nil, err
	}

	return animals, nil
}
