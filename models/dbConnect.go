package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Inits sqlite database for dev purposes.

func connectDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}
