// Package app defines the shared application context used across all commands
// and service handlers.
package app

import (
	"github.com/0xdevren/netsentry/internal/telemetry"
)

// Version, Commit, and BuildDate are injected at link time by the Makefile.
var (
	Version   = "dev"
	Commit    = "unknown"
	BuildDate = "unknown"
)

// Context holds all application-wide shared dependencies. It is constructed
// once at startup and passed down through the call chain. There is no global
// state; all fields are accessed via this struct.
type Context struct {
	// Logger is the application-wide structured logger.
	Logger *telemetry.Logger
	// Metrics holds the Prometheus metric instruments.
	Metrics *telemetry.Metrics
	// Config holds the application runtime configuration.
	Config RuntimeConfig
}

// RuntimeConfig holds configuration values parsed from flags and environment.
type RuntimeConfig struct {
	// LogLevel is the minimum log level.
	LogLevel string
	// LogJSON enables structured JSON logging.
	LogJSON bool
	// MetricsEnabled enables the Prometheus metrics endpoint.
	MetricsEnabled bool
	// TracingEnabled enables OpenTelemetry tracing.
	TracingEnabled bool
	// APIAddr is the listen address for the HTTP API server.
	APIAddr string
}

// NewContext constructs an application Context with the given configuration.
func NewContext(cfg RuntimeConfig) *Context {
	logger := telemetry.NewLogger(telemetry.LogOptions{
		Level: cfg.LogLevel,
		JSON:  cfg.LogJSON,
	})
	metrics := telemetry.NewMetrics("netsentry")
	return &Context{
		Logger:  logger,
		Metrics: metrics,
		Config:  cfg,
	}
}
