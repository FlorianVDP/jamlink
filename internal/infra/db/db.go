package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func ConnectDB() *gorm.DB {

	db, err := gorm.Open(postgres.Open("host=localhost user=tindermals password=password dbname=tindermals port=5432 sslmode=disable"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	return db
}
