package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/0xdevren/netsentry/internal/config"
	"github.com/0xdevren/netsentry/internal/model"
	"github.com/0xdevren/netsentry/internal/parser"
	"github.com/0xdevren/netsentry/internal/topology"
)

func newTopologyCmd() *cobra.Command {
	var configs []string

	cmd := &cobra.Command{
		Use:   "topology",
		Short: "Build and analyze the network topology from device configurations",
		Example: `  netsentry topology --config r1.conf --config r2.conf --config r3.conf`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(configs) == 0 {
				return fmt.Errorf("at least one --config is required")
			}

			detector := config.NewDetector()
			var parsedConfigs []*model.ConfigModel

			for _, cfgPath := range configs {
				data, err := os.ReadFile(cfgPath)
				if err != nil {
					return fmt.Errorf("cannot read %q: %w", cfgPath, err)
				}
				deviceType := detector.Detect(data)
				device := model.Device{ID: cfgPath, Type: deviceType}
				parsed, err := parser.Parse(cmd.Context(), deviceType, data, device)
				if err != nil {
					return fmt.Errorf("parse %q: %w", cfgPath, err)
				}
				parsedConfigs = append(parsedConfigs, parsed)
			}

			builder := topology.NewBuilder()
			graph := builder.Build(parsedConfigs)

			if err := graph.Validate(); err != nil {
				return fmt.Errorf("topology validation: %w", err)
			}

			analyzer := topology.NewAnalyzer()
			result := analyzer.Analyze(graph)

			fmt.Printf("Topology: %d devices, %d links\n", len(graph.Devices()), len(graph.Links()))
			fmt.Print(result.String())

			if result.HasIssues {
				os.Exit(1)
			}
			return nil
		},
	}

	cmd.Flags().StringArrayVar(&configs, "config", nil, "Device configuration file (repeatable)")
	return cmd
}
