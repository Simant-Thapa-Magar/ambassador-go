package controllers

import (
	"ambassador/src/database"
	"ambassador/src/middlewares"
	"ambassador/src/models"

	"github.com/bxcodec/faker/v4"
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

type CreateLinkRequest struct {
	Product []int `json:"products"`
}

func CreateLink(c *fiber.Ctx) error {
	var request CreateLinkRequest
	if err := c.BodyParser(&request); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Error parsing data",
		})
	}

	userId, e := middlewares.GetAuthenticatedUserId(c)

	if e != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Who are you ?",
		})
	}

	link := models.Link{
		Code:   faker.Username(),
		UserId: userId,
	}

	for _, productId := range request.Product {
		product := models.Product{}
		product.Id = uint(productId)
		link.Products = append(link.Products, product)
	}

	database.DB.Create(&link)
	return c.JSON(link)
}

func Stats(c *fiber.Ctx) error {
	var Links []models.Link
	var Orders []models.Order
	var result []interface{}
	userId, _ := middlewares.GetAuthenticatedUserId(c)

	database.DB.Find(&Links, models.Link{
		UserId: userId,
	})

	for _, link := range Links {
		database.DB.Preload("OrderItems").Find(&Orders, models.Order{
			Code:     link.Code,
			Complete: true,
		})

		revenue := 0.0

		for _, order := range Orders {
			revenue += order.GetTotal()
		}

		result = append(result, fiber.Map{
			"code":    link.Code,
			"count":   len(Orders),
			"revenue": revenue,
		})
	}

	return c.JSON(result)
}
