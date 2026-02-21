package api

import (
	"encoding/json"
	"net/http"

	"github.com/0xdevren/netsentry/internal/app"
)

// HealthResponse is the response body for the health endpoint.
type HealthResponse struct {
	Status    string `json:"status"`
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildDate string `json:"build_date"`
}

// HealthCheckHandler returns the full health status with build metadata.
func HealthCheckHandler(appCtx *app.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(HealthResponse{ //nolint:errcheck
			Status:    "ok",
			Version:   app.Version,
			Commit:    app.Commit,
			BuildDate: app.BuildDate,
		})
	}
}
