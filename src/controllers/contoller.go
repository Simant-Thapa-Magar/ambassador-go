package controllers

import (
	"ambassador/src/database"
	"ambassador/src/models"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func Register(c *fiber.Ctx) error {
	var data models.User

	if err := c.BodyParser(&data); err != nil {
		fmt.Println("Error while parsing body")
		return err
	}

	user := models.User{
		FirstName:    data.FirstName,
		LastName:     data.LastName,
		Email:        data.Email,
		IsAmbassador: false,
	}

	user.SetPassword(data.Password)

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

	if user.VerifyPassword(data.Password) != nil {
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

func User(c *fiber.Ctx) error {
	var user models.User
	cookie := c.Cookies("ambassadorJWT")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if err != nil || !token.Valid {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthorized access",
		})
	}

	payload := token.Claims.(*jwt.StandardClaims)

	database.DB.Where("id=?", payload.Subject).First(&user)

	return c.JSON(user)
}
