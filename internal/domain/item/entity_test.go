package item_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lucianosilva/sample-redis-app/internal/domain/item"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name        string
		itemName    string
		description string
		wantErr     bool
	}{
		{name: "valid item", itemName: "Widget", description: "A useful widget", wantErr: false},
		{name: "empty description", itemName: "Widget", description: "", wantErr: false},
		{name: "empty name", itemName: "", description: "desc", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := item.New(tt.itemName, tt.description)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.False(t, got.ID().IsZero())
			assert.Equal(t, tt.itemName, got.Name())
			assert.Equal(t, tt.description, got.Description())
			assert.WithinDuration(t, time.Now(), got.CreatedAt(), time.Second)
		})
	}
}

func TestReconstitute(t *testing.T) {
	id := item.NewItemID()
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	got := item.Reconstitute(id, "Test", "Desc", now)

	assert.Equal(t, id, got.ID())
	assert.Equal(t, "Test", got.Name())
	assert.Equal(t, "Desc", got.Description())
	assert.Equal(t, now, got.CreatedAt())
}
