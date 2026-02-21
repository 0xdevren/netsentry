package policy

import (
	"fmt"

	"github.com/0xdevren/netsentry/internal/model"
)

// Evaluator applies a single Rule to a ConfigModel and returns the result.
type Evaluator struct {
	matcher *Matcher
}

// NewEvaluator constructs an Evaluator with the given Matcher.
func NewEvaluator(m *Matcher) *Evaluator {
	return &Evaluator{matcher: m}
}

// Evaluate applies the rule to the configuration and returns a ValidationResult.
func (e *Evaluator) Evaluate(rule Rule, cfg *model.ConfigModel) ValidationResult {
	result := ValidationResult{
		RuleID:          rule.ID,
		RuleDescription: rule.Description,
		Device:          cfg.Device,
		Severity:        rule.Severity,
	}

	if !rule.IsEnabled() {
		result.Status = StatusSkip
		result.Message = "rule is disabled"
		return result
	}

	matched, err := e.matcher.Match(rule.Match, cfg)
	if err != nil {
		result.Status = StatusError
		result.Message = fmt.Sprintf("evaluation error: %s", err.Error())
		return result
	}

	switch {
	case rule.Action.Deny && matched:
		result.Status = StatusFail
		result.Message = fmt.Sprintf("rule %s violated: %s", rule.ID, rule.Description)
		result.Remediation = rule.Action.Remediation
	case rule.Action.Warn && matched:
		result.Status = StatusWarn
		result.Message = fmt.Sprintf("rule %s warning: %s", rule.ID, rule.Description)
		result.Remediation = rule.Action.Remediation
	case !rule.Action.Deny && !rule.Action.Warn && !matched:
		// required_block / required match: absence is a failure
		result.Status = StatusFail
		result.Message = fmt.Sprintf("rule %s: required condition not met: %s", rule.ID, rule.Description)
		result.Remediation = rule.Action.Remediation
	default:
		result.Status = StatusPass
		result.Message = fmt.Sprintf("rule %s passed", rule.ID)
	}

	return result
}
