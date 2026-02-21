package api

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/0xdevren/netsentry/internal/app"
	"github.com/0xdevren/netsentry/internal/policy"
	"github.com/0xdevren/netsentry/internal/policy/dsl"
)

// PolicyListHandler handles GET /api/v1/policy.
// Returns a list of policy file paths found in the configured policies directory.
func PolicyListHandler(appCtx *app.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dir := r.URL.Query().Get("dir")
		if dir == "" {
			dir = "policies"
		}
		entries, err := os.ReadDir(dir)
		if err != nil {
			jsonError(w, "cannot read policy directory: "+err.Error(), http.StatusInternalServerError)
			return
		}
		var files []string
		for _, e := range entries {
			if !e.IsDir() {
				files = append(files, e.Name())
			}
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"policies": files}) //nolint:errcheck
	}
}

// LintRequest is the JSON body for POST /api/v1/policy/lint.
type LintRequest struct {
	// PolicyYAML is the raw policy YAML to lint.
	PolicyYAML string `json:"policy_yaml"`
}

// PolicyLintHandler handles POST /api/v1/policy/lint.
func PolicyLintHandler(appCtx *app.Context) http.HandlerFunc {
	parser := dsl.NewParser()
	validator := dsl.NewValidator()
	loader := policy.NewLoader()

	return func(w http.ResponseWriter, r *http.Request) {
		var req LintRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			jsonError(w, "invalid request body", http.StatusBadRequest)
			return
		}

		raw, err := parser.ParseBytes([]byte(req.PolicyYAML))
		if err != nil {
			jsonError(w, "yaml parse error: "+err.Error(), http.StatusBadRequest)
			return
		}

		errs := validator.Validate(raw)
		if len(errs) > 0 {
			type lintError struct {
				Field   string `json:"field"`
				Message string `json:"message"`
				RuleID  string `json:"rule_id,omitempty"`
			}
			var lintErrs []lintError
			for _, e := range errs {
				lintErrs = append(lintErrs, lintError{Field: e.Field, Message: e.Message, RuleID: e.RuleID})
			}
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(map[string]interface{}{"valid": false, "errors": lintErrs}) //nolint:errcheck
			return
		}

		_, err = loader.LoadBytes([]byte(req.PolicyYAML))
		if err != nil {
			jsonError(w, "policy validation error: "+err.Error(), http.StatusUnprocessableEntity)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{"valid": true, "errors": nil}) //nolint:errcheck
	}
}
