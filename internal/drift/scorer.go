package drift

// DriftScore summarises the magnitude of configuration drift.
type DriftScore struct {
	// DeviceID is the assessed device.
	DeviceID string `json:"device_id"`
	// LinesAdded is the number of new lines in the current configuration.
	LinesAdded int `json:"lines_added"`
	// LinesRemoved is the number of lines removed relative to the baseline.
	LinesRemoved int `json:"lines_removed"`
	// TotalChanges is the total line change count.
	TotalChanges int `json:"total_changes"`
	// DriftPercent is the percentage of the baseline changed (0-100).
	DriftPercent float64 `json:"drift_percent"`
	// Significant indicates whether the drift exceeds the configured threshold.
	Significant bool `json:"significant"`
}

// DriftScorer computes a DriftScore from a DiffResult.
type DriftScorer struct {
	// Threshold is the drift percentage above which drift is considered significant.
	Threshold float64
}

// NewDriftScorer constructs a DriftScorer with the given significance threshold.
func NewDriftScorer(threshold float64) *DriftScorer {
	if threshold <= 0 {
		threshold = 5.0
	}
	return &DriftScorer{Threshold: threshold}
}

// Score computes a DriftScore from a DiffResult and the baseline line count.
func (s *DriftScorer) Score(diff *DiffResult, baselineLineCount int) *DriftScore {
	ds := &DriftScore{
		DeviceID:     diff.DeviceID,
		LinesAdded:   len(diff.Added),
		LinesRemoved: len(diff.Removed),
	}
	ds.TotalChanges = ds.LinesAdded + ds.LinesRemoved

	if baselineLineCount > 0 {
		ds.DriftPercent = float64(ds.TotalChanges) / float64(baselineLineCount) * 100
	}
	ds.Significant = ds.DriftPercent >= s.Threshold
	return ds
}
