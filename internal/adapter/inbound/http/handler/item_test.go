package handler_test

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

func TestItemHandler_Create(t *testing.T) {
	repo := newStubRepo()
	h := handler.NewItemHandler(
		itemapp.NewCreateItemHandler(repo),
		itemapp.NewGetItemHandler(repo),
	)

	body := `{"name":"Widget","description":"A test widget"}`
	req := httptest.NewRequest(http.MethodPost, "/items", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Create(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]string
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.NotEmpty(t, resp["id"])
}

func TestItemHandler_Create_InvalidBody(t *testing.T) {
	repo := newStubRepo()
	h := handler.NewItemHandler(
		itemapp.NewCreateItemHandler(repo),
		itemapp.NewGetItemHandler(repo),
	)

	req := httptest.NewRequest(http.MethodPost, "/items", bytes.NewBufferString("not json"))
	w := httptest.NewRecorder()

	h.Create(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestItemHandler_Create_EmptyName(t *testing.T) {
	repo := newStubRepo()
	h := handler.NewItemHandler(
		itemapp.NewCreateItemHandler(repo),
		itemapp.NewGetItemHandler(repo),
	)

	body := `{"name":"","description":"desc"}`
	req := httptest.NewRequest(http.MethodPost, "/items", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Create(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestItemHandler_Get(t *testing.T) {
	repo := newStubRepo()
	fixedID := item.NewItemID()
	repo.items[fixedID.String()] = item.Reconstitute(fixedID, "Widget", "desc", time.Now())

	h := handler.NewItemHandler(
		itemapp.NewCreateItemHandler(repo),
		itemapp.NewGetItemHandler(repo),
	)

	req := httptest.NewRequest(http.MethodGet, "/items/"+fixedID.String(), nil)
	req.SetPathValue("id", fixedID.String())
	w := httptest.NewRecorder()

	h.Get(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var dto itemapp.ItemDTO
	require.NoError(t, json.NewDecoder(w.Body).Decode(&dto))
	assert.Equal(t, fixedID.String(), dto.ID)
	assert.Equal(t, "Widget", dto.Name)
}

func TestItemHandler_Get_NotFound(t *testing.T) {
	repo := newStubRepo()
	h := handler.NewItemHandler(
		itemapp.NewCreateItemHandler(repo),
		itemapp.NewGetItemHandler(repo),
	)

	id := item.NewItemID()
	req := httptest.NewRequest(http.MethodGet, "/items/"+id.String(), nil)
	req.SetPathValue("id", id.String())
	w := httptest.NewRecorder()

	h.Get(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestItemHandler_Get_InvalidID(t *testing.T) {
	repo := newStubRepo()
	h := handler.NewItemHandler(
		itemapp.NewCreateItemHandler(repo),
		itemapp.NewGetItemHandler(repo),
	)

	req := httptest.NewRequest(http.MethodGet, "/items/bad-id", nil)
	req.SetPathValue("id", "bad-id")
	w := httptest.NewRecorder()

	h.Get(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
