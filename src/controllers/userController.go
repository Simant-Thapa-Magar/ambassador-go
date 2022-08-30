package controllers

import (
	"ambassador/src/database"
	"ambassador/src/models"

	"github.com/gofiber/fiber/v2"
)

func GetAmbassadors(c *fiber.Ctx) error {
	var ambassadors []models.User

	database.DB.Where("is_ambassador=true").Find(&ambassadors)

	return c.JSON(ambassadors)
}
