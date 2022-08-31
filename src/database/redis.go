package database

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var Cache *redis.Client
var CacheChannel chan string

func RedisSetup() {
	Cache = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		DB:   0,
	})
}

func SetupCcheChannel() {
	CacheChannel = make(chan string)
	go func(ch chan string) {
		for {
			key := <-ch
			Cache.Del(context.Background(), key)
		}
	}(CacheChannel)
}

func ClearCache(keys ...string) {
	for _, key := range keys {
		CacheChannel <- key
	}
}
