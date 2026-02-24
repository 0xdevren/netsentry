package drift

import "strings"

// Comparator computes line-level diffs between two raw configurations.
type Comparator struct{}

// NewComparator constructs a Comparator.
func NewComparator() *Comparator { return &Comparator{} }

// Compare produces a DiffResult between baselineData and currentData for deviceID.
func (c *Comparator) Compare(deviceID string, baselineData, currentData []byte) *DiffResult {
	baseLines := splitLines(baselineData)
	currLines := splitLines(currentData)

	// Build line sets with pre-allocated capacity.
	baseSet := make(map[string]struct{}, len(baseLines))
	for _, line := range baseLines {
		baseSet[line] = struct{}{}
	}

	currSet := make(map[string]struct{}, len(currLines))
	for _, line := range currLines {
		currSet[line] = struct{}{}
	}

	// Pre-allocate diff slices.
	maxChanges := len(baseLines) + len(currLines)
	result := &DiffResult{
		DeviceID: deviceID,
		Added:    make([]LineDiff, 0, maxChanges/4),
		Removed:  make([]LineDiff, 0, maxChanges/4),
	}

	// Lines removed (in baseline but not in current).
	for i, line := range baseLines {
		if _, ok := currSet[line]; !ok {
			result.Removed = append(result.Removed, LineDiff{
				Type:       ChangeRemoved,
				Line:       line,
				LineNumber: i + 1,
			})
		}
	}

	// Lines added (in current but not in baseline).
	for i, line := range currLines {
		if _, ok := baseSet[line]; !ok {
			result.Added = append(result.Added, LineDiff{
				Type:       ChangeAdded,
				Line:       line,
				LineNumber: i + 1,
			})
		}
	}

	result.HasChanges = len(result.Added) > 0 || len(result.Removed) > 0
	return result
}

// splitLines splits raw bytes into trimmed, non-empty lines.
func splitLines(data []byte) []string {
	raw := strings.Split(string(data), "\n")
	// Pre-allocate with estimated capacity (assume 80% non-empty).
	out := make([]string, 0, len(raw)*4/5)
	for _, l := range raw {
		l = strings.TrimRight(l, "\r")
		if strings.TrimSpace(l) != "" {
			out = append(out, l)
		}
	}
	return out
}
