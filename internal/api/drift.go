package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/0xdevren/netsentry/internal/app"
	"github.com/0xdevren/netsentry/internal/drift"
)

// driftStore is a simple in-memory store of configuration snapshots.
// In production this would be backed by a persistent store.
var driftStore = struct {
	baseline map[string][]byte
	current  map[string][]byte
}{
	baseline: make(map[string][]byte),
	current:  make(map[string][]byte),
}

// DriftHandler handles GET /api/v1/drift/{deviceID}.
func DriftHandler(appCtx *app.Context) http.HandlerFunc {
	comparator := drift.NewComparator()
	scorer := drift.NewDriftScorer(5.0)

	return func(w http.ResponseWriter, r *http.Request) {
		deviceID := chi.URLParam(r, "deviceID")
		if deviceID == "" {
			jsonError(w, "deviceID is required", http.StatusBadRequest)
			return
		}

		baseline, ok := driftStore.baseline[deviceID]
		if !ok {
			jsonError(w, "no baseline found for device "+deviceID, http.StatusNotFound)
			return
		}
		current, ok := driftStore.current[deviceID]
		if !ok {
			jsonError(w, "no current config found for device "+deviceID, http.StatusNotFound)
			return
		}

		diff := comparator.Compare(deviceID, baseline, current)
		score := scorer.Score(diff, len(splitLinesToCount(baseline)))

		json.NewEncoder(w).Encode(map[string]interface{}{ //nolint:errcheck
			"diff":  diff,
			"score": score,
		})
	}
}

func splitLinesToCount(data []byte) []string {
	var lines []string
	s := string(data)
	current := ""
	for _, c := range s {
		if c == '\n' {
			if current != "" {
				lines = append(lines, current)
			}
			current = ""
		} else {
			current += string(c)
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}
