// +build prod

package models

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func connectDB() *gorm.DB {

	log.Println("Opening up database connection.")

	db, err := gorm.Open(postgres.Open(os.Getenv("DB_CONNECTION_STRING")), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}
