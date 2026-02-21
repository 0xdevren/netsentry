// Package api provides the HTTP REST API server for NetSentry.
package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/0xdevren/netsentry/internal/app"
	"github.com/0xdevren/netsentry/internal/telemetry"
)

// Server is the NetSentry HTTP API server.
type Server struct {
	appCtx  *app.Context
	httpSrv *http.Server
}

// NewServer constructs a Server bound to addr.
func NewServer(appCtx *app.Context, addr string) *Server {
	r := NewRouter(appCtx)
	return &Server{
		appCtx: appCtx,
		httpSrv: &http.Server{
			Addr:         addr,
			Handler:      r,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}

// Start begins listening for HTTP connections. It blocks until the server is
// shut down.
func (s *Server) Start() error {
	s.appCtx.Logger.Info("API server listening", "addr", s.httpSrv.Addr)
	if err := s.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("api server: %w", err)
	}
	return nil
}

// Shutdown gracefully drains active connections within the given timeout.
func (s *Server) Shutdown(ctx context.Context) error {
	s.appCtx.Logger.Info("shutting down API server")
	return s.httpSrv.Shutdown(ctx)
}

// NewRouter assembles the chi router with all API routes and middleware.
func NewRouter(appCtx *app.Context) http.Handler {
	r := chi.NewRouter()
	r.Use(RequestIDMiddleware)
	r.Use(LoggingMiddleware(appCtx.Logger))
	r.Use(RecoveryMiddleware(appCtx.Logger))

	r.Get("/healthz", HealthHandler(appCtx))
	r.Get("/readyz", ReadyHandler())
	r.Handle("/metrics", telemetry.MetricsHandler())

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/validate", ValidateHandler(appCtx))
		r.Get("/policy", PolicyListHandler(appCtx))
		r.Post("/policy/lint", PolicyLintHandler(appCtx))
		r.Get("/drift/{deviceID}", DriftHandler(appCtx))
	})

	return r
}
