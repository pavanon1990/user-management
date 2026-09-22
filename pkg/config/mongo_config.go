package config

import (
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type MongoConfig struct {
	URI               string
	Database          string
	MaxPoolSize       uint64
	MinPoolSize       uint64
	MaxIdleTime       time.Duration
	ConnectionTimeout time.Duration
}

const DB_NAME = "users"

func LoadMongoConfig() MongoConfig {
	uri := os.Getenv("MONGO_URI")

	replacer := strings.NewReplacer(
		"{user}", os.Getenv("MONGO_USER"),
		"{password}", url.QueryEscape(os.Getenv("MONGO_PASSWORD")),
		"{host}", os.Getenv("MONGO_HOST"),
		"{port}", os.Getenv("MONGO_PORT"),
	)

	return MongoConfig{
		URI:               replacer.Replace(uri),
		Database:          "user_management",
		MinPoolSize:       resolvePoolSize("MONGO_MINPOOL_SIZE", 3),
		MaxPoolSize:       resolvePoolSize("MONGO_MAXPOOL_SIZE", 30),
		MaxIdleTime:       resolveDurationTime("MONGO_MAX_IDLE_TIME", 30*time.Second),
		ConnectionTimeout: resolveDurationTime("MONGO_CONNECT_TIMEOUT", 10*time.Second),
	}
}

func resolveDurationTime(key string, defaultValue time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	d, err := time.ParseDuration(val)
	if err != nil {
		return defaultValue
	}
	return d
}

func resolvePoolSize(key string, defaultValue uint64) uint64 {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	d, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return defaultValue
	}
	return d
}
