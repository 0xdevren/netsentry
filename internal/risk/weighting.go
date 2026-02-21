package risk

import "github.com/0xdevren/netsentry/internal/policy"

// WeightTable maps severity levels to multiplier factors for risk scoring.
type WeightTable struct {
	factors map[policy.Severity]float64
}

// DefaultWeights returns the standard NetSentry severity weight table.
func DefaultWeights() *WeightTable {
	return &WeightTable{factors: map[policy.Severity]float64{
		policy.SeverityCritical: 4.0,
		policy.SeverityHigh:     2.5,
		policy.SeverityMedium:   1.5,
		policy.SeverityLow:      1.0,
		policy.SeverityInfo:     0.1,
	}}
}

// Factor returns the weight multiplier for a severity level.
func (w *WeightTable) Factor(s policy.Severity) float64 {
	if f, ok := w.factors[s]; ok {
		return f
	}
	return 1.0
}

// WithFactor returns a new WeightTable with the given severity factor overridden.
func (w *WeightTable) WithFactor(s policy.Severity, factor float64) *WeightTable {
	newFactors := make(map[policy.Severity]float64, len(w.factors))
	for k, v := range w.factors {
		newFactors[k] = v
	}
	newFactors[s] = factor
	return &WeightTable{factors: newFactors}
}
