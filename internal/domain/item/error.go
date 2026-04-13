package item

import "errors"

// Sentinel errors for the item aggregate.
var (
	ErrNotFound   = errors.New("item not found")
	ErrValidation = errors.New("validation error")
)
