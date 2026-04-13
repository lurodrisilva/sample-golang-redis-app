package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lucianosilva/sample-redis-app/internal/infrastructure/config"
)

func TestLoad_Defaults(t *testing.T) {
	t.Setenv("HTTP_PORT", "")
	t.Setenv("REDIS_ADDR", "")
	t.Setenv("REDIS_PASSWORD", "")
	t.Setenv("REDIS_DB", "")

	cfg, err := config.Load()
	require.NoError(t, err)

	assert.Equal(t, "8080", cfg.HTTPPort)
	assert.Equal(t, "localhost:6379", cfg.RedisAddr)
	assert.Equal(t, "", cfg.RedisPassword)
	assert.Equal(t, 0, cfg.RedisDB)
}

func TestLoad_CustomValues(t *testing.T) {
	t.Setenv("HTTP_PORT", "9090")
	t.Setenv("REDIS_ADDR", "redis:6380")
	t.Setenv("REDIS_PASSWORD", "secret")
	t.Setenv("REDIS_DB", "2")

	cfg, err := config.Load()
	require.NoError(t, err)

	assert.Equal(t, "9090", cfg.HTTPPort)
	assert.Equal(t, "redis:6380", cfg.RedisAddr)
	assert.Equal(t, "secret", cfg.RedisPassword)
	assert.Equal(t, 2, cfg.RedisDB)
}

func TestLoad_InvalidRedisDB(t *testing.T) {
	t.Setenv("REDIS_DB", "not-a-number")

	_, err := config.Load()
	require.Error(t, err)
}
