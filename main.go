package main

import (
	"ambassador/src/database"
	route "ambassador/src/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	database.Connect()
	database.AutoMigrate()
	database.RedisSetup()
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))
	route.SetUp(app)

	log.Fatal(app.Listen(":4000"))
}
