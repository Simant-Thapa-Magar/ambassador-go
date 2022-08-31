package main

import (
	"ambassador/src/database"
	"ambassador/src/models"
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func main() {
	database.Connect()
	database.RedisSetup()

	var ambassadorUsers []models.User

	database.DB.Find(&ambassadorUsers, models.User{
		IsAmbassador: true,
	})

	for _, ambassadorUser := range ambassadorUsers {
		ambassador := models.Ambassador(ambassadorUser)
		ambassador.CalculateRevenue(database.DB)

		fmt.Printf("Name %s score %f", ambassadorUser.GetFullname(), *ambassador.Revenue)

		database.Cache.ZAdd(context.Background(), "rankings", &redis.Z{
			Score:  *ambassador.Revenue,
			Member: ambassadorUser.GetFullname(),
		})
	}
}
