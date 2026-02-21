package validator

import (
	"context"
	"fmt"
	"sort"

	"github.com/0xdevren/netsentry/internal/policy"
)

// DeviceValidatorOptions configures the DeviceValidator.
type DeviceValidatorOptions struct {
	// Concurrency is the number of parallel evaluation workers.
	Concurrency int
}

// DeviceValidator validates a single device configuration against a policy.
type DeviceValidator struct {
	engine *policy.Engine
}

// NewDeviceValidator constructs a DeviceValidator with the given options.
func NewDeviceValidator(opts DeviceValidatorOptions) (*DeviceValidator, error) {
	concurrency := opts.Concurrency
	if concurrency <= 0 {
		concurrency = 4
	}
	return &DeviceValidator{
		engine: policy.NewEngine(policy.EngineOptions{Concurrency: concurrency}),
	}, nil
}

// Validate runs the policy engine against the config model and assembles
// a structured Report.
func (v *DeviceValidator) Validate(ctx context.Context, req ValidationRequest) (*policy.Report, error) {
	if req.Config == nil {
		return nil, fmt.Errorf("device validator: ConfigModel is required")
	}
	if req.Policy == nil {
		return nil, fmt.Errorf("device validator: Policy is required")
	}

	results, err := v.engine.Run(ctx, req.Policy, req.Config)
	if err != nil {
		return nil, fmt.Errorf("device validator: engine run: %w", err)
	}

	// Sort results by rule ID for deterministic output.
	sort.Slice(results, func(i, j int) bool {
		return results[i].RuleID < results[j].RuleID
	})

	summary := policy.ComputeSummary(results)

	report := &policy.Report{
		Device:        req.Config.Device,
		Policy:        req.Policy.Name,
		PolicyVersion: req.Policy.Version,
		Results:       results,
		Summary:       summary,
	}

	return report, nil
}
