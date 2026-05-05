package differ

// Result holds the outcome of diffing two env file maps.
type Result struct {
	// MissingInRight contains keys present in left but absent in right.
	MissingInRight []string
	// MissingInLeft contains keys present in right but absent in left.
	MissingInLeft []string
	// Mismatched contains keys present in both files but with different values.
	Mismatched []MismatchedKey
}

// MismatchedKey describes a key whose value differs between two env files.
type MismatchedKey struct {
	Key        string
	LeftValue  string
	RightValue string
}

// Diff compares two parsed env maps (key -> value) and returns a Result
// describing any missing or mismatched keys.
func Diff(left, right map[string]string) Result {
	var result Result

	for k, lv := range left {
		rv, ok := right[k]
		if !ok {
			result.MissingInRight = append(result.MissingInRight, k)
			continue
		}
		if lv != rv {
			result.Mismatched = append(result.Mismatched, MismatchedKey{
				Key:        k,
				LeftValue:  lv,
				RightValue: rv,
			})
		}
	}

	for k := range right {
		if _, ok := left[k]; !ok {
			result.MissingInLeft = append(result.MissingInLeft, k)
		}
	}

	sortStrings(result.MissingInRight)
	sortStrings(result.MissingInLeft)
	sortMismatched(result.Mismatched)

	return result
}

// HasDiff returns true when the Result contains any differences.
func (r Result) HasDiff() bool {
	return len(r.MissingInRight) > 0 || len(r.MissingInLeft) > 0 || len(r.Mismatched) > 0
}
