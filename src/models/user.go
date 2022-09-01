package models

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Model
	FirstName    string   `json:"first_name"`
	LastName     string   `json:"last_name"`
	Email        string   `json:"email" gorm:"unique"`
	Password     string   `json:"password"`
	IsAmbassador bool     `json:"-"`
	Revenue      *float64 `json:"revenue,omitempty" gorm:"-"`
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

func (user *User) GetFullname() string {
	return user.FirstName + " " + user.LastName
}

type Admin User

func (admin *Admin) CalculateRevenue(db *gorm.DB) {
	var orders []Order
	var revenue float64

	db.Preload("OrderItems").Where("user_id=? and complete=true", admin.Id).Find(&orders)

	for _, order := range orders {
		for _, orderItem := range order.OrderItems {
			revenue += orderItem.AdminRevenue
		}
	}

	admin.Revenue = &revenue
}

type Ambassador User

func (ambassador *Ambassador) CalculateRevenue(db *gorm.DB) {
	var orders []Order
	var revenue float64

	db.Preload("OrderItems").Where("user_id=? and complete=true", ambassador.Id).Find(&orders)

	for _, order := range orders {
		for _, orderItem := range order.OrderItems {
			revenue += orderItem.AmbassadorRevenue
		}
	}

	ambassador.Revenue = &revenue
}
