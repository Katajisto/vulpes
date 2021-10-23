package models

import (
	"log"
	"os"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func CreateTestUserService() (*UserService, error) {
	err := os.Remove("test.db")
	if err != nil {
		log.Println(err)
	}
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
	  Logger: logger.Default.LogMode(logger.Silent),
	})


	db.AutoMigrate(&User{})
	db.AutoMigrate(&Session{})

	if err != nil {
		return nil, err
	}

	return NewUserService(db), nil
}

func TestCreateUser(t *testing.T) {
	us, err := CreateTestUserService()

	if err != nil {
		t.Error(err)
	}

	user := &User{
		Username: "test",
		Password: "test",
	}

	err = us.AddUser(user)

	if err != nil {
		t.Error(err)
	}

	gotUser, err := us.GetUserByUsername("test")

	if err != nil {
		t.Error(err)
	}

	if gotUser.Username != "test" {
		t.Errorf("Expected username: test. Got: %v", gotUser.Username)
	}

	if gotUser.Password != "test" {
		t.Errorf("Expected password: test. Got: %v", gotUser.Password)
	}
}

func TestCreateOverlappingUser(t *testing.T) {
	us, err := CreateTestUserService()

	if err != nil {
		t.Error(err)
	}

	user := &User{
		Username: "test",
		Password: "test",
	}

	user2 := &User{
		Username: "test",
		Password: "test2",
	}

	err = us.AddUser(user)

	if err != nil {
		t.Error(err)
	}

	gotUser, err := us.GetUserByUsername("test")

	if err != nil {
		t.Error(err)
	}

	if gotUser.Username != "test" {
		t.Errorf("Expected username: test. Got: %v", gotUser.Username)
	}

	if gotUser.Password != "test" {
		t.Errorf("Expected password: test. Got: %v", gotUser.Password)
	}

	err = us.AddUser(user2)

	if err == nil {
		t.Error("Expected error, created 2 users with same username.")
	}
}

func TestCreateManyUser(t *testing.T) {
	us, err := CreateTestUserService()

	if err != nil {
		t.Error(err)
	}

	user := &User{
		Username: "test",
		Password: "test",
	}
	user2 := &User{
		Username: "test2",
		Password: "test2",
	}
	user3 := &User{
		Username: "test3",
		Password: "test3",
	}

	users := []*User{user, user2, user3}

	for _, user := range users {
		err = us.AddUser(user)

		if err != nil {
			t.Error(err)
		}

		gotUser, err := us.GetUserByUsername(user.Username)

		if err != nil {
			t.Error(err)
		}

		if gotUser.Username != user.Username {
			t.Errorf("Expected username: test. Got: %v", gotUser.Username)
		}

		if gotUser.Password != user.Password {
			t.Errorf("Expected password: test. Got: %v", gotUser.Password)
		}
	}

	allUsers, err := us.GetAllUsers()

	if err != nil {
		t.Error(err)
	}

	if len(allUsers) != len(users) {
		t.Errorf("Expected all users (%v) amount to be same as created users (%v).", len(allUsers), len(users))
	}

	for i, user := range allUsers {
		if user.Username != users[i].Username {
			t.Errorf("Expected username to be: %v. Got: %v", users[i].Username, user.Username)
		}

		if user.Password != users[i].Password {
			t.Errorf("Expected password to be: %v. Got: %v", users[i].Password, user.Password)
		}
	}

}

func TestSession(t *testing.T) {
	us, err := CreateTestUserService()

	if err != nil {
		t.Error(err)
	}

	user := &User{
		Username: "test",
		Password: "test",
	}

	err = us.AddUser(user)

	if err != nil {
		t.Error(err)
	}

	gotUser, err := us.GetUserByUsername("test")

	if err != nil {
		t.Error(err)
	}

	session := Session{
		UserID:    gotUser.ID,
		TokenHash: "HASHHASHHASH",
	}

	err = us.AddSession(&session)

	if err != nil {
		t.Error(err)
	}

	user2, err := us.BySession("HASHHASHHASH")

	if err != nil {
		t.Error(err)
	}

	if user2.Username != user.Username {
		t.Errorf("Expected got user to have username: %v. Got: %v", user.Username, user2.Username)
	}

	user3, err := us.BySession("HASHHasdasdadsASHHASH")

	if user3 != nil || err == nil {
		t.Error("Expected error, got user and no error with wrong hash.")
	}

	badSession := Session{
		UserID:    1337,
		TokenHash: "BAD",
	}

	err = us.AddSession(&badSession)

	if err != nil {
		t.Error(err)
	}

	user4, err := us.BySession("BAD")

	if user4 != nil || err == nil {
		t.Error("Expected error, got user and no error with wrong hash.")
	}

}

func TestDeleteUser(t *testing.T) {
	us, err := CreateTestUserService()

	if err != nil {
		t.Error(err)
	}

	user := &User{
		Username: "test",
		Password: "test",
	}

	err = us.AddUser(user)

	if err != nil {
		t.Error(err)
	}

	gotUser, err := us.GetUserByUsername("test")

	if err != nil {
		t.Error(err)
	}

	us.DeleteUser(gotUser.ID)

	_, err = us.GetUserByUsername("test")

	if err == nil {
		t.Error("Asked for deleted user. Expected error, got no error.")
	}
}
