package controllers

import (
	"ambassador/src/database"
	"ambassador/src/models"

	"github.com/gofiber/fiber/v2"
)

func GetUserLinks(c *fiber.Ctx) error {
	var links []models.Link
	userId := c.Params("id")

	database.DB.Where("user_id=?", userId).Find(&links)

	for i, link := range links {
		var orders []models.Order
		database.DB.Where("code = ? and complete = true", link.Code).Find(&orders)
		links[i].Orders = orders
	}

	return c.JSON(links)
}
