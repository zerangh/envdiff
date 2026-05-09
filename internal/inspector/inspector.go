// Package inspector provides functionality to inspect a single .env file
// and report statistics such as key count, empty values, duplicate keys,
// and keys with special characters.
package inspector

import (
	"fmt"
	"strings"
	"unicode"
)

// Report holds the inspection results for a single .env file.
type Report struct {
	TotalKeys      int
	EmptyValues    []string
	DuplicateKeys  []string
	SpecialCharKeys []string
}

// Format returns a human-readable summary of the inspection report.
func (r Report) Format() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Total keys: %d\n", r.TotalKeys)

	if len(r.EmptyValues) == 0 {
		sb.WriteString("Empty values: none\n")
	} else {
		fmt.Fprintf(&sb, "Empty values (%d): %s\n", len(r.EmptyValues), strings.Join(r.EmptyValues, ", "))
	}

	if len(r.DuplicateKeys) == 0 {
		sb.WriteString("Duplicate keys: none\n")
	} else {
		fmt.Fprintf(&sb, "Duplicate keys (%d): %s\n", len(r.DuplicateKeys), strings.Join(r.DuplicateKeys, ", "))
	}

	if len(r.SpecialCharKeys) == 0 {
		sb.WriteString("Keys with special characters: none\n")
	} else {
		fmt.Fprintf(&sb, "Keys with special characters (%d): %s\n", len(r.SpecialCharKeys), strings.Join(r.SpecialCharKeys, ", "))
	}

	return sb.String()
}

// Inspect analyses the provided key-value map (as parsed from a .env file)
// along with the raw ordered keys slice to detect duplicates, empty values,
// and keys containing non-alphanumeric/underscore characters.
func Inspect(env map[string]string, orderedKeys []string) Report {
	seen := make(map[string]int)
	for _, k := range orderedKeys {
		seen[k]++
	}

	var duplicates []string
	for k, count := range seen {
		if count > 1 {
			duplicates = append(duplicates, k)
		}
	}

	var emptyVals []string
	var specialKeys []string
	for k, v := range env {
		if v == "" {
			emptyVals = append(emptyVals, k)
		}
		if hasSpecialChars(k) {
			specialKeys = append(specialKeys, k)
		}
	}

	return Report{
		TotalKeys:       len(env),
		EmptyValues:     emptyVals,
		DuplicateKeys:   duplicates,
		SpecialCharKeys: specialKeys,
	}
}

// hasSpecialChars returns true if the key contains characters other than
// uppercase/lowercase letters, digits, and underscores.
func hasSpecialChars(key string) bool {
	for _, r := range key {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return true
		}
	}
	return false
}
