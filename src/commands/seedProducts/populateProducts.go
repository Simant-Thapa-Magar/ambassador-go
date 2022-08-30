package main

import (
	"ambassador/src/database"
	"ambassador/src/models"
	"math/rand"

	"github.com/bxcodec/faker/v4"
)

func main() {
	database.Connect()
	for i := 0; i < 30; i++ {
		product := models.Product{
			Title:       faker.Username(),
			Description: faker.Username(),
			Image:       faker.Username(),
			Price:       float64(rand.Intn(100)),
		}

		database.DB.Create(&product)
	}
}
