package policy

import "github.com/0xdevren/netsentry/internal/model"

// MatchType enumerates the supported rule matching strategies.
type MatchType string

const (
	// MatchContains passes when any configuration line contains the given substring.
	MatchContains MatchType = "contains"
	// MatchNotContains passes when no configuration line contains the given substring.
	MatchNotContains MatchType = "not_contains"
	// MatchRegex passes when any configuration line matches the given regular expression.
	MatchRegex MatchType = "regex"
	// MatchRequiredBlock passes when the configuration contains a required block by prefix.
	MatchRequiredBlock MatchType = "required_block"
)

// MatchSpec defines how a rule evaluates the device configuration.
type MatchSpec struct {
	// Contains is the substring to search for. Used with MatchContains.
	Contains string `json:"contains,omitempty" yaml:"contains,omitempty"`
	// NotContains is the substring that must be absent. Used with MatchNotContains.
	NotContains string `json:"not_contains,omitempty" yaml:"not_contains,omitempty"`
	// Regex is the regular expression pattern. Used with MatchRegex.
	Regex string `json:"regex,omitempty" yaml:"regex,omitempty"`
	// RequiredBlock is the configuration block prefix that must be present.
	RequiredBlock string `json:"required_block,omitempty" yaml:"required_block,omitempty"`
}

// ActionSpec defines the action to take when a rule matches.
type ActionSpec struct {
	// Deny causes the rule to produce a FAIL result when the condition is met.
	Deny bool `json:"deny,omitempty" yaml:"deny,omitempty"`
	// Warn causes the rule to produce a WARN result when the condition is met.
	Warn bool `json:"warn,omitempty" yaml:"warn,omitempty"`
	// Remediation provides a suggested corrective action message.
	Remediation string `json:"remediation,omitempty" yaml:"remediation,omitempty"`
}

// Rule defines a single compliance rule within a policy.
type Rule struct {
	// ID is the unique identifier for this rule (e.g. "SNMP-002").
	ID string `json:"id" yaml:"id"`
	// Description explains what the rule validates.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// Severity is the risk level of a violation of this rule.
	Severity Severity `json:"severity" yaml:"severity"`
	// Tags is an optional set of labels for categorisation and filtering.
	Tags []string `json:"tags,omitempty" yaml:"tags,omitempty"`
	// Match defines the condition evaluated against the device configuration.
	Match MatchSpec `json:"match" yaml:"match"`
	// Action defines the outcome when the condition is satisfied.
	Action ActionSpec `json:"action" yaml:"action"`
	// Enabled can disable a rule without removing it from the policy file.
	Enabled *bool `json:"enabled,omitempty" yaml:"enabled,omitempty"`
}

// IsEnabled reports whether the rule is active.
func (r Rule) IsEnabled() bool {
	if r.Enabled == nil {
		return true
	}
	return *r.Enabled
}

// Policy is the top-level policy definition loaded from a YAML file.
type Policy struct {
	// Name is the human-readable policy identifier (e.g. "CIS-Baseline").
	Name string `json:"name" yaml:"name"`
	// Version is the semantic version of the policy file.
	Version string `json:"version,omitempty" yaml:"version,omitempty"`
	// Description explains the purpose of the policy.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// Author is the policy author or organisation.
	Author string `json:"author,omitempty" yaml:"author,omitempty"`
	// Rules is the ordered list of compliance rules.
	Rules []Rule `json:"rules" yaml:"rules"`
}

// ValidationStatus represents the outcome of a single rule evaluation.
type ValidationStatus string

const (
	// StatusPass indicates the device is compliant with the rule.
	StatusPass ValidationStatus = "PASS"
	// StatusFail indicates the device violates the rule.
	StatusFail ValidationStatus = "FAIL"
	// StatusWarn indicates a non-critical finding for the rule.
	StatusWarn ValidationStatus = "WARN"
	// StatusSkip indicates the rule was not evaluated (e.g. disabled).
	StatusSkip ValidationStatus = "SKIP"
	// StatusError indicates an internal error occurred during evaluation.
	StatusError ValidationStatus = "ERROR"
)

// ValidationResult is the outcome of evaluating a single rule against a single device.
type ValidationResult struct {
	// RuleID is the ID of the evaluated rule.
	RuleID string `json:"rule_id" yaml:"rule_id"`
	// RuleDescription is the human-readable rule description.
	RuleDescription string `json:"rule_description,omitempty" yaml:"rule_description,omitempty"`
	// Device is the device that was evaluated.
	Device model.Device `json:"device" yaml:"device"`
	// Status is the evaluation outcome.
	Status ValidationStatus `json:"status" yaml:"status"`
	// Severity is the severity of the rule.
	Severity Severity `json:"severity" yaml:"severity"`
	// Message is a human-readable explanation of the result.
	Message string `json:"message,omitempty" yaml:"message,omitempty"`
	// Remediation is a suggested corrective action (populated on FAIL/WARN).
	Remediation string `json:"remediation,omitempty" yaml:"remediation,omitempty"`
}

// Report is the top-level output of a validation run against a device.
type Report struct {
	// Device is the evaluated device.
	Device model.Device `json:"device" yaml:"device"`
	// Policy is the policy name used for this validation.
	Policy string `json:"policy" yaml:"policy"`
	// PolicyVersion is the version of the evaluated policy.
	PolicyVersion string `json:"policy_version,omitempty" yaml:"policy_version,omitempty"`
	// Results is the complete list of rule evaluation results.
	Results []ValidationResult `json:"results" yaml:"results"`
	// Summary provides aggregate compliance metrics.
	Summary ReportSummary `json:"summary" yaml:"summary"`
}

// ReportSummary aggregates the compliance metrics for a validation report.
type ReportSummary struct {
	// Total is the number of rules evaluated.
	Total int `json:"total" yaml:"total"`
	// Passed is the number of rules that passed.
	Passed int `json:"passed" yaml:"passed"`
	// Failed is the number of rules that failed.
	Failed int `json:"failed" yaml:"failed"`
	// Warnings is the number of rules that produced a warning.
	Warnings int `json:"warnings" yaml:"warnings"`
	// Skipped is the number of rules that were skipped.
	Skipped int `json:"skipped" yaml:"skipped"`
	// Errors is the number of rules that encountered an internal error.
	Errors int `json:"errors" yaml:"errors"`
	// Score is the compliance percentage (0-100).
	Score float64 `json:"score" yaml:"score"`
}

// ComputeSummary calculates aggregate statistics from a slice of results.
func ComputeSummary(results []ValidationResult) ReportSummary {
	s := ReportSummary{Total: len(results)}
	for _, r := range results {
		switch r.Status {
		case StatusPass:
			s.Passed++
		case StatusFail:
			s.Failed++
		case StatusWarn:
			s.Warnings++
		case StatusSkip:
			s.Skipped++
		case StatusError:
			s.Errors++
		}
	}
	evaluated := s.Passed + s.Failed + s.Warnings
	if evaluated > 0 {
		s.Score = float64(s.Passed) / float64(evaluated) * 100
	}
	return s
}
