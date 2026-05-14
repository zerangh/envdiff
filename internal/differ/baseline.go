package differ

import (
	"github.com/user/envdiff/internal/parser"
)

// BaselineResult holds the comparison between a current env and a stored baseline.
type BaselineResult struct {
	// Added contains keys present in current but not in baseline.
	Added []string
	// Removed contains keys present in baseline but not in current.
	Removed []string
	// Changed contains keys whose values differ between baseline and current.
	Changed []MismatchedKey
	// Unchanged is the count of keys identical in both.
	Unchanged int
}

// CompareToBaseline diffs a current env map against a baseline env map.
// It returns a BaselineResult describing what has changed since the baseline
// was recorded.
func CompareToBaseline(baseline, current map[string]string) BaselineResult {
	var added, removed []string
	var changed []MismatchedKey
	unchanged := 0

	for k, cv := range current {
		bv, ok := baseline[k]
		if !ok {
			added = append(added, k)
		} else if bv != cv {
			changed = append(changed, MismatchedKey{Key: k, LeftValue: bv, RightValue: cv})
		} else {
			unchanged++
		}
	}

	for k := range baseline {
		if _, ok := current[k]; !ok {
			removed = append(removed, k)
		}
	}

	sortStrings(added)
	sortStrings(removed)
	sortMismatched(changed)

	return BaselineResult{
		Added:     added,
		Removed:   removed,
		Changed:   changed,
		Unchanged: unchanged,
	}
}

// CompareFilesToBaseline is a convenience wrapper that parses both files
// before calling CompareToBaseline.
func CompareFilesToBaseline(baselineFile, currentFile string) (BaselineResult, error) {
	base, err := parser.ParseFile(baselineFile)
	if err != nil {
		return BaselineResult{}, err
	}
	curr, err := parser.ParseFile(currentFile)
	if err != nil {
		return BaselineResult{}, err
	}
	return CompareToBaseline(base, curr), nil
}
