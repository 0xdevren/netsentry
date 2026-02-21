package model

// VLAN represents a single VLAN definition on a device.
type VLAN struct {
	// ID is the VLAN identifier (1-4094).
	ID int `json:"id" yaml:"id"`
	// Name is the operator-assigned VLAN name.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// State is the administrative state: "active" or "suspend".
	State string `json:"state,omitempty" yaml:"state,omitempty"`
}
