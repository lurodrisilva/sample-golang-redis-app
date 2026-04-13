package http

import (
	"net/http"

	"github.com/lucianosilva/sample-redis-app/internal/adapter/inbound/http/handler"
	"github.com/lucianosilva/sample-redis-app/internal/adapter/inbound/http/middleware"
)

// NewRouter creates and configures the HTTP router.
func NewRouter(itemHandler *handler.ItemHandler, healthHandler *handler.HealthHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /items", itemHandler.Create)
	mux.HandleFunc("GET /items/{id}", itemHandler.Get)
	mux.HandleFunc("GET /health/live", healthHandler.Live)

	var h http.Handler = mux
	h = middleware.Recoverer(h)

	return h
}
