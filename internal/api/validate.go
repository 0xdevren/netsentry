package api

import (
	"encoding/json"
	"net/http"

	"github.com/0xdevren/netsentry/internal/app"
	"github.com/0xdevren/netsentry/internal/config"
	"github.com/0xdevren/netsentry/internal/model"
	"github.com/0xdevren/netsentry/internal/parser"
	"github.com/0xdevren/netsentry/internal/policy"
	"github.com/0xdevren/netsentry/internal/validator"
)

// ValidateRequest is the JSON body for POST /api/v1/validate.
type ValidateRequest struct {
	// Config is the raw device configuration text.
	Config string `json:"config"`
	// PolicyYAML is the raw policy YAML text.
	PolicyYAML string `json:"policy_yaml"`
	// Strict treats warnings as failures.
	Strict bool `json:"strict"`
	// Concurrency overrides the default worker count.
	Concurrency int `json:"concurrency"`
}

// ValidateHandler handles POST /api/v1/validate.
func ValidateHandler(appCtx *app.Context) http.HandlerFunc {
	detector := config.NewDetector()
	policyLoader := policy.NewLoader()

	return func(w http.ResponseWriter, r *http.Request) {
		var req ValidateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			jsonError(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}
		if req.Config == "" {
			jsonError(w, "config is required", http.StatusBadRequest)
			return
		}
		if req.PolicyYAML == "" {
			jsonError(w, "policy_yaml is required", http.StatusBadRequest)
			return
		}

		rawData := []byte(req.Config)
		deviceType := detector.Detect(rawData)

		device := model.Device{
			ID:   "api-request",
			Type: deviceType,
		}
		parsedCfg, err := parser.Parse(r.Context(), deviceType, rawData, device)
		if err != nil {
			jsonError(w, "parse error: "+err.Error(), http.StatusUnprocessableEntity)
			return
		}

		pol, err := policyLoader.LoadBytes([]byte(req.PolicyYAML))
		if err != nil {
			jsonError(w, "policy error: "+err.Error(), http.StatusBadRequest)
			return
		}

		concurrency := req.Concurrency
		if concurrency <= 0 {
			concurrency = 4
		}

		report, err := validator.Validate(r.Context(), validator.ValidationRequest{
			Config:      parsedCfg,
			Policy:      pol,
			Strict:      req.Strict,
			Concurrency: concurrency,
		})
		if err != nil {
			appCtx.Logger.Error("validate handler: validation error", err)
			jsonError(w, "validation error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(report) //nolint:errcheck
	}
}

// HealthHandler handles GET /healthz.
func HealthHandler(appCtx *app.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{ //nolint:errcheck
			"status":  "ok",
			"version": app.Version,
			"commit":  app.Commit,
		})
	}
}

// ReadyHandler handles GET /readyz.
func ReadyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ready"}`)) //nolint:errcheck
	}
}

// jsonError writes a JSON-formatted error response.
func jsonError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message}) //nolint:errcheck
}


