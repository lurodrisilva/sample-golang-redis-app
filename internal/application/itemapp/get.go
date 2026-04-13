package itemapp

import (
	"context"
	"fmt"

	"github.com/lucianosilva/sample-redis-app/internal/domain/item"
)

// GetItemQuery represents a request to retrieve an item.
type GetItemQuery struct {
	ItemID string
}

// ItemDTO is the read-model output.
type ItemDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

// GetItemHandler handles GetItemQuery.
type GetItemHandler struct {
	repo item.Repository
}

// NewGetItemHandler creates a new GetItemHandler.
func NewGetItemHandler(repo item.Repository) *GetItemHandler {
	return &GetItemHandler{repo: repo}
}

// Handle retrieves an item and maps to DTO.
func (h *GetItemHandler) Handle(ctx context.Context, q GetItemQuery) (ItemDTO, error) {
	id, err := item.ParseItemID(q.ItemID)
	if err != nil {
		return ItemDTO{}, fmt.Errorf("get item: %w", err)
	}

	i, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return ItemDTO{}, fmt.Errorf("get item: %w", err)
	}

	return toItemDTO(i), nil
}

func toItemDTO(i *item.Item) ItemDTO {
	return ItemDTO{
		ID:          i.ID().String(),
		Name:        i.Name(),
		Description: i.Description(),
		CreatedAt:   i.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}
}
