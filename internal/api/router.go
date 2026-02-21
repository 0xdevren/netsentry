package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/0xdevren/netsentry/internal/app"
	"github.com/0xdevren/netsentry/internal/telemetry"
)

// router.go delegates to server.go's NewRouter. This file provides the explicit
// route table as documentation.
//
// Route map:
//   GET  /healthz               – liveness probe
//   GET  /readyz                – readiness probe
//   GET  /metrics               – Prometheus metrics
//   POST /api/v1/validate       – run a validation job
//   GET  /api/v1/policy         – list available policies
//   POST /api/v1/policy/lint    – lint a policy file
//   GET  /api/v1/drift/{id}     – retrieve drift report for a device

// BuildRouter is the canonical router constructor used by tests and main server.
func BuildRouter(appCtx *app.Context) http.Handler {
	r := chi.NewRouter()
	r.Use(RequestIDMiddleware)
	r.Use(LoggingMiddleware(appCtx.Logger))
	r.Use(RecoveryMiddleware(appCtx.Logger))
	r.Use(corsMiddleware)

	r.Get("/healthz", HealthHandler(appCtx))
	r.Get("/readyz", ReadyHandler())
	r.Handle("/metrics", telemetry.MetricsHandler())

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(JSONContentType)
		r.Post("/validate", ValidateHandler(appCtx))
		r.Get("/policy", PolicyListHandler(appCtx))
		r.Post("/policy/lint", PolicyLintHandler(appCtx))
		r.Get("/drift/{deviceID}", DriftHandler(appCtx))
	})

	return r
}

// corsMiddleware adds permissive CORS headers for API responses.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// JSONContentType sets the Content-Type response header to application/json.
func JSONContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

