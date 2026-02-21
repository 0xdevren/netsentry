// Package policy defines the severity levels used across the NetSentry
// policy and validation subsystems.
package policy

// Severity represents the risk level of a policy violation.
type Severity string

const (
	// SeverityCritical indicates an immediate security risk requiring urgent remediation.
	SeverityCritical Severity = "CRITICAL"
	// SeverityHigh indicates a significant risk with high remediation priority.
	SeverityHigh Severity = "HIGH"
	// SeverityMedium indicates a moderate risk requiring scheduled remediation.
	SeverityMedium Severity = "MEDIUM"
	// SeverityLow indicates a minor risk or best-practice deviation.
	SeverityLow Severity = "LOW"
	// SeverityInfo indicates an informational finding with no immediate risk.
	SeverityInfo Severity = "INFO"
)

// Weight returns the numeric weight of the severity for scoring purposes.
// Higher values indicate greater severity.
func (s Severity) Weight() int {
	switch s {
	case SeverityCritical:
		return 100
	case SeverityHigh:
		return 75
	case SeverityMedium:
		return 50
	case SeverityLow:
		return 25
	case SeverityInfo:
		return 5
	default:
		return 0
	}
}

// String returns the string representation of the severity.
func (s Severity) String() string {
	return string(s)
}

// IsValid reports whether the severity value is a recognised level.
func (s Severity) IsValid() bool {
	switch s {
	case SeverityCritical, SeverityHigh, SeverityMedium, SeverityLow, SeverityInfo:
		return true
	default:
		return false
	}
}

// ParseSeverity converts a string to a Severity, returning SeverityInfo for
// unrecognised values.
func ParseSeverity(s string) Severity {
	sv := Severity(s)
	if sv.IsValid() {
		return sv
	}
	return SeverityInfo
}
