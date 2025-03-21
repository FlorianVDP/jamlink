package db

import (
	"gorm.io/gorm"
	"log"
	animalinfra "tindermals-backend/internal/modules/animal/infra"
	userinfra "tindermals-backend/internal/modules/user/infra"
)

func MigrateDB(db *gorm.DB) {
	log.Println("🚀 Running global database migrations...")

	userinfra.MigrateUserTable(db)
	animalinfra.MigrateAnimalTable(db)

	log.Println("✅ All migrations completed successfully!")
}
