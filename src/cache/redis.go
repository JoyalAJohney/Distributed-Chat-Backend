package cache

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var RedisConnection *redis.Client

func InitRedis() {
	RedisConnection = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})

	// Test connection
	ctx := context.Background()
	_, err := RedisConnection.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
}

