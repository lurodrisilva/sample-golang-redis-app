package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/lucianosilva/sample-redis-app/internal/application/itemapp"
	"github.com/lucianosilva/sample-redis-app/internal/domain/item"
)

// ItemHandler handles HTTP requests for items.
type ItemHandler struct {
	createItem *itemapp.CreateItemHandler
	getItem    *itemapp.GetItemHandler
}

// NewItemHandler creates a new ItemHandler.
func NewItemHandler(
	create *itemapp.CreateItemHandler,
	get *itemapp.GetItemHandler,
) *ItemHandler {
	return &ItemHandler{createItem: create, getItem: get}
}

// Create handles POST /items.
func (h *ItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	cmd := itemapp.CreateItemCommand{
		Name:        req.Name,
		Description: req.Description,
	}

	id, err := h.createItem.Handle(r.Context(), cmd)
	if err != nil {
		handleAppError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

// Get handles GET /items/{id}.
func (h *ItemHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	q := itemapp.GetItemQuery{ItemID: id}

	dto, err := h.getItem.Handle(r.Context(), q)
	if err != nil {
		handleAppError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto)
}

type createItemRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func handleAppError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, item.ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, item.ErrValidation):
		writeError(w, http.StatusBadRequest, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "internal error")
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
