package inventory

import (
	"context"
	"fmt"
	"sync"

	"github.com/0xdevren/netsentry/internal/model"
)

// StaticInventory is an in-memory inventory populated from a static device list.
type StaticInventory struct {
	mu      sync.RWMutex
	devices map[string]model.Device
}

// NewStaticInventory constructs a StaticInventory from a list of devices.
func NewStaticInventory(devices []model.Device) *StaticInventory {
	inv := &StaticInventory{devices: make(map[string]model.Device, len(devices))}
	for _, d := range devices {
		id := d.ID
		if id == "" {
			id = d.Hostname
		}
		inv.devices[id] = d
	}
	return inv
}

// List returns all devices in the static inventory.
func (s *StaticInventory) List(_ context.Context) ([]model.Device, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]model.Device, 0, len(s.devices))
	for _, d := range s.devices {
		out = append(out, d)
	}
	return out, nil
}

// Get returns the device with the given ID.
func (s *StaticInventory) Get(_ context.Context, id string) (model.Device, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	d, ok := s.devices[id]
	if !ok {
		return model.Device{}, fmt.Errorf("static inventory: device %q not found", id)
	}
	return d, nil
}
