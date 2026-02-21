package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/0xdevren/netsentry/internal/telemetry"
)

// Runtime manages the application lifecycle: startup, signal handling,
// and graceful shutdown.
type Runtime struct {
	appCtx       *Context
	shutdownHooks []func(context.Context) error
}

// NewRuntime constructs a Runtime from an application Context.
func NewRuntime(appCtx *Context) *Runtime {
	return &Runtime{appCtx: appCtx}
}

// RegisterShutdownHook registers a function to be called during graceful shutdown.
// Hooks are invoked in LIFO order.
func (r *Runtime) RegisterShutdownHook(fn func(context.Context) error) {
	r.shutdownHooks = append(r.shutdownHooks, fn)
}

// WaitForSignal blocks until SIGINT or SIGTERM is received, then invokes all
// registered shutdown hooks.
func (r *Runtime) WaitForSignal(ctx context.Context) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		r.appCtx.Logger.Info("signal received, shutting down", "signal", sig.String())
	case <-ctx.Done():
		r.appCtx.Logger.Info("context cancelled, shutting down")
	}

	r.runShutdownHooks(ctx)
}

// runShutdownHooks executes all registered hooks in reverse registration order.
func (r *Runtime) runShutdownHooks(ctx context.Context) {
	for i := len(r.shutdownHooks) - 1; i >= 0; i-- {
		if err := r.shutdownHooks[i](ctx); err != nil {
			r.appCtx.Logger.Error("shutdown hook error", err)
		}
	}
}

// BuildInfo returns a map of build-time metadata for display or logging.
func BuildInfo() map[string]string {
	return map[string]string{
		"version":    Version,
		"commit":     Commit,
		"build_date": BuildDate,
	}
}

// NewDefaultContext constructs an application Context suitable for CLI commands.
func NewDefaultContext(logLevel string, logJSON bool) *Context {
	return NewContext(RuntimeConfig{
		LogLevel: logLevel,
		LogJSON:  logJSON,
	})
}

// SetupTracing initialises the OTel tracer and registers its shutdown hook.
func SetupTracing(appCtx *Context, r *Runtime, serviceName string) {
	_, shutdown, err := telemetry.InitTracer(telemetry.TracerOptions{
		ServiceName:    serviceName,
		ServiceVersion: Version,
		Enabled:        appCtx.Config.TracingEnabled,
	})
	if err != nil {
		appCtx.Logger.Warn("tracing init failed", "error", err.Error())
		return
	}
	r.RegisterShutdownHook(shutdown)
}
