// Package cmd – validate command (clean, no unused imports)
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/0xdevren/netsentry/internal/config"
	"github.com/0xdevren/netsentry/internal/model"
	"github.com/0xdevren/netsentry/internal/parser"
	"github.com/0xdevren/netsentry/internal/policy"
	"github.com/0xdevren/netsentry/internal/report"
	"github.com/0xdevren/netsentry/internal/validator"
)

func newValidateCmd() *cobra.Command {
	var (
		configPath  string
		policyPath  string
		format      string
		outputPath  string
		strict      bool
		timeout     time.Duration
		concurrency int
	)

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate a device configuration against a policy",
		Long: `Validate parses a device configuration file and evaluates it against
the specified policy definition. Exit codes:

  0  All rules passed (fully compliant)
  1  Policy violations detected
  2  Execution error
  3  Invalid input
  4  Timeout`,
		Example: `  netsentry validate --config router.conf --policy baseline.yaml
  netsentry validate --config router.conf --policy baseline.yaml --format json --output report.json
  netsentry validate --config router.conf --policy baseline.yaml --strict --timeout 30s`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if timeout > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, timeout)
				defer cancel()
			}

			rawData, err := os.ReadFile(configPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: cannot read config %q: %v\n", configPath, err)
				os.Exit(3)
			}

			detector := config.NewDetector()
			deviceType := detector.Detect(rawData)

			device := model.Device{ID: configPath, Type: deviceType}
			parsedCfg, err := parser.Parse(ctx, deviceType, rawData, device)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: parse failed: %v\n", err)
				os.Exit(3)
			}

			loader := policy.NewLoader()
			pol, err := loader.LoadFile(policyPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: policy load failed: %v\n", err)
				os.Exit(3)
			}

			if concurrency <= 0 {
				concurrency = 4
			}
			rep, err := validator.Validate(ctx, validator.ValidationRequest{
				Config:      parsedCfg,
				Policy:      pol,
				Strict:      strict,
				Concurrency: concurrency,
			})
			if err != nil {
				if ctx.Err() != nil {
					fmt.Fprintln(os.Stderr, "error: validation timed out")
					os.Exit(4)
				}
				fmt.Fprintf(os.Stderr, "error: validation failed: %v\n", err)
				os.Exit(2)
			}

			reporter, err := report.New(report.Options{
				Format:     report.Format(format),
				OutputPath: outputPath,
			})
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: reporter: %v\n", err)
				os.Exit(2)
			}
			if err := reporter.Generate(rep); err != nil {
				fmt.Fprintf(os.Stderr, "error: report generation: %v\n", err)
				os.Exit(2)
			}

			os.Exit(validator.ExitCode(rep, strict))
			return nil
		},
	}

	cmd.Flags().StringVar(&configPath, "config", "", "Path to device configuration file (required)")
	cmd.Flags().StringVar(&policyPath, "policy", "", "Path to policy YAML file (required)")
	cmd.Flags().StringVar(&format, "format", "table", "Output format: table|json|yaml|html")
	cmd.Flags().StringVar(&outputPath, "output", "", "Write report to file path")
	cmd.Flags().BoolVar(&strict, "strict", false, "Treat warnings as violations (exit 1)")
	cmd.Flags().DurationVar(&timeout, "timeout", 0, "Validation timeout (e.g. 30s)")
	cmd.Flags().IntVar(&concurrency, "concurrency", 4, "Parallel rule evaluation workers")
	_ = cmd.MarkFlagRequired("config")
	_ = cmd.MarkFlagRequired("policy")
	return cmd
}
