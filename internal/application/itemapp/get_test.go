package itemapp_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lucianosilva/sample-redis-app/internal/application/itemapp"
	"github.com/lucianosilva/sample-redis-app/internal/domain/item"
)

func TestGetItemHandler_Handle(t *testing.T) {
	fixedID := item.NewItemID()
	fixedTime := time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)
	existing := item.Reconstitute(fixedID, "Widget", "A widget", fixedTime)

	tests := []struct {
		name    string
		query   itemapp.GetItemQuery
		repo    *mockRepo
		wantErr bool
	}{
		{
			name:  "found",
			query: itemapp.GetItemQuery{ItemID: fixedID.String()},
			repo: &mockRepo{findByIDFunc: func(_ context.Context, _ item.ItemID) (*item.Item, error) {
				return existing, nil
			}},
			wantErr: false,
		},
		{
			name:  "not found",
			query: itemapp.GetItemQuery{ItemID: fixedID.String()},
			repo: &mockRepo{findByIDFunc: func(_ context.Context, _ item.ItemID) (*item.Item, error) {
				return nil, item.ErrNotFound
			}},
			wantErr: true,
		},
		{
			name:    "invalid id",
			query:   itemapp.GetItemQuery{ItemID: "bad"},
			repo:    &mockRepo{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := itemapp.NewGetItemHandler(tt.repo)

			dto, err := h.Handle(context.Background(), tt.query)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, fixedID.String(), dto.ID)
			assert.Equal(t, "Widget", dto.Name)
			assert.Equal(t, "A widget", dto.Description)
		})
	}
}
