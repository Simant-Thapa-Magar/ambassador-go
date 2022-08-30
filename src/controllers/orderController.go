package controllers

import (
	"ambassador/src/database"
	"ambassador/src/models"

	"github.com/gofiber/fiber/v2"
)

func GetOrders(c *fiber.Ctx) error {
	var orders []models.Order

	database.DB.Preload("OrderItems").Find(&orders)

	for i, order := range orders {
		orders[i].Name = order.GetFullname()
		orders[i].Total = order.GetTotal()
	}

	return c.JSON(orders)
}
