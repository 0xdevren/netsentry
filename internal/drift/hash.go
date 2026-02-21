// Package drift provides configuration drift detection between two snapshots.
package drift

import (
	"github.com/0xdevren/netsentry/internal/util"
)

// HashEntry records the SHA-256 digest of a device configuration snapshot.
type HashEntry struct {
	// DeviceID is the identifier of the device.
	DeviceID string `json:"device_id"`
	// Hash is the SHA-256 hex digest of the raw configuration bytes.
	Hash string `json:"hash"`
}

// HashConfig computes a HashEntry for the given raw configuration bytes.
func HashConfig(deviceID string, data []byte) HashEntry {
	return HashEntry{
		DeviceID: deviceID,
		Hash:     util.SHA256Bytes(data),
	}
}

// HasChanged reports whether the configuration for a device has changed
// between baseline and current snapshots.
func HasChanged(baseline, current HashEntry) bool {
	return baseline.Hash != current.Hash
}
