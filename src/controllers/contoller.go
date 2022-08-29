package controllers

import (
	"ambassador/src/database"
	"ambassador/src/models"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	var data models.User

	if err := c.BodyParser(&data); err != nil {
		fmt.Println("Error while parsing body")
		return err
	}

	password, e := bcrypt.GenerateFromPassword([]byte(data.Password), 12)

	if e != nil {
		fmt.Println("Error on password bcrypt")
		return e
	}

	user := models.User{
		FirstName:    data.FirstName,
		LastName:     data.LastName,
		Email:        data.Email,
		Password:     string(password),
		IsAmbassador: false,
	}

	database.DB.Create(&user)

	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	var user, data models.User

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	database.DB.Where("email = ?", data.Email).First(&user)

	if user.Id == 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "No user found",
		})
	}

	if passError := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); passError != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Wrong password",
		})
	}

	payload := jwt.StandardClaims{
		Subject:   strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	}

	token, e := jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString([]byte("secret"))

	if e != nil {
		return e
	}

	cookie := fiber.Cookie{
		Name:     "ambassadorJWT",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "Token generated successfully",
	})
}
