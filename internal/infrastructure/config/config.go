package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds the application configuration loaded from environment variables.
type Config struct {
	HTTPPort      string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

// Load reads configuration from environment variables with sensible defaults.
func Load() (Config, error) {
	redisDB, err := parseIntEnv("REDIS_DB", 0)
	if err != nil {
		return Config{}, fmt.Errorf("load config: %w", err)
	}

	return Config{
		HTTPPort:      envOrDefault("HTTP_PORT", "8080"),
		RedisAddr:     envOrDefault("REDIS_ADDR", "localhost:6379"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       redisDB,
	}, nil
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func parseIntEnv(key string, fallback int) (int, error) {
	v := os.Getenv(key)
	if v == "" {
		return fallback, nil
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("parse %s: %w", key, err)
	}
	return n, nil
}
