package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/0xdevren/netsentry/internal/config"
	"github.com/0xdevren/netsentry/internal/drift"
	"github.com/0xdevren/netsentry/internal/model"
	"github.com/0xdevren/netsentry/internal/parser"
)

func newDriftCmd() *cobra.Command {
	var (
		baselinePath string
		currentPath  string
		threshold    float64
	)

	cmd := &cobra.Command{
		Use:   "drift",
		Short: "Detect configuration drift between two snapshots",
		Example: `  netsentry drift --baseline router-2024-01-01.conf --current router.conf
  netsentry drift --baseline baseline.conf --current current.conf --threshold 10`,
		RunE: func(cmd *cobra.Command, args []string) error {
			baselineData, err := os.ReadFile(baselinePath)
			if err != nil {
				return fmt.Errorf("cannot read baseline %q: %w", baselinePath, err)
			}
			currentData, err := os.ReadFile(currentPath)
			if err != nil {
				return fmt.Errorf("cannot read current %q: %w", currentPath, err)
			}

			// Hash comparison.
			bHash := drift.HashConfig("baseline", baselineData)
			cHash := drift.HashConfig("current", currentData)
			if !drift.HasChanged(bHash, cHash) {
				fmt.Println("No configuration drift detected.")
				return nil
			}

			// Line diff.
			comparator := drift.NewComparator()
			diff := comparator.Compare("device", baselineData, currentData)

			// Score.
			scorer := drift.NewDriftScorer(threshold)
			baselineLines := splitData(baselineData)
			score := scorer.Score(diff, len(baselineLines))

			// Detect device type for context.
			detector := config.NewDetector()
			deviceType := detector.Detect(currentData)
			_, _ = parser.Parse(cmd.Context(), deviceType, currentData, model.Device{})

			fmt.Printf("Configuration drift detected for device.\n\n")
			fmt.Printf("Lines added   : %d\n", score.LinesAdded)
			fmt.Printf("Lines removed : %d\n", score.LinesRemoved)
			fmt.Printf("Total changes : %d\n", score.TotalChanges)
			fmt.Printf("Drift percent : %.1f%%\n", score.DriftPercent)
			if score.Significant {
				fmt.Printf("Status        : SIGNIFICANT (threshold: %.1f%%)\n", threshold)
			} else {
				fmt.Printf("Status        : within threshold (%.1f%%)\n", threshold)
			}

			fmt.Println("\nDiff:")
			fmt.Println(diff.String())

			if score.Significant {
				os.Exit(1)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&baselinePath, "baseline", "", "Path to baseline configuration file (required)")
	cmd.Flags().StringVar(&currentPath, "current", "", "Path to current configuration file (required)")
	cmd.Flags().Float64Var(&threshold, "threshold", 5.0, "Drift percentage threshold for significance")
	_ = cmd.MarkFlagRequired("baseline")
	_ = cmd.MarkFlagRequired("current")
	return cmd
}

func splitData(data []byte) []string {
	var lines []string
	line := ""
	for _, c := range string(data) {
		if c == '\n' {
			if line != "" {
				lines = append(lines, line)
			}
			line = ""
		} else {
			line += string(c)
		}
	}
	return lines
}
