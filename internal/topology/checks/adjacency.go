// Package checks defines topology check interfaces and built-in implementations.
package checks

import "github.com/0xdevren/netsentry/internal/model"

// Issue represents a single topology problem detected by a check.
type Issue struct {
	// Code is a short machine-readable identifier (e.g. "DUP-IP-001").
	Code string `json:"code"`
	// Severity is the risk level of the issue.
	Severity string `json:"severity"`
	// Message is the human-readable description.
	Message string `json:"message"`
	// DeviceID is the device involved (if applicable).
	DeviceID string `json:"device_id,omitempty"`
}

// Check is the interface all topology checks implement.
type Check interface {
	// Run evaluates the topology graph and returns any detected issues.
	Run(g *model.TopologyGraph) []Issue
}
