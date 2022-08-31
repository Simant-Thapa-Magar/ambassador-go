package controllers

import (
	"ambassador/src/database"
	"ambassador/src/models"
	"context"
	"encoding/json"
	"sort"
	"strconv"
	"strings"
	"time"

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

func ProductFrontend(c *fiber.Ctx) error {
	var products []models.Product

	ctx := context.Background()

	cachedProducts, err := database.Cache.Get(ctx, "products_frontend").Result()

	if err != nil {
		// from db
		database.DB.Find(&products)
		productBytes, errMarshal := json.Marshal(products)
		if errMarshal != nil {
			panic(errMarshal)
		}
		if errKey := database.Cache.Set(ctx, "products_frontend", productBytes, 30*time.Minute).Err(); errKey != nil {
			panic(errKey)
		}
	} else {
		// from cache
		json.Unmarshal([]byte(cachedProducts), &products)
	}

	return c.JSON(products)
}

func ProductBackend(c *fiber.Ctx) error {
	var products, searchedProducts []models.Product

	ctx := context.Background()

	cachedProducts, err := database.Cache.Get(ctx, "products_backend").Result()

	if err != nil {
		database.DB.Find(&products)
		productBytes, errMarshal := json.Marshal(products)
		if errMarshal != nil {
			panic(errMarshal)
		}
		if errKey := database.Cache.Set(ctx, "products_backend", productBytes, 30*time.Minute).Err(); errKey != nil {
			panic(errKey)
		}
	} else {
		json.Unmarshal([]byte(cachedProducts), &products)
	}

	if s := c.Query("q"); s != "" {
		lower := strings.ToLower(s)
		for _, product := range products {
			if strings.Contains(strings.ToLower(product.Title), lower) || strings.Contains(strings.ToLower(product.Description), lower) {
				searchedProducts = append(searchedProducts, product)
			}
		}
	} else {
		searchedProducts = products
	}

	if sortParam := c.Query("sort"); sortParam != "" {
		lower := strings.ToLower(sortParam)
		if lower == "asc" {
			sort.Slice(searchedProducts, func(i, j int) bool {
				return searchedProducts[i].Price < searchedProducts[j].Price
			})
		} else if lower == "desc" {
			sort.Slice(searchedProducts, func(i, j int) bool {
				return searchedProducts[i].Price > searchedProducts[j].Price
			})
		}
	}

	return c.JSON(searchedProducts)
}
