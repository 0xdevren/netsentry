package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/0xdevren/netsentry/internal/policy"
	"github.com/0xdevren/netsentry/internal/policy/dsl"
)

func newPolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "policy",
		Short: "Manage and inspect policy definitions",
	}
	cmd.AddCommand(newPolicyListCmd(), newPolicyValidateCmd(), newPolicyLintCmd())
	return cmd
}

func newPolicyListCmd() *cobra.Command {
	var dir string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List policy files in a directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := os.ReadDir(dir)
			if err != nil {
				return fmt.Errorf("cannot read directory %q: %w", dir, err)
			}
			fmt.Printf("Policies in %s:\n", dir)
			for _, e := range entries {
				if !e.IsDir() {
					fmt.Printf("  %s\n", e.Name())
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&dir, "dir", "policies", "Directory to search for policy files")
	return cmd
}

func newPolicyValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate <policy.yaml>",
		Short: "Validate a policy file against the schema",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			loader := policy.NewLoader()
			pol, err := loader.LoadFile(args[0])
			if err != nil {
				return fmt.Errorf("policy invalid: %w", err)
			}
			fmt.Printf("Policy %q is valid (%d rules defined).\n", pol.Name, len(pol.Rules))
			return nil
		},
	}
}

func newPolicyLintCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "lint <policy.yaml>",
		Short: "Lint a policy file for structural and semantic errors",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			p := dsl.NewParser()
			v := dsl.NewValidator()

			raw, err := p.ParseFile(args[0])
			if err != nil {
				return fmt.Errorf("parse error: %w", err)
			}

			errs := v.Validate(raw)
			if len(errs) == 0 {
				fmt.Println("No issues found.")
				return nil
			}
			fmt.Fprintf(os.Stderr, "%d issue(s) found:\n", len(errs))
			for _, e := range errs {
				fmt.Fprintf(os.Stderr, "  [%s] %s\n", e.Field, e.Message)
			}
			os.Exit(1)
			return nil
		},
	}
}
