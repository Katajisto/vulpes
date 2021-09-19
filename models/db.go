package models

import "gorm.io/gorm"

func InitDB() *gorm.DB {
	db := connectDB()

	// Migrate models
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Session{})
	return db
}
