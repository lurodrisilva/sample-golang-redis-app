package main

import (
	"log/slog"
	"os"

	"github.com/redis/go-redis/v9"

	inboundhttp "github.com/lucianosilva/sample-redis-app/internal/adapter/inbound/http"
	"github.com/lucianosilva/sample-redis-app/internal/adapter/inbound/http/handler"
	"github.com/lucianosilva/sample-redis-app/internal/adapter/outbound/persistence"
	"github.com/lucianosilva/sample-redis-app/internal/application/itemapp"
	"github.com/lucianosilva/sample-redis-app/internal/infrastructure/config"
	"github.com/lucianosilva/sample-redis-app/internal/infrastructure/server"
)

func main() {
	if err := run(); err != nil {
		slog.Error("application failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// Connect to Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	defer redisClient.Close()

	// Build layers (dependency injection)
	repo := persistence.NewRedisItemRepo(redisClient)

	createHandler := itemapp.NewCreateItemHandler(repo)
	getHandler := itemapp.NewGetItemHandler(repo)

	itemHTTPHandler := handler.NewItemHandler(createHandler, getHandler)
	healthHTTPHandler := handler.NewHealthHandler()

	router := inboundhttp.NewRouter(itemHTTPHandler, healthHTTPHandler)

	// Start server
	return server.Run(":"+cfg.HTTPPort, router)
}
