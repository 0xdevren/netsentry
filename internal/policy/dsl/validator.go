package dsl

import (
	"fmt"
)

// ValidationError describes a structural error in a raw policy document.
type ValidationError struct {
	// RuleIndex is the zero-based index of the offending rule (-1 for document-level errors).
	RuleIndex int
	// RuleID is the rule identifier if available.
	RuleID string
	// Field is the field name that contains the error.
	Field string
	// Message is the human-readable error description.
	Message string
}

func (e ValidationError) Error() string {
	if e.RuleID != "" {
		return fmt.Sprintf("rule %q field %q: %s", e.RuleID, e.Field, e.Message)
	}
	return fmt.Sprintf("document field %q: %s", e.Field, e.Message)
}

// Validator performs structural and semantic validation on a RawPolicy.
type Validator struct{}

// NewValidator constructs a new DSL Validator.
func NewValidator() *Validator {
	return &Validator{}
}

// Validate checks a RawPolicy for structural integrity and returns any errors.
func (v *Validator) Validate(p *RawPolicy) []ValidationError {
	var errs []ValidationError

	if p.Name == "" {
		errs = append(errs, ValidationError{RuleIndex: -1, Field: "name", Message: "policy name is required"})
	}

	seen := make(map[string]struct{}, len(p.Rules))
	for i, r := range p.Rules {
		if r.ID == "" {
			errs = append(errs, ValidationError{RuleIndex: i, Field: "id", Message: "rule id is required"})
			continue
		}
		if _, dup := seen[r.ID]; dup {
			errs = append(errs, ValidationError{RuleIndex: i, RuleID: r.ID, Field: "id", Message: "duplicate rule id"})
		}
		seen[r.ID] = struct{}{}

		if !containsKey(ValidSeverities, r.Severity) {
			errs = append(errs, ValidationError{
				RuleIndex: i, RuleID: r.ID, Field: "severity",
				Message: fmt.Sprintf("invalid severity %q; must be one of %v", r.Severity, ValidSeverities),
			})
		}

		if len(r.Match) == 0 {
			errs = append(errs, ValidationError{RuleIndex: i, RuleID: r.ID, Field: "match", Message: "match block is required"})
		} else {
			for k := range r.Match {
				if !containsKey(SupportedMatchKeys, k) {
					errs = append(errs, ValidationError{
						RuleIndex: i, RuleID: r.ID, Field: "match." + k,
						Message: fmt.Sprintf("unsupported match key %q", k),
					})
				}
			}
		}

		if len(r.Action) == 0 {
			errs = append(errs, ValidationError{RuleIndex: i, RuleID: r.ID, Field: "action", Message: "action block is required"})
		} else {
			for k := range r.Action {
				if !containsKey(SupportedActionKeys, k) {
					errs = append(errs, ValidationError{
						RuleIndex: i, RuleID: r.ID, Field: "action." + k,
						Message: fmt.Sprintf("unsupported action key %q", k),
					})
				}
			}
		}
	}

	return errs
}
