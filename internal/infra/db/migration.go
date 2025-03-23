package db

import (
	"gorm.io/gorm"
	animalinfra "jamlink-backend/internal/modules/animal/infra"
	userinfra "jamlink-backend/internal/modules/user/infra"
	"log"
)

func MigrateDB(db *gorm.DB) {
	log.Println("🚀 Running global database migrations...")

	userinfra.MigrateUserTable(db)
	animalinfra.MigrateAnimalTable(db)

	log.Println("✅ All migrations completed successfully!")
}
