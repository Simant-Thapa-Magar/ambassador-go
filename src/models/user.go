package models

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id           uint
	FirstName    string `json:"FirstName"`
	LastName     string `json:"LastName"`
	Email        string `json:"Email" gorm:"unique"`
	Password     string `json:"Password"`
	IsAmbassador bool   `json:"-"`
}

func (user *User) SetPassword(pwd string) error {
	fmt.Println("Update password for user ", user, " with ", pwd)
	password, e := bcrypt.GenerateFromPassword([]byte(pwd), 12)
	user.Password = string(password)
	return e
}

func (user *User) VerifyPassword(pwd string) error {
	passError := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pwd))
	return passError
}
