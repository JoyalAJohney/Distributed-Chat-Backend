package config

import (
	"os"
	"log"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	RedisHost string
	RedisPort string
}

var Config AppConfig

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize Configuration
	Config = AppConfig{
		RedisHost: getEnv("REDIS_HOST", "localhost"),
		RedisPort: getEnv("REDIS_PORT", "6379"),
	}
}

func getEnv(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}