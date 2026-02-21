package netsentry_test

import (
	"testing"

	"github.com/0xdevren/netsentry/internal/policy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoader_LoadBytes_Valid(t *testing.T) {
	yaml := `
name: test-policy
version: "1.0"
rules:
  - id: SNMP-001
    description: No public SNMP
    severity: HIGH
    match:
      contains: "snmp-server community public"
    action:
      deny: true
`
	loader := policy.NewLoader()
	pol, err := loader.LoadBytes([]byte(yaml))
	require.NoError(t, err)
	assert.Equal(t, "test-policy", pol.Name)
	assert.Len(t, pol.Rules, 1)
	assert.Equal(t, "SNMP-001", pol.Rules[0].ID)
	assert.Equal(t, policy.SeverityHigh, pol.Rules[0].Severity)
}

func TestLoader_LoadBytes_MissingName(t *testing.T) {
	yaml := `
rules:
  - id: X-001
    severity: LOW
    match:
      contains: "foo"
    action:
      deny: true
`
	loader := policy.NewLoader()
	_, err := loader.LoadBytes([]byte(yaml))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name is required")
}

func TestLoader_LoadBytes_DuplicateRuleID(t *testing.T) {
	yaml := `
name: dup-test
rules:
  - id: DUP-001
    severity: LOW
    match:
      contains: "foo"
    action:
      deny: true
  - id: DUP-001
    severity: HIGH
    match:
      contains: "bar"
    action:
      deny: true
`
	loader := policy.NewLoader()
	_, err := loader.LoadBytes([]byte(yaml))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate rule id")
}

func TestLoader_LoadBytes_InvalidSeverity(t *testing.T) {
	yaml := `
name: sev-test
rules:
  - id: SEV-001
    severity: EXTREME
    match:
      contains: "foo"
    action:
      deny: true
`
	loader := policy.NewLoader()
	_, err := loader.LoadBytes([]byte(yaml))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid severity")
}

func TestSeverity_Weight(t *testing.T) {
	tests := []struct {
		severity policy.Severity
		weight   int
	}{
		{policy.SeverityCritical, 100},
		{policy.SeverityHigh, 75},
		{policy.SeverityMedium, 50},
		{policy.SeverityLow, 25},
		{policy.SeverityInfo, 5},
	}
	for _, tt := range tests {
		t.Run(string(tt.severity), func(t *testing.T) {
			assert.Equal(t, tt.weight, tt.severity.Weight())
		})
	}
}

func TestComputeSummary(t *testing.T) {
	results := []policy.ValidationResult{
		{Status: policy.StatusPass},
		{Status: policy.StatusPass},
		{Status: policy.StatusFail},
		{Status: policy.StatusWarn},
		{Status: policy.StatusSkip},
	}
	s := policy.ComputeSummary(results)
	assert.Equal(t, 5, s.Total)
	assert.Equal(t, 2, s.Passed)
	assert.Equal(t, 1, s.Failed)
	assert.Equal(t, 1, s.Warnings)
	assert.Equal(t, 1, s.Skipped)
	// Score = 2 / (2+1+1) * 100 = 50
	assert.InDelta(t, 50.0, s.Score, 0.01)
}
