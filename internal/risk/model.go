// Package risk provides risk scoring for policy violations.
package risk

import "github.com/0xdevren/netsentry/internal/policy"

// RiskLevel describes the aggregate risk tier of a device.
type RiskLevel string

const (
	RiskLevelCritical RiskLevel = "CRITICAL"
	RiskLevelHigh     RiskLevel = "HIGH"
	RiskLevelMedium   RiskLevel = "MEDIUM"
	RiskLevelLow      RiskLevel = "LOW"
	RiskLevelNone     RiskLevel = "NONE"
)

// RiskModel is the result of a risk assessment for a single device.
type RiskModel struct {
	// DeviceID is the identifier of the assessed device.
	DeviceID string `json:"device_id"`
	// RawScore is the unweighted aggregate violation score.
	RawScore float64 `json:"raw_score"`
	// WeightedScore is the weighted aggregate violation score (0-100).
	WeightedScore float64 `json:"weighted_score"`
	// Level is the derived risk tier.
	Level RiskLevel `json:"level"`
	// ViolationBreakdown maps severity string to violation count.
	ViolationBreakdown map[string]int `json:"violation_breakdown"`
}

// FromReport derives a RiskModel from a validation report.
func FromReport(rep *policy.Report, weights *WeightTable) *RiskModel {
	if weights == nil {
		weights = DefaultWeights()
	}
	rm := &RiskModel{
		DeviceID:           rep.Device.ID,
		ViolationBreakdown: make(map[string]int),
	}

	for _, r := range rep.Results {
		if r.Status != policy.StatusFail {
			continue
		}
		rm.ViolationBreakdown[string(r.Severity)]++
		rm.RawScore += float64(r.Severity.Weight())
		rm.WeightedScore += float64(r.Severity.Weight()) * weights.Factor(r.Severity)
	}

	rm.Level = classifyRisk(rm.WeightedScore)
	return rm
}

func classifyRisk(score float64) RiskLevel {
	switch {
	case score >= 500:
		return RiskLevelCritical
	case score >= 200:
		return RiskLevelHigh
	case score >= 75:
		return RiskLevelMedium
	case score > 0:
		return RiskLevelLow
	default:
		return RiskLevelNone
	}
}
