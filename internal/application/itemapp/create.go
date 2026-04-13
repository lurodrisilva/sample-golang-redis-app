package itemapp

import (
	"context"
	"fmt"

	"github.com/lucianosilva/sample-redis-app/internal/domain/item"
)

// CreateItemCommand represents the intent to create a new item.
type CreateItemCommand struct {
	Name        string
	Description string
}

// CreateItemHandler handles CreateItemCommand.
type CreateItemHandler struct {
	repo item.Repository
}

// NewCreateItemHandler creates a new CreateItemHandler.
func NewCreateItemHandler(repo item.Repository) *CreateItemHandler {
	return &CreateItemHandler{repo: repo}
}

// Handle executes the create item use case.
func (h *CreateItemHandler) Handle(ctx context.Context, cmd CreateItemCommand) (string, error) {
	i, err := item.New(cmd.Name, cmd.Description)
	if err != nil {
		return "", fmt.Errorf("create item: %w", err)
	}

	if err := h.repo.Save(ctx, i); err != nil {
		return "", fmt.Errorf("create item: save: %w", err)
	}

	return i.ID().String(), nil
}
