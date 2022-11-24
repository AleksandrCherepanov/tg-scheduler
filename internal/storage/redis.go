package storage

import "github.com/go-redis/redis"

var client *redis.Client

func GetRedisClient() *redis.Client {
	if client == nil {
		client = redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
			Password: "",
			DB: 0,
		})
	}

	return client
}
