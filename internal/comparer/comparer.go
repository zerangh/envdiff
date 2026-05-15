// Package comparer provides value-level comparison utilities for env entries,
// including case-insensitive, trimmed, and normalized comparison modes.
package comparer

import (
	"strings"

	"github.com/user/envdiff/internal/differ"
)

// Options controls how values are compared.
type Options struct {
	// CaseInsensitive treats values as equal if they differ only in case.
	CaseInsensitive bool
	// TrimSpace ignores leading/trailing whitespace when comparing.
	TrimSpace bool
	// NormalizeBools treats "true"/"1"/"yes" and "false"/"0"/"no" as equivalent.
	NormalizeBools bool
}

// DefaultOptions returns the default comparison options.
func DefaultOptions() Options {
	return Options{
		CaseInsensitive: false,
		TrimSpace:       true,
		NormalizeBools:  false,
	}
}

// Equal reports whether two values are considered equal under the given options.
func Equal(a, b string, opts Options) bool {
	if opts.TrimSpace {
		a = strings.TrimSpace(a)
		b = strings.TrimSpace(b)
	}
	if opts.NormalizeBools {
		a = normalizeBool(a)
		b = normalizeBool(b)
	}
	if opts.CaseInsensitive {
		return strings.EqualFold(a, b)
	}
	return a == b
}

// FilterMismatches returns a new differ.Result where mismatched entries that
// are considered equal under opts are removed from the Mismatched list.
func FilterMismatches(result differ.Result, opts Options) differ.Result {
	filtered := make([]differ.Mismatch, 0, len(result.Mismatched))
	for _, m := range result.Mismatched {
		if !Equal(m.LeftValue, m.RightValue, opts) {
			filtered = append(filtered, m)
		}
	}
	return differ.Result{
		MissingInLeft:  result.MissingInLeft,
		MissingInRight: result.MissingInRight,
		Mismatched:     filtered,
	}
}

func normalizeBool(v string) string {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "true", "1", "yes":
		return "true"
	case "false", "0", "no":
		return "false"
	}
	return v
}
