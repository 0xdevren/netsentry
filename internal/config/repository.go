package config

import (
	"fmt"
	"sync"
)

// ConfigEntry is a cached raw configuration record.
type ConfigEntry struct {
	// DeviceID is the identifier of the device this configuration belongs to.
	DeviceID string
	// Data is the raw configuration bytes.
	Data []byte
	// Checksum is the SHA-256 hex digest of Data for integrity verification.
	Checksum string
}

// Repository provides an in-memory store for loaded device configurations.
type Repository struct {
	mu      sync.RWMutex
	entries map[string]ConfigEntry
}

// NewRepository constructs an empty Repository.
func NewRepository() *Repository {
	return &Repository{entries: make(map[string]ConfigEntry)}
}

// Store saves a ConfigEntry, keyed on its DeviceID.
func (r *Repository) Store(entry ConfigEntry) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[entry.DeviceID] = entry
}

// Get retrieves the ConfigEntry for the given deviceID. Returns an error if
// no entry exists.
func (r *Repository) Get(deviceID string) (ConfigEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	e, ok := r.entries[deviceID]
	if !ok {
		return ConfigEntry{}, fmt.Errorf("repository: no config for device %q", deviceID)
	}
	return e, nil
}

// Delete removes the entry for the given deviceID.
func (r *Repository) Delete(deviceID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.entries, deviceID)
}

// All returns a snapshot of all stored entries.
func (r *Repository) All() []ConfigEntry {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]ConfigEntry, 0, len(r.entries))
	for _, e := range r.entries {
		out = append(out, e)
	}
	return out
}

// Count returns the number of stored entries.
func (r *Repository) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.entries)
}
