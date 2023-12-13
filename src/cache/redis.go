package cache

import (
	"fmt"
	"context"

	"github.com/go-redis/redis/v8"

	"realtime-chat/src/config"
)

type RedisConfig struct {
	Host string
	Port string
}
var redisConfig RedisConfig
var RedisClient *redis.Client


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

	// Test connection
	ctx := context.Background()
	if err := RedisClient.Ping(ctx).Err(); err != nil {
		panic(err)
	}
}

