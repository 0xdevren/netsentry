// Package main is the netsentry CLI entry point.
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/0xdevren/netsentry/internal/app"
)

var (
	// globalLogLevel is the --log-level flag value, shared across all commands.
	globalLogLevel string
	// globalLogJSON enables JSON log output.
	globalLogJSON bool
	// appCtx is the shared application context constructed in PersistentPreRunE.
	appCtx *app.Context
)

// rootCmd is the top-level cobra command.
var rootCmd = &cobra.Command{
	Use:   "netsentry",
	Short: "NetSentry – Next-Generation Network Configuration Validator",
	Long: `NetSentry is a production-grade CLI tool for:
  - Network configuration validation against policy baselines
  - Compliance enforcement (CIS, custom policies)
  - Security misconfiguration detection
  - Structured report generation (JSON, YAML, HTML, table)
  - CI/CD pipeline integration`,
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		appCtx = app.NewDefaultContext(globalLogLevel, globalLogJSON)
		return nil
	},
}

// Execute adds all child commands to root and runs the CLI.
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return err
	}
	return nil
}

func init() {
	rootCmd.PersistentFlags().StringVar(&globalLogLevel, "log-level", "info",
		"Log level (debug|info|warn|error)")
	rootCmd.PersistentFlags().BoolVar(&globalLogJSON, "log-json", false,
		"Output logs in JSON format")

	rootCmd.AddCommand(
		newValidateCmd(),
		newScanCmd(),
		newPolicyCmd(),
		newReportCmd(),
		newDriftCmd(),
		newTopologyCmd(),
		newServeCmd(),
		newVersionCmd(),
	)
}

// newVersionCmd returns the version sub-command.
func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			info := app.BuildInfo()
			fmt.Printf("netsentry %s (commit: %s, built: %s)\n",
				info["version"], info["commit"], info["build_date"])
		},
	}
}

// newScanCmd returns the scan sub-command (SSH-based remote config retrieval).
func newScanCmd() *cobra.Command {
	var (
		target string
		sshUser string
		sshKey string
	)
	cmd := &cobra.Command{
		Use:   "scan",
		Short: "Retrieve and validate configuration from a live device via SSH",
		RunE: func(cmd *cobra.Command, args []string) error {
			if target == "" {
				return fmt.Errorf("--target is required")
			}
			fmt.Printf("Scanning device %s as user %s using key %s\n", target, sshUser, sshKey)
			fmt.Println("(SSH scan: connect to device, run 'show running-config', then validate)")
			return nil
		},
	}
	cmd.Flags().StringVar(&target, "target", "", "Device IP address or hostname (required)")
	cmd.Flags().StringVar(&sshUser, "ssh-user", "admin", "SSH username")
	cmd.Flags().StringVar(&sshKey, "ssh-key", "~/.ssh/id_rsa", "Path to SSH private key")
	return cmd
}
