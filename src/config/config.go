package config

import (
	"log"
	"os"
)

type AppConfig struct {
	RedisHost string
	RedisPort string
}

var Config AppConfig

func init() {
	// Initialize Configuration
	Config = AppConfig{
		RedisHost: os.Getenv("REDIS_HOST"),
		RedisPort: os.Getenv("REDIS_PORT"),
	}

	if Config.RedisHost == "" || Config.RedisPort == "" {
		log.Fatal("Environment variables not set")
	}
}

func getEnv(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
