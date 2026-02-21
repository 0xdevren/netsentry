package risk

import "github.com/0xdevren/netsentry/internal/policy"

// Scorer computes risk assessments from validation reports.
type Scorer struct {
	weights *WeightTable
}

// NewScorer constructs a Scorer with the given weight table.
// Passing nil uses the DefaultWeights.
func NewScorer(weights *WeightTable) *Scorer {
	if weights == nil {
		weights = DefaultWeights()
	}
	return &Scorer{weights: weights}
}

// Score computes a RiskModel for the given validation report.
func (s *Scorer) Score(report *policy.Report) *RiskModel {
	return FromReport(report, s.weights)
}

// ScoreMany computes risk models for multiple validation reports.
func (s *Scorer) ScoreMany(reports []*policy.Report) []*RiskModel {
	out := make([]*RiskModel, 0, len(reports))
	for _, r := range reports {
		out = append(out, s.Score(r))
	}
	return out
}
