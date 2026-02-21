package drift

import "strings"

// Comparator computes line-level diffs between two raw configurations.
type Comparator struct{}

// NewComparator constructs a Comparator.
func NewComparator() *Comparator { return &Comparator{} }

// Compare produces a DiffResult between baselineData and currentData for deviceID.
func (c *Comparator) Compare(deviceID string, baselineData, currentData []byte) *DiffResult {
	baseLines := toLineSet(baselineData)
	currLines := toLineSet(currentData)

	result := &DiffResult{DeviceID: deviceID}

	// Lines removed (in baseline but not in current).
	for i, line := range splitLines(baselineData) {
		if _, ok := currLines[line]; !ok {
			result.Removed = append(result.Removed, LineDiff{
				Type:       ChangeRemoved,
				Line:       line,
				LineNumber: i + 1,
			})
		}
	}

	// Lines added (in current but not in baseline).
	for i, line := range splitLines(currentData) {
		if _, ok := baseLines[line]; !ok {
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

// toLineSet converts raw config bytes into a set of non-empty trimmed lines.
func toLineSet(data []byte) map[string]struct{} {
	set := make(map[string]struct{})
	for _, line := range splitLines(data) {
		set[line] = struct{}{}
	}
	return set
}

// splitLines splits raw bytes into trimmed, non-empty lines.
func splitLines(data []byte) []string {
	raw := strings.Split(string(data), "\n")
	out := make([]string, 0, len(raw))
	for _, l := range raw {
		l = strings.TrimRight(l, "\r")
		if strings.TrimSpace(l) != "" {
			out = append(out, l)
		}
	}
	return out
}
