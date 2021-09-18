package security

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

const RememberTokenBytes = 32

// Hashes a password.
func Hash(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hash)
}

// Compares two password hashes and returns ture if they match.
func Compare(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	log.Println("PASSWD COMPARE ERR: ", err)
	return err == nil
}
