package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	Id           uint
	FirstName    string `json:"FirstName"`
	LastName     string `json:"LastName"`
	Email        string `json:"Email"`
	Password     string `json:"Password"`
	IsAmbassador bool   `json:"IsAmbassador"`
}

func (user *User) SetPassword(pwd string) error {
	password, e := bcrypt.GenerateFromPassword([]byte(pwd), 12)
	user.Password = string(password)
	return e
}

func (user *User) VerifyPassword(pwd string) error {
	passError := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pwd))
	return passError
}
