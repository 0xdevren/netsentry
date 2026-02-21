package topology

import (
	"fmt"

	"github.com/0xdevren/netsentry/internal/topology/checks"
)

// AnalysisResult is the aggregate output of all topology checks.
type AnalysisResult struct {
	// Issues is the list of detected topology problems.
	Issues []checks.Issue
	// HasIssues indicates whether any issues were found.
	HasIssues bool
}

// String formats the analysis result for human-readable output.
func (a *AnalysisResult) String() string {
	if !a.HasIssues {
		return "Topology analysis: no issues detected.\n"
	}
	out := fmt.Sprintf("Topology analysis: %d issue(s) detected.\n", len(a.Issues))
	for _, issue := range a.Issues {
		out += fmt.Sprintf("  [%s] %s: %s\n", issue.Severity, issue.Code, issue.Message)
	}
	return out
}

// Analyzer runs all registered topology checks against a Graph.
type Analyzer struct {
	checkers []checks.Check
}

// NewAnalyzer constructs an Analyzer pre-loaded with all built-in checks.
func NewAnalyzer() *Analyzer {
	return &Analyzer{
		checkers: []checks.Check{
			&checks.DuplicateIPCheck{},
			&checks.SubnetOverlapCheck{},
			&checks.LoopCheck{},
			&checks.AdjacencyCheck{},
		},
	}
}

// Analyze runs all checks against the given Graph and returns the aggregate result.
func (a *Analyzer) Analyze(g *Graph) *AnalysisResult {
	result := &AnalysisResult{}
	tg := g.ToModel()
	for _, chk := range a.checkers {
		issues := chk.Run(tg)
		result.Issues = append(result.Issues, issues...)
	}
	result.HasIssues = len(result.Issues) > 0
	return result
}
