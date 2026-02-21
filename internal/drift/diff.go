package drift

import (
	"fmt"
	"strings"
)

// ChangeType classifies the nature of a configuration line change.
type ChangeType string

const (
	ChangeAdded   ChangeType = "ADDED"
	ChangeRemoved ChangeType = "REMOVED"
)

// LineDiff represents a single line-level change between two configurations.
type LineDiff struct {
	// Type indicates whether the line was added or removed.
	Type ChangeType `json:"type"`
	// Line is the configuration line content.
	Line string `json:"line"`
	// LineNumber is the approximate line number in the new (ADDED) or old (REMOVED) config.
	LineNumber int `json:"line_number"`
}

// DiffResult is the full set of differences between two raw configurations.
type DiffResult struct {
	// DeviceID is the device these configs belong to.
	DeviceID string `json:"device_id"`
	// Added is the set of lines present in current but absent in baseline.
	Added []LineDiff `json:"added,omitempty"`
	// Removed is the set of lines present in baseline but absent in current.
	Removed []LineDiff `json:"removed,omitempty"`
	// HasChanges indicates whether any differences were detected.
	HasChanges bool `json:"has_changes"`
}

// String returns a unified-diff-style human-readable representation.
func (d *DiffResult) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "--- baseline/%s\n+++ current/%s\n", d.DeviceID, d.DeviceID)
	for _, r := range d.Removed {
		fmt.Fprintf(&sb, "-%s\n", r.Line)
	}
	for _, a := range d.Added {
		fmt.Fprintf(&sb, "+%s\n", a.Line)
	}
	return sb.String()
}
