package item_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lucianosilva/sample-redis-app/internal/domain/item"
)

func TestParseItemID(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "valid uuid", input: "550e8400-e29b-41d4-a716-446655440000", wantErr: false},
		{name: "invalid uuid", input: "not-a-uuid", wantErr: true},
		{name: "empty", input: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := item.ParseItemID(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.input, got.String())
		})
	}
}

func TestNewItemID(t *testing.T) {
	id := item.NewItemID()
	assert.False(t, id.IsZero())
}

func TestItemIDIsZero(t *testing.T) {
	var zero item.ItemID
	assert.True(t, zero.IsZero())
}
