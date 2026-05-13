package differ

import "sort"

// sortStrings sorts a slice of strings in place.
func sortStrings(s []string) {
	sort.Strings(s)
}

// sortMismatched sorts a slice of MismatchedKey by key name.
func sortMismatched(m []MismatchedKey) {
	sort.Slice(m, func(i, j int) bool {
		return m[i].Key < m[j].Key
	})
}
