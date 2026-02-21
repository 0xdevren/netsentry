// Package model defines the core domain models for NetSentry.
// All types in this package represent the canonical internal representation
// of network device state, independent of any vendor-specific syntax.
package model

import "time"

// DeviceType represents the vendor and platform family of a network device.
type DeviceType string

const (
	DeviceTypeCiscoIOS  DeviceType = "cisco-ios"
	DeviceTypeCiscoNXOS DeviceType = "cisco-nxos"
	DeviceTypeJuniperOS DeviceType = "juniper-junos"
	DeviceTypeAristaEOS DeviceType = "arista-eos"
	DeviceTypeUnknown   DeviceType = "unknown"
)

// Device represents a network device instance with its metadata.
type Device struct {
	// ID is a unique identifier for this device, typically its hostname or UUID.
	ID string `json:"id" yaml:"id"`
	// Hostname is the configured hostname of the device.
	Hostname string `json:"hostname" yaml:"hostname"`
	// Type identifies the vendor and platform.
	Type DeviceType `json:"type" yaml:"type"`
	// ManagementIP is the primary management IP address.
	ManagementIP string `json:"management_ip,omitempty" yaml:"management_ip,omitempty"`
	// Version is the operating system version string.
	Version string `json:"version,omitempty" yaml:"version,omitempty"`
	// Site is an optional location or datacenter identifier.
	Site string `json:"site,omitempty" yaml:"site,omitempty"`
	// Role is a logical classification (e.g. "core", "edge", "access").
	Role string `json:"role,omitempty" yaml:"role,omitempty"`
	// Tags is an arbitrary set of key-value metadata labels.
	Tags map[string]string `json:"tags,omitempty" yaml:"tags,omitempty"`
	// DiscoveredAt is the timestamp when this device record was created.
	DiscoveredAt time.Time `json:"discovered_at" yaml:"discovered_at"`
}

// String returns the device identifier for display purposes.
func (d Device) String() string {
	if d.Hostname != "" {
		return d.Hostname
	}
	return d.ID
}
