package database

import (
	"github.com/go-redis/redis/v9"
)

var Client *redis.Client

func RedisSetup() {
	Client = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		DB:   0,
	})
}
