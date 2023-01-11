package storage

import (
	"os"

	"github.com/go-redis/redis"
)

var client *redis.Client

func GetRedisClient() *redis.Client {
	redisHost := os.Getenv("redisHost")

	if redisHost == "" {
		redisHost = "localhost"
	}

	if client == nil {
		client = redis.NewClient(&redis.Options{
			Addr:     redisHost + ":6379",
			Password: "",
			DB:       0,
		})
	}

	return client
}
