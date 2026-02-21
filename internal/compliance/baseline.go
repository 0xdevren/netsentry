// Package compliance provides CIS benchmark and baseline compliance tracking.
package compliance

import (
	"fmt"
	"time"

	"github.com/0xdevren/netsentry/internal/policy"
)

// BaselineEntry is a snapshot of a device's compliance state at a point in time.
type BaselineEntry struct {
	// DeviceID is the evaluated device.
	DeviceID string `json:"device_id"`
	// PolicyName is the policy used.
	PolicyName string `json:"policy_name"`
	// Score is the compliance percentage at baseline time.
	Score float64 `json:"score"`
	// RecordedAt is when the baseline was captured.
	RecordedAt time.Time `json:"recorded_at"`
	// Summary is the rule evaluation summary.
	Summary policy.ReportSummary `json:"summary"`
}

// BaselineStore maintains an in-memory collection of baseline entries.
type BaselineStore struct {
	entries map[string]BaselineEntry
}

// NewBaselineStore constructs an empty BaselineStore.
func NewBaselineStore() *BaselineStore {
	return &BaselineStore{entries: make(map[string]BaselineEntry)}
}

// Record stores a new baseline entry for the given device.
func (b *BaselineStore) Record(report *policy.Report) {
	b.entries[report.Device.ID] = BaselineEntry{
		DeviceID:   report.Device.ID,
		PolicyName: report.Policy,
		Score:      report.Summary.Score,
		RecordedAt: time.Now().UTC(),
		Summary:    report.Summary,
	}
}

// Get returns the baseline entry for the given device ID.
func (b *BaselineStore) Get(deviceID string) (BaselineEntry, error) {
	e, ok := b.entries[deviceID]
	if !ok {
		return BaselineEntry{}, fmt.Errorf("baseline: no entry for device %q", deviceID)
	}
	return e, nil
}

// Compare returns the change in compliance score since the baseline.
func (b *BaselineStore) Compare(current *policy.Report) (delta float64, err error) {
	baseline, err := b.Get(current.Device.ID)
	if err != nil {
		return 0, err
	}
	return current.Summary.Score - baseline.Score, nil
}
