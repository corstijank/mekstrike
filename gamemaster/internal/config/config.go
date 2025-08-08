package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port        int
	PubsubName  string
	Environment string
}

func Load() *Config {
	return &Config{
		Port:        getEnvAsInt("PORT", 7011),
		PubsubName:  getEnv("PUBSUB_NAME", "redis-pubsub"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}