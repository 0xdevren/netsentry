// Package inventory provides device discovery and inventory management.
package inventory

import (
	"context"
	"github.com/0xdevren/netsentry/internal/model"
)

// Provider is the interface implemented by all inventory backends.
type Provider interface {
	// List returns all known devices from this inventory source.
	List(ctx context.Context) ([]model.Device, error)
	// Get returns a single device by its ID.
	Get(ctx context.Context, id string) (model.Device, error)
}
