package controllers

import (
	"ambassador/src/database"
	"ambassador/src/models"
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func GetAmbassadors(c *fiber.Ctx) error {
	var ambassadors []models.User

	database.DB.Where("is_ambassador=true").Find(&ambassadors)

	return c.JSON(ambassadors)
}

func GetRanking(c *fiber.Ctx) error {
	var result []interface{}
	rankings, err := database.Cache.ZRevRangeByScoreWithScores(context.Background(), "rankings", &redis.ZRangeBy{
		Min: "-inf",
		Max: "+inf",
	}).Result()

	if err != nil {
		panic(err)
	}

	for _, rank := range rankings {
		result = append(result, fiber.Map{
			"Name":  rank.Member,
			"Score": rank.Score,
		})
	}

	return c.JSON(result)
}
