package database

import (
	"github.com/go-redis/redis/v9"
)

var Cache *redis.Client

func RedisSetup() {
	Cache = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		DB:   0,
	})
}
