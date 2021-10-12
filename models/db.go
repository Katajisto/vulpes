package models

import "gorm.io/gorm"

func InitDB() *gorm.DB {
	db := connectDB()

	// Migrate models
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Session{})
	db.AutoMigrate(&TgNotifyTarget{})
	MigrateData(db)
	return db
}
