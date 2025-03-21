package animalinfra

import (
	"log"
	"tindermals-backend/internal/modules/animal/domain"

	"gorm.io/gorm"
)

func MigrateAnimalTable(db *gorm.DB) {
	log.Println("🚀 Running Animal Table Migration...")

	err := db.AutoMigrate(&animalDomain.Animal{})
	if err != nil {
		log.Fatalf("❌ Animal table migration failed: %v", err)
	}

	log.Println("✅ Animal Table Migration completed successfully!")
}
