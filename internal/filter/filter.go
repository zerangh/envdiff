// Package filter provides utilities for filtering diff results
// based on key patterns, prefixes, or exclusion rules.
package filter

import (
	"strings"

	"github.com/user/envdiff/internal/differ"
)

// Options holds filtering configuration.
type Options struct {
	// Prefix restricts results to keys that start with this string.
	Prefix string
	// Exclude is a list of exact key names to omit from results.
	Exclude []string
	// OnlyMissing limits results to keys missing in either environment.
	OnlyMissing bool
	// OnlyMismatched limits results to keys present in both but with differing values.
	OnlyMismatched bool
}

// Apply returns a new DiffResult containing only the entries that satisfy
// the given Options. If no filtering options are set the original result
// is returned unchanged.
func Apply(result differ.DiffResult, opts Options) differ.DiffResult {
	excludeSet := make(map[string]bool, len(opts.Exclude))
	for _, k := range opts.Exclude {
		excludeSet[k] = true
	}

	out := differ.DiffResult{}

	if !opts.OnlyMismatched {
		for _, k := range result.MissingInRight {
			if matchKey(k, opts.Prefix, excludeSet) {
				out.MissingInRight = append(out.MissingInRight, k)
			}
		}
		for _, k := range result.MissingInLeft {
			if matchKey(k, opts.Prefix, excludeSet) {
				out.MissingInLeft = append(out.MissingInLeft, k)
			}
		}
	}

	if !opts.OnlyMissing {
		for _, m := range result.Mismatched {
			if matchKey(m.Key, opts.Prefix, excludeSet) {
				out.Mismatched = append(out.Mismatched, m)
			}
		}
	}

	return out
}

// matchKey returns true when key satisfies the prefix constraint and is not
// in the exclusion set.
func matchKey(key, prefix string, exclude map[string]bool) bool {
	if exclude[key] {
		return false
	}
	if prefix != "" && !strings.HasPrefix(key, prefix) {
		return false
	}
	return true
}
