package controllers

import (
	"ambassador/src/database"
	"ambassador/src/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func CreateProduct(c *fiber.Ctx) error {
	var product models.Product

	if err := c.BodyParser(&product); err != nil {
		return err
	}

	database.DB.Create(&product)

	return c.JSON(product)
}

func GetProduct(c *fiber.Ctx) error {
	var product models.Product
	pId := c.Params("id")
	productId, _ := strconv.Atoi(pId)
	database.DB.Where("id=?", productId).Find(&product)

	return c.JSON(product)
}

func UpdateProduct(c *fiber.Ctx) error {
	var product models.Product

	if err := c.BodyParser(&product); err != nil {
		return err
	}

	pId := c.Params("id")
	productId, _ := strconv.Atoi(pId)
	product.Id = uint(productId)

	database.DB.Updates(&product)

	return c.JSON(product)
}

func DeleteProduct(c *fiber.Ctx) error {
	var product models.Product
	pId := c.Params("id")
	productId, _ := strconv.Atoi(pId)
	product.Id = uint(productId)

	database.DB.Delete(&product)

	return c.JSON(fiber.Map{
		"message": "Product deleted",
	})
}
