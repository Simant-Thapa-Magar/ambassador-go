package controllers

import (
	"ambassador/src/database"
	"ambassador/src/middlewares"
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

	userId, err := middlewares.GetAuthenticatedUserId(c)

	if err != nil || userId == 0 {
		return c.JSON(fiber.Map{
			"message": "Cannot fetch user",
		})
	}

	database.DB.Where("id=?", userId).First(&user)

	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "ambassadorJWT",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "Logged Out",
	})
}

func UpdateUser(c *fiber.Ctx) error {
	var user models.User

	userId, err := middlewares.GetAuthenticatedUserId(c)

	if err != nil {
		return c.JSON(fiber.Map{
			"message": "Error getting user",
		})
	}

	if err := c.BodyParser(&user); err != nil {
		return c.JSON(fiber.Map{
			"message": "Error getting user",
		})
	}

	user.Id = userId

	database.DB.Updates(&user)

	return c.JSON(user)
}

func UpdatePassword(c *fiber.Ctx) error {
	var data models.User

	userId, err := middlewares.GetAuthenticatedUserId(c)

	if err != nil {
		return c.JSON(fiber.Map{
			"message": "Error getting user",
		})
	}

	if err := c.BodyParser(&data); err != nil {
		return c.JSON(fiber.Map{
			"message": "Error getting user",
		})
	}

	user := models.User{
		Id: userId,
	}

	e := user.SetPassword(data.Password)

	if e != nil {
		return c.JSON(fiber.Map{
			"message": "Error updating password",
		})
	}

	database.DB.Updates(&user)

	return c.JSON(user)
}
