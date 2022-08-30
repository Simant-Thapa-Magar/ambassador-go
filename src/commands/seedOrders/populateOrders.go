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
		var orderItems []models.OrderItem
		randomOrderItemCount := rand.Intn(5) + 1
		for j := 0; j < randomOrderItemCount; j++ {
			price := rand.Intn(90) + 10
			quantity := rand.Intn(10) + 1
			orderItems = append(orderItems, models.OrderItem{
				ProductTitle:      faker.Word(),
				Price:             float64(price),
				Quantity:          uint(quantity),
				AdminRevenue:      0.9 * float64(price) * float64(quantity),
				AmbassadorRevenue: 0.1 * float64(price) * float64(quantity),
			})
		}

		order := models.Order{
			TransactionId:   faker.Word(),
			UserId:          uint(rand.Intn(30) + 1),
			Code:            faker.Username(),
			AmbassadorEmail: faker.Email(),
			FirstName:       faker.FirstName(),
			LastName:        faker.LastName(),
			Email:           faker.Email(),
			Complete:        true,
			OrderItems:      orderItems,
		}

		database.DB.Create(&order)
	}
}
