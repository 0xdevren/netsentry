package app

import (
	"context"
	"fmt"
	"time"

	"github.com/0xdevren/netsentry/internal/config"
	"github.com/0xdevren/netsentry/internal/model"
	"github.com/0xdevren/netsentry/internal/parser"
	"github.com/0xdevren/netsentry/internal/policy"
	"github.com/0xdevren/netsentry/internal/report"
	"github.com/0xdevren/netsentry/internal/validator"
	"github.com/prometheus/client_golang/prometheus"
)

// ValidateCommandOptions encapsulates all parameters for a validate command run.
type ValidateCommandOptions struct {
	// ConfigPath is the filesystem path to the device configuration file.
	ConfigPath string
	// PolicyPath is the filesystem path to the policy YAML file.
	PolicyPath string
	// Format is the output format ("table", "json", "yaml", "html").
	Format string
	// OutputPath writes the report to a file instead of stdout.
	OutputPath string
	// Strict treats warnings as failures for exit code purposes.
	Strict bool
	// Timeout is the maximum allowed duration for the validation run.
	Timeout time.Duration
	// Concurrency is the number of parallel rule evaluation workers.
	Concurrency int
}

// Orchestrator coordinates the high-level use cases for the application.
type Orchestrator struct {
	appCtx     *Context
	cfgLoader  *config.Loader
	detector   *config.Detector
	policyLoader *policy.Loader
}

// NewOrchestrator constructs an Orchestrator from the application context.
func NewOrchestrator(appCtx *Context) *Orchestrator {
	return &Orchestrator{
		appCtx:       appCtx,
		cfgLoader:    config.NewLoader(),
		detector:     config.NewDetector(),
		policyLoader: policy.NewLoader(),
	}
}

// RunValidate executes the full validation pipeline and returns the report and
// the process exit code.
func (o *Orchestrator) RunValidate(ctx context.Context, opts ValidateCommandOptions) (*policy.Report, int, error) {
	if opts.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
	}

	o.appCtx.Logger.Info("loading device configuration", "path", opts.ConfigPath)
	rawData, err := o.cfgLoader.Load(ctx, config.LoadOptions{
		Source: "filesystem",
		Path:   opts.ConfigPath,
	})
	if err != nil {
		return nil, 3, fmt.Errorf("orchestrator: load config: %w", err)
	}

	deviceType := o.detector.Detect(rawData)
	o.appCtx.Logger.Info("detected device type", "type", string(deviceType))

	o.appCtx.Logger.Info("parsing configuration")
	device := model.Device{ID: opts.ConfigPath, Type: deviceType}
	parsedCfg, err := parser.Parse(ctx, deviceType, rawData, device)
	if err != nil {
		return nil, 3, fmt.Errorf("orchestrator: parse config: %w", err)
	}

	o.appCtx.Logger.Info("loading policy", "path", opts.PolicyPath)
	pol, err := o.policyLoader.LoadFile(opts.PolicyPath)
	if err != nil {
		return nil, 3, fmt.Errorf("orchestrator: load policy: %w", err)
	}

	timer := prometheus.NewTimer(o.appCtx.Metrics.ValidationDuration)
	defer timer.ObserveDuration()

	o.appCtx.Metrics.ActiveValidations.Inc()
	defer o.appCtx.Metrics.ActiveValidations.Dec()

	rep, err := validator.Validate(ctx, validator.ValidationRequest{
		Config:      parsedCfg,
		Policy:      pol,
		Strict:      opts.Strict,
		Concurrency: opts.Concurrency,
	})
	if err != nil {
		if ctx.Err() != nil {
			return nil, 4, fmt.Errorf("orchestrator: validation timed out: %w", ctx.Err())
		}
		return nil, 2, fmt.Errorf("orchestrator: validation: %w", err)
	}

	o.appCtx.Metrics.ValidationTotal.Add(1)
	for _, res := range rep.Results {
		if res.Status == policy.StatusFail || res.Status == policy.StatusWarn {
			o.appCtx.Metrics.PolicyViolations.WithLabelValues(string(res.Severity)).Inc()
		}
	}

	o.appCtx.Logger.Info("validation complete",
		"passed", rep.Summary.Passed,
		"failed", rep.Summary.Failed,
		"score", fmt.Sprintf("%.0f%%", rep.Summary.Score),
	)

	exitCode := validator.ExitCode(rep, opts.Strict)

	reporter, err := report.New(report.Options{
		Format:     report.Format(opts.Format),
		OutputPath: opts.OutputPath,
	})
	if err != nil {
		return rep, 2, fmt.Errorf("orchestrator: reporter: %w", err)
	}
	if err := reporter.Generate(rep); err != nil {
		return rep, 2, fmt.Errorf("orchestrator: generate report: %w", err)
	}

	return rep, exitCode, nil
}


