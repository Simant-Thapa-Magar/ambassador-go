package main

import (
	"ambassador/src/database"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	database.Connect()
	database.AutoMigrate()
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World! This is Golang")
	})

	log.Fatal(app.Listen(":4000"))
}
