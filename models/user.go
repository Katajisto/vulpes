package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex"`
	Password string
	Sessions []Session
}

type Session struct {
	gorm.Model
	UserID    uint
	TokenHash string
	Device    string
}

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (us *UserService) AddSession(session *Session) error {
	return us.db.Create(session).Error
}

func (us *UserService) GetUserByUsername(username string) (*User, error) {
	user := &User{}
	err := us.db.Where("username = ?", username).First(user).Error
	return user, err
}

func (us *UserService) GetAllUsers() ([]User, error) {
	var users []User
	err := us.db.Find(&users).Error
	return users, err
}

func (us *UserService) AddUser(user *User) error {
	return us.db.Create(user).Error
}

func (us *UserService) DeleteUser(id uint) error {
	return us.db.Delete(&User{}, id).Error
}

func (us *UserService) BySession(tokenHash string) (*User, error) {
	session := &Session{}
	err := us.db.First(&session, "token_hash = ?", tokenHash).Error
	if err != nil {
		return nil, err
	}
	user := &User{}
	err = us.db.First(user, session.UserID).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
