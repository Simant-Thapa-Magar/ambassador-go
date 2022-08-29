package models

type User struct {
	Id           uint
	FirstName    string `json:"FirstName"`
	LastName     string `json:"LastName"`
	Email        string `json:"Email"`
	Password     string `json:"Password"`
	IsAmbassador bool   `json:"IsAmbassador"`
}
