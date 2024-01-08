package cache

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"

	"realtime-chat/src/config"
)

type RedisConfig struct {
	Host string
	Port string
}

var redisConfig RedisConfig
var RedisClient *redis.Client
var PubSubConnection *redis.PubSub

func init() {
	redisConfig = RedisConfig{
		Host: config.Config.RedisHost,
		Port: config.Config.RedisPort,
	}
}

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisConfig.Host, redisConfig.Port),
	})

	// Initialize PubSub connection
	ctx := context.Background()
	PubSubConnection = RedisClient.Subscribe(ctx)

	// Ping Redis to check if connection is established
	if err := RedisClient.Ping(ctx).Err(); err != nil {
		panic(err)
	}
}
