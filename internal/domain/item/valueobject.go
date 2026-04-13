package item

import (
	"fmt"

	"github.com/google/uuid"
)

// ItemID uniquely identifies an item.
type ItemID uuid.UUID

// NewItemID creates a new random ItemID.
func NewItemID() ItemID {
	return ItemID(uuid.New())
}

// ParseItemID parses a string into an ItemID.
func ParseItemID(s string) (ItemID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return ItemID{}, fmt.Errorf("parse item id: %w", err)
	}
	return ItemID(id), nil
}

// String returns the string representation.
func (id ItemID) String() string {
	return uuid.UUID(id).String()
}

// IsZero reports whether the ID is the zero value.
func (id ItemID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}
