package item

import (
	"fmt"
	"time"
)

// Item is the aggregate root for the item bounded context.
type Item struct {
	id          ItemID
	name        string
	description string
	createdAt   time.Time
}

// New creates a new Item. Returns an error if validation fails.
func New(name, description string) (*Item, error) {
	if name == "" {
		return nil, fmt.Errorf("%w: item name is required", ErrValidation)
	}

	return &Item{
		id:          NewItemID(),
		name:        name,
		description: description,
		createdAt:   time.Now(),
	}, nil
}

// Reconstitute rebuilds an Item from persistence. No validation.
func Reconstitute(id ItemID, name, description string, createdAt time.Time) *Item {
	return &Item{
		id:          id,
		name:        name,
		description: description,
		createdAt:   createdAt,
	}
}

// ID returns the item identifier.
func (i *Item) ID() ItemID           { return i.id }
func (i *Item) Name() string         { return i.name }
func (i *Item) Description() string  { return i.description }
func (i *Item) CreatedAt() time.Time { return i.createdAt }
