package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func ConnectDB() *gorm.DB {
	dsn := fmt.Sprintf("host=localhost user=tindermals password=password dbname=tindermals port=5432 sslmode=disable")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	return db
}
