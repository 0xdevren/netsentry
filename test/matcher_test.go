package netsentry_test

import (
	"testing"

	"github.com/0xdevren/netsentry/internal/model"
	"github.com/0xdevren/netsentry/internal/policy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeConfig(lines ...string) *model.ConfigModel {
	return &model.ConfigModel{
		Device: model.Device{ID: "R1", Hostname: "R1"},
		Lines:  lines,
	}
}

func TestMatcher_Contains_Match(t *testing.T) {
	m := policy.NewMatcher()
	cfg := makeConfig("snmp-server community public RO")
	matched, err := m.Match(policy.MatchSpec{Contains: "snmp-server community public"}, cfg)
	require.NoError(t, err)
	assert.True(t, matched)
}

func TestMatcher_Contains_NoMatch(t *testing.T) {
	m := policy.NewMatcher()
	cfg := makeConfig("hostname R1", "ip ssh version 2")
	matched, err := m.Match(policy.MatchSpec{Contains: "snmp-server community public"}, cfg)
	require.NoError(t, err)
	assert.False(t, matched)
}

func TestMatcher_NotContains_Pass(t *testing.T) {
	m := policy.NewMatcher()
	cfg := makeConfig("hostname R1", "ip ssh version 2")
	matched, err := m.Match(policy.MatchSpec{NotContains: "snmp-server community public"}, cfg)
	require.NoError(t, err)
	assert.True(t, matched)
}

func TestMatcher_NotContains_Fail(t *testing.T) {
	m := policy.NewMatcher()
	cfg := makeConfig("snmp-server community public RO")
	matched, err := m.Match(policy.MatchSpec{NotContains: "snmp-server community public"}, cfg)
	require.NoError(t, err)
	assert.False(t, matched)
}

func TestMatcher_Regex_Match(t *testing.T) {
	m := policy.NewMatcher()
	cfg := makeConfig("enable password cisco123")
	matched, err := m.Match(policy.MatchSpec{Regex: `^enable password`}, cfg)
	require.NoError(t, err)
	assert.True(t, matched)
}

func TestMatcher_Regex_InvalidPattern(t *testing.T) {
	m := policy.NewMatcher()
	cfg := makeConfig("hostname R1")
	_, err := m.Match(policy.MatchSpec{Regex: `[invalid`}, cfg)
	assert.Error(t, err)
}

func TestMatcher_RequiredBlock_Present(t *testing.T) {
	m := policy.NewMatcher()
	cfg := makeConfig("ntp server 10.0.0.1", "hostname R1")
	matched, err := m.Match(policy.MatchSpec{RequiredBlock: "ntp server"}, cfg)
	require.NoError(t, err)
	assert.True(t, matched)
}

func TestMatcher_RequiredBlock_Absent(t *testing.T) {
	m := policy.NewMatcher()
	cfg := makeConfig("hostname R1", "ip route 0.0.0.0 0.0.0.0 10.0.0.1")
	matched, err := m.Match(policy.MatchSpec{RequiredBlock: "ntp server"}, cfg)
	require.NoError(t, err)
	assert.False(t, matched)
}

func TestMatcher_EmptySpec_Error(t *testing.T) {
	m := policy.NewMatcher()
	cfg := makeConfig("hostname R1")
	_, err := m.Match(policy.MatchSpec{}, cfg)
	assert.Error(t, err)
}

func TestEvaluator_DenyOnMatch(t *testing.T) {
	enabled := true
	rule := policy.Rule{
		ID:       "SNMP-001",
		Severity: policy.SeverityHigh,
		Match:    policy.MatchSpec{Contains: "snmp-server community public"},
		Action:   policy.ActionSpec{Deny: true},
		Enabled:  &enabled,
	}
	evaluator := policy.NewEvaluator(policy.NewMatcher())
	cfg := makeConfig("snmp-server community public RO")
	result := evaluator.Evaluate(rule, cfg)
	assert.Equal(t, policy.StatusFail, result.Status)
}

func TestEvaluator_PassWhenNoMatch(t *testing.T) {
	enabled := true
	rule := policy.Rule{
		ID:       "SNMP-001",
		Severity: policy.SeverityHigh,
		Match:    policy.MatchSpec{Contains: "snmp-server community public"},
		Action:   policy.ActionSpec{Deny: true},
		Enabled:  &enabled,
	}
	evaluator := policy.NewEvaluator(policy.NewMatcher())
	cfg := makeConfig("hostname R1", "ip ssh version 2")
	result := evaluator.Evaluate(rule, cfg)
	assert.Equal(t, policy.StatusPass, result.Status)
}

func TestEvaluator_SkipsDisabledRule(t *testing.T) {
	disabled := false
	rule := policy.Rule{
		ID:       "SNMP-001",
		Severity: policy.SeverityHigh,
		Match:    policy.MatchSpec{Contains: "anything"},
		Action:   policy.ActionSpec{Deny: true},
		Enabled:  &disabled,
	}
	evaluator := policy.NewEvaluator(policy.NewMatcher())
	cfg := makeConfig("anything")
	result := evaluator.Evaluate(rule, cfg)
	assert.Equal(t, policy.StatusSkip, result.Status)
}
