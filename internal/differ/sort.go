package differ

import "sort"

// sortStrings sorts a string slice in place.
func sortStrings(s []string) {
	sort.Strings(s)
}

// sortMismatched sorts a MismatchedKey slice by Key in place.
func sortMismatched(m []MismatchedKey) {
	sort.Slice(m, func(i, j int) bool {
		return m[i].Key < m[j].Key
	})
}
