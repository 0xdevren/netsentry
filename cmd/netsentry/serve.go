package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/0xdevren/netsentry/internal/api"
	"github.com/0xdevren/netsentry/internal/app"
)

func newServeCmd() *cobra.Command {
	var addr string

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the NetSentry HTTP API server",
		Long: `Serve starts the NetSentry REST API on the specified address.

Endpoints:
  GET  /healthz              – liveness probe
  GET  /readyz               – readiness probe
  GET  /metrics              – Prometheus metrics
  POST /api/v1/validate      – validate a device configuration
  GET  /api/v1/policy        – list available policies
  POST /api/v1/policy/lint   – lint a policy definition
  GET  /api/v1/drift/{id}    – retrieve drift report`,
		Example: `  netsentry serve
  netsentry serve --addr 0.0.0.0:9090`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			srv := api.NewServer(appCtx, addr)
			runtime := app.NewRuntime(appCtx)
			runtime.RegisterShutdownHook(func(ctx context.Context) error {
				return srv.Shutdown(ctx)
			})

			errCh := make(chan error, 1)
			go func() { errCh <- srv.Start() }()

			select {
			case err := <-errCh:
				if err != nil {
					return fmt.Errorf("server error: %w", err)
				}
			case <-ctx.Done():
				runtime.WaitForSignal(ctx)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&addr, "addr", ":8080", "Listen address for the HTTP API server")
	return cmd
}
