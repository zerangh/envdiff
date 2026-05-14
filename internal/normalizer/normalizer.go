// Package normalizer provides utilities for normalizing .env file values
// before comparison, such as trimming whitespace, lowercasing booleans,
// and collapsing equivalent representations.
package normalizer

import (
	"strings"

	"github.com/your-org/envdiff/internal/differ"
)

// Options controls which normalizations are applied.
type Options struct {
	// TrimSpace removes leading and trailing whitespace from values.
	TrimSpace bool
	// NormalizeBools converts true/false/yes/no/1/0 to a canonical form.
	NormalizeBools bool
	// LowercaseValues converts all values to lowercase.
	LowercaseValues bool
}

// DefaultOptions returns the recommended normalization settings.
func DefaultOptions() Options {
	return Options{
		TrimSpace:      true,
		NormalizeBools: true,
		LowercaseValues: false,
	}
}

// NormalizeValue applies the selected normalizations to a single value.
func NormalizeValue(v string, opts Options) string {
	if opts.TrimSpace {
		v = strings.TrimSpace(v)
	}
	if opts.LowercaseValues {
		v = strings.ToLower(v)
	}
	if opts.NormalizeBools {
		v = normalizeBool(v)
	}
	return v
}

// NormalizeEnv applies normalization to every value in an env map.
func NormalizeEnv(env map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = NormalizeValue(v, opts)
	}
	return out
}

// NormalizeResult applies normalization to the values inside a differ.Result,
// re-evaluating mismatches after normalization so that equivalent values are
// no longer reported as different.
func NormalizeResult(r differ.Result, opts Options) differ.Result {
	normalized := differ.Result{
		MissingInLeft:  r.MissingInLeft,
		MissingInRight: r.MissingInRight,
	}
	for _, m := range r.Mismatched {
		lv := NormalizeValue(m.LeftValue, opts)
		rv := NormalizeValue(m.RightValue, opts)
		if lv == rv {
			// Values are equivalent after normalization — skip.
			continue
		}
		normalized.Mismatched = append(normalized.Mismatched, differ.Mismatch{
			Key:        m.Key,
			LeftValue:  lv,
			RightValue: rv,
		})
	}
	return normalized
}

// normalizeBool converts common boolean-like strings to a canonical "true"/"false".
func normalizeBool(v string) string {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "yes", "1", "on", "true":
		return "true"
	case "no", "0", "off", "false":
		return "false"
	}
	return v
}
