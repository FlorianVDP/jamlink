package userinfra

import (
	"jamlink-backend/internal/modules/user/domain"
	"log"

	"gorm.io/gorm"
)

func MigrateUserTable(db *gorm.DB) {
	log.Println("🚀 Running User Table Migration...")

	err := db.AutoMigrate(&userDomain.User{})
	if err != nil {
		log.Fatalf("❌ User table migration failed: %v", err)
	}

	log.Println("✅ User Table Migration completed successfully!")
}
