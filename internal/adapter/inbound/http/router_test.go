package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	inboundhttp "github.com/lucianosilva/sample-redis-app/internal/adapter/inbound/http"
	"github.com/lucianosilva/sample-redis-app/internal/adapter/inbound/http/handler"
	"github.com/lucianosilva/sample-redis-app/internal/application/itemapp"
	"github.com/lucianosilva/sample-redis-app/internal/domain/item"
)

type stubRepo struct {
	items map[string]*item.Item
}

func newStubRepo() *stubRepo {
	return &stubRepo{items: make(map[string]*item.Item)}
}

func (s *stubRepo) Save(_ context.Context, i *item.Item) error {
	s.items[i.ID().String()] = i
	return nil
}

func (s *stubRepo) FindByID(_ context.Context, id item.ItemID) (*item.Item, error) {
	i, ok := s.items[id.String()]
	if !ok {
		return nil, item.ErrNotFound
	}
	return i, nil
}

func newRouter(t *testing.T) http.Handler {
	t.Helper()
	repo := newStubRepo()
	itemH := handler.NewItemHandler(
		itemapp.NewCreateItemHandler(repo),
		itemapp.NewGetItemHandler(repo),
	)
	healthH := handler.NewHealthHandler()
	return inboundhttp.NewRouter(itemH, healthH)
}

func TestRouter_HealthLive(t *testing.T) {
	router := newRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/health/live", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRouter_CreateItem(t *testing.T) {
	router := newRouter(t)

	body := `{"name":"Widget","description":"A widget"}`
	req := httptest.NewRequest(http.MethodPost, "/items", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]string
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.NotEmpty(t, resp["id"])
}

func TestRouter_GetItem(t *testing.T) {
	repo := newStubRepo()
	fixedID := item.NewItemID()
	repo.items[fixedID.String()] = item.Reconstitute(fixedID, "Widget", "desc", time.Now())

	itemH := handler.NewItemHandler(
		itemapp.NewCreateItemHandler(repo),
		itemapp.NewGetItemHandler(repo),
	)
	healthH := handler.NewHealthHandler()
	router := inboundhttp.NewRouter(itemH, healthH)

	req := httptest.NewRequest(http.MethodGet, "/items/"+fixedID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRouter_PanicRecovery(t *testing.T) {
	router := newRouter(t)

	// Requesting a non-existent route should not panic
	req := httptest.NewRequest(http.MethodDelete, "/nonexistent", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// 404 or 405 is expected, just verify no panic
	assert.True(t, w.Code >= 400)
}
