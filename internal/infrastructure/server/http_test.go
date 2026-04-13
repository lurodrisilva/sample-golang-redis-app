package server_test

import (
	"net/http"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lucianosilva/sample-redis-app/internal/infrastructure/server"
)

func TestRun_GracefulShutdown(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.Run(":0", handler)
	}()

	// Give the server time to start
	time.Sleep(100 * time.Millisecond)

	// Send SIGINT to trigger graceful shutdown
	require.NoError(t, syscall.Kill(syscall.Getpid(), syscall.SIGINT))

	select {
	case err := <-errCh:
		assert.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("server did not shut down in time")
	}
}

func TestRun_InvalidAddr(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Use an invalid address to trigger ListenAndServe error
	err := server.Run(":-1", handler)
	assert.Error(t, err)
}
