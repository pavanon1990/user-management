package config

import (
	"os"
	"time"
)

type AppConfig struct {
	Port      string
	SecretKey string
	TokenTTL  time.Duration
}

func LoadAppConfig() AppConfig {
	return AppConfig{
		Port:      getEnvOrDefault("APP_PORT", "8080"),
		SecretKey: os.Getenv("SECRET_KEY"),
		TokenTTL:  resolveDurationTime("TOKEN_TTL", 24*time.Hour),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val
}
