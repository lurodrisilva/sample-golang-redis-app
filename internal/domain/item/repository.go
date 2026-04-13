package item

import "context"

// Repository defines persistence operations for the Item aggregate.
// Defined in the domain layer — implemented by outbound adapters.
type Repository interface {
	Save(ctx context.Context, item *Item) error
	FindByID(ctx context.Context, id ItemID) (*Item, error)
}
