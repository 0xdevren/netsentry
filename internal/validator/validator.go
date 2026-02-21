// Package validator implements the validation pipeline that orchestrates
// configuration loading, parsing, policy evaluation, and report generation.
package validator

import (
	"context"
	"fmt"

	"github.com/0xdevren/netsentry/internal/model"
	"github.com/0xdevren/netsentry/internal/policy"
)

// Validator is the top-level interface for device configuration validation.
type Validator interface {
	// Validate runs the full validation pipeline and returns a Report.
	Validate(ctx context.Context, req ValidationRequest) (*policy.Report, error)
}

// ValidationRequest encapsulates all parameters for a validation run.
type ValidationRequest struct {
	// Config is the parsed configuration model.
	Config *model.ConfigModel
	// Policy is the loaded policy to validate against.
	Policy *policy.Policy
	// Strict causes warnings to be treated as failures for exit-code purposes.
	Strict bool
	// Concurrency is the number of parallel rule evaluation workers.
	Concurrency int
}

// ExitCode computes the appropriate process exit code from a Report.
//
//	0 = fully compliant
//	1 = policy violations
//	2 = execution error
//	3 = invalid input
//	4 = timeout
func ExitCode(report *policy.Report, strict bool) int {
	if report == nil {
		return 2
	}
	s := report.Summary
	if s.Errors > 0 {
		return 2
	}
	if s.Failed > 0 {
		return 1
	}
	if strict && s.Warnings > 0 {
		return 1
	}
	return 0
}

// Validate is a convenience function that runs the pipeline using the default
// DeviceValidator implementation.
func Validate(ctx context.Context, req ValidationRequest) (*policy.Report, error) {
	v, err := NewDeviceValidator(DeviceValidatorOptions{
		Concurrency: req.Concurrency,
	})
	if err != nil {
		return nil, fmt.Errorf("validator: %w", err)
	}
	return v.Validate(ctx, req)
}
