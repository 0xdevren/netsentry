package checks

import (
	"fmt"
	"github.com/0xdevren/netsentry/internal/model"
)

// AdjacencyCheck detects asymmetric links: links that are unidirectional
// where bidirectional connectivity is expected.
type AdjacencyCheck struct{}

// Run detects asymmetric topology links.
func (a *AdjacencyCheck) Run(g *model.TopologyGraph) []Issue {
	// Build reverse link set with pre-allocated capacity.
	nLinks := len(g.Links)
	type linkKey struct{ src, dst string }
	existing := make(map[linkKey]struct{}, nLinks)
	for _, link := range g.Links {
		existing[linkKey{link.SourceDevice, link.TargetDevice}] = struct{}{}
	}

	issues := make([]Issue, 0, nLinks/4) // Pre-allocate for typical case (25% asymmetric)
	for _, link := range g.Links {
		reverse := linkKey{link.TargetDevice, link.SourceDevice}
		if _, ok := existing[reverse]; !ok {
			issues = append(issues, Issue{
				Code:     "ADJ-ASYMMETRIC-001",
				Severity: "MEDIUM",
				Message:  fmt.Sprintf("asymmetric link: %s -> %s has no return path", link.SourceDevice, link.TargetDevice),
				DeviceID: link.SourceDevice,
			})
		}
	}
	return issues
}
