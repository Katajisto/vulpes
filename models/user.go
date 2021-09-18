package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string
	Password string
	Sessions []Session
}

type Session struct {
	gorm.Model
	UserID    uint
	TokenHash string
	Device    string
}

func AddSession(db *gorm.DB, session *Session) error {
	return db.Create(session).Error
}

func GetUserByUsername(db *gorm.DB, username string) (*User, error) {
	user := &User{}
	err := db.Where("username = ?", username).First(user).Error
	return user, err
}

func GetAllUsers(db *gorm.DB) ([]User, error) {
	var users []User
	err := db.Find(&users).Error
	return users, err
}

func AddUser(db *gorm.DB, user *User) error {
	return db.Create(user).Error
}

func DeleteUser(db *gorm.DB, id uint) error {
	return db.Delete(&User{}, id).Error
}
