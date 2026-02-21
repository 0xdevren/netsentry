package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/0xdevren/netsentry/internal/config"
	"github.com/0xdevren/netsentry/internal/model"
	"github.com/0xdevren/netsentry/internal/parser"
	"github.com/0xdevren/netsentry/internal/policy"
	"github.com/0xdevren/netsentry/internal/report"
	"github.com/0xdevren/netsentry/internal/validator"
)

func newReportCmd() *cobra.Command {
	var (
		configPath string
		policyPath string
		format     string
		outputPath string
	)

	cmd := &cobra.Command{
		Use:   "report",
		Short: "Generate a compliance report without enforcing exit codes",
		Long: `Report runs the full validation pipeline and produces a report in the
specified format. Unlike validate, it always exits with code 0 unless an
execution error occurs, making it suitable for report-only pipelines.`,
		Example: `  netsentry report --config router.conf --policy baseline.yaml --format html --output report.html
  netsentry report --config router.conf --policy baseline.yaml --format json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			rawData, err := os.ReadFile(configPath)
			if err != nil {
				return fmt.Errorf("cannot read config %q: %w", configPath, err)
			}

			detector := config.NewDetector()
			deviceType := detector.Detect(rawData)
			device := model.Device{ID: configPath, Type: deviceType}

			parsedCfg, err := parser.Parse(cmd.Context(), deviceType, rawData, device)
			if err != nil {
				return fmt.Errorf("parse error: %w", err)
			}

			loader := policy.NewLoader()
			pol, err := loader.LoadFile(policyPath)
			if err != nil {
				return fmt.Errorf("policy error: %w", err)
			}

			rep, err := validator.Validate(cmd.Context(), validator.ValidationRequest{
				Config: parsedCfg, Policy: pol, Concurrency: 4,
			})
			if err != nil {
				return fmt.Errorf("validation error: %w", err)
			}

			reporter, err := report.New(report.Options{
				Format:     report.Format(format),
				OutputPath: outputPath,
			})
			if err != nil {
				return fmt.Errorf("reporter error: %w", err)
			}
			return reporter.Generate(rep)
		},
	}

	cmd.Flags().StringVar(&configPath, "config", "", "Path to device configuration file (required)")
	cmd.Flags().StringVar(&policyPath, "policy", "", "Path to policy YAML file (required)")
	cmd.Flags().StringVar(&format, "format", "table", "Output format: table|json|yaml|html")
	cmd.Flags().StringVar(&outputPath, "output", "", "Write report to file path")
	_ = cmd.MarkFlagRequired("config")
	_ = cmd.MarkFlagRequired("policy")
	return cmd
}
