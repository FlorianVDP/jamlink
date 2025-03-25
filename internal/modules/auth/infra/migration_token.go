package userinfra

import (
	"jamlink-backend/internal/modules/auth/domain/token"
	"log"

	"gorm.io/gorm"
)

func MigrateTokenTable(db *gorm.DB) {
	log.Println("🚀 Running Token Table Migration...")

	err := db.AutoMigrate(&token.Token{})
	if err != nil {
		log.Fatalf("❌ Token table migration failed: %v", err)
	}

	log.Println("✅ Token Table Migration completed successfully!")
}
