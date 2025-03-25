package db

import (
	"gorm.io/gorm"
	userinfra "jamlink-backend/internal/modules/auth/infra"
	"log"
)

func MigrateDB(db *gorm.DB) {
	log.Println("🚀 Running global database migrations...")

	userinfra.MigrateUserTable(db)
	userinfra.MigrateTokenTable(db)

	log.Println("✅ All migrations completed successfully!")
}
