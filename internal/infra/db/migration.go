package db

import (
	"gorm.io/gorm"
	"log"
	"tindermals-backend/internal/domain"
)

func MigrateDB(db *gorm.DB) {
	err := db.AutoMigrate(&domain.Animal{})

	if err != nil {
		log.Fatal("Migration error:", err)
	}
}
