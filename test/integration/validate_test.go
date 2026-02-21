package integration

import (
	"context"
	"testing"

	"github.com/0xdevren/netsentry/internal/app"
	"github.com/0xdevren/netsentry/internal/policy"
	_ "github.com/0xdevren/netsentry/internal/parser/cisco"
)

func TestHardenedValidationPipeline(t *testing.T) {
	ctx := context.Background()
	opts := app.ValidateCommandOptions{
		ConfigPath:  "../data/test_router.conf",
		PolicyPath:  "../data/test_policy.yaml",
		Concurrency: 4,
	}

	orchestrator := app.NewOrchestrator(app.NewDefaultContext("debug", false))

	t.Run("execute_full_validation", func(t *testing.T) {
		report, exitCode, err := orchestrator.RunValidate(ctx, opts)
		if err != nil {
			t.Fatalf("Validation pipeline failed unexpectedly: %v", err)
		}

		if report == nil {
			t.Fatal("Expected validation report, got nil")
		}

		// Based on test_router.conf and test_policy.yaml:
		// ROUTE-OSPF-PASSIVE-DEFAULT -> should PASS
		// SEC-SSH-V2-ONLY -> should PASS
		// SEC-SNMP-V3-ONLY -> should FAIL (because 'snmp-server community public ro' is present)
		if report.Summary.Total != 3 {
			t.Errorf("Expected 3 total rules evaluated, got %d", report.Summary.Total)
		}
		if report.Summary.Passed != 3 {
			t.Errorf("Expected 3 passed rules, got %d", report.Summary.Passed)
		}
		if report.Summary.Failed != 0 {
			t.Errorf("Expected 0 failed rule, got %d", report.Summary.Failed)
		}

		// Exit code should be 0 since there are no failures
		if exitCode != 0 {
			t.Errorf("Expected exit code 0 due to all passing, got %d", exitCode)
		}

		// Verify SEC-SNMP-V3-ONLY passed
		foundSNMP := false
		for _, res := range report.Results {
			if res.RuleID == "SEC-SNMP-V3-ONLY" {
				if res.Status != policy.StatusPass {
					t.Errorf("Expected SEC-SNMP-V3-ONLY to PASS, but status was %s", res.Status)
				}
				foundSNMP = true
			}
		}
		if !foundSNMP {
			t.Error("Expected to find result for SEC-SNMP-V3-ONLY")
		}
	})
}
