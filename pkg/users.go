package users

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	// "github.com/go-sql-driver/mysql"
)

type User struct {
	username string
	password string
}

func getHashedPassword(i string) string {
	password := []byte(i)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hashedPassword)
}
