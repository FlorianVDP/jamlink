package repository

import (
	"gorm.io/gorm"
	"tindermals-backend/internal/domain"
)

type PostgresAnimalRepository struct {
	db *gorm.DB
}

func NewPostgresAnimalRepository(db *gorm.DB) *PostgresAnimalRepository {
	return &PostgresAnimalRepository{db: db}
}

func (r *PostgresAnimalRepository) Create(animal *domain.Animal) error {
	return r.db.Create(animal).Error
}

func (r *PostgresAnimalRepository) FindByID(id string) (*domain.Animal, error) {
	var animal domain.Animal

	if err := r.db.First(&animal, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &animal, nil
}

func (r *PostgresAnimalRepository) FindAll() ([]*domain.Animal, error) {
	var animals []*domain.Animal

	if err := r.db.Find(&animals).Error; err != nil {
		return nil, err
	}

	return animals, nil
}
