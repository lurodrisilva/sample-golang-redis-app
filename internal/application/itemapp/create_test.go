package itemapp_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lucianosilva/sample-redis-app/internal/application/itemapp"
	"github.com/lucianosilva/sample-redis-app/internal/domain/item"
)

type mockRepo struct {
	saveFunc     func(ctx context.Context, i *item.Item) error
	findByIDFunc func(ctx context.Context, id item.ItemID) (*item.Item, error)
}

func (m *mockRepo) Save(ctx context.Context, i *item.Item) error {
	if m.saveFunc != nil {
		return m.saveFunc(ctx, i)
	}
	return nil
}

func (m *mockRepo) FindByID(ctx context.Context, id item.ItemID) (*item.Item, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	return nil, item.ErrNotFound
}

func TestCreateItemHandler_Handle(t *testing.T) {
	tests := []struct {
		name    string
		cmd     itemapp.CreateItemCommand
		saveErr error
		wantErr bool
	}{
		{
			name:    "success",
			cmd:     itemapp.CreateItemCommand{Name: "Widget", Description: "A widget"},
			wantErr: false,
		},
		{
			name:    "empty name",
			cmd:     itemapp.CreateItemCommand{Name: "", Description: "desc"},
			wantErr: true,
		},
		{
			name:    "save error",
			cmd:     itemapp.CreateItemCommand{Name: "Widget", Description: "desc"},
			saveErr: errors.New("redis down"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepo{saveFunc: func(_ context.Context, _ *item.Item) error {
				return tt.saveErr
			}}
			h := itemapp.NewCreateItemHandler(repo)

			id, err := h.Handle(context.Background(), tt.cmd)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, id)
		})
	}
}
