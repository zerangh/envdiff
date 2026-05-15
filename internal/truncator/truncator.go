// Package truncator provides utilities for truncating long values in env maps
// and diff results to a maximum length, making output more readable.
package truncator

import (
	"strings"

	"github.com/envdiff/envdiff/internal/differ"
)

const defaultMaxLen = 64
const ellipsis = "..."

// Options controls truncation behaviour.
type Options struct {
	// MaxLen is the maximum number of characters allowed before truncation.
	// Defaults to 64 if zero.
	MaxLen int
}

// DefaultOptions returns sensible truncation defaults.
func DefaultOptions() Options {
	return Options{MaxLen: defaultMaxLen}
}

// TruncateValue shortens a single string value if it exceeds MaxLen.
func TruncateValue(v string, opts Options) string {
	max := opts.MaxLen
	if max <= 0 {
		max = defaultMaxLen
	}
	if len(v) <= max {
		return v
	}
	return v[:max] + ellipsis
}

// TruncateEnv returns a new map with all values truncated.
func TruncateEnv(env map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = TruncateValue(v, opts)
	}
	return out
}

// TruncateResult returns a copy of the differ.Result with long values truncated
// in the MismatchedValues entries.
func TruncateResult(r differ.Result, opts Options) differ.Result {
	out := differ.Result{
		MissingInLeft:    append([]string(nil), r.MissingInLeft...),
		MissingInRight:   append([]string(nil), r.MissingInRight...),
		MismatchedValues: make([]differ.Mismatch, len(r.MismatchedValues)),
	}
	for i, m := range r.MismatchedValues {
		out.MismatchedValues[i] = differ.Mismatch{
			Key:        m.Key,
			LeftValue:  TruncateValue(m.LeftValue, opts),
			RightValue: TruncateValue(m.RightValue, opts),
		}
	}
	return out
}

// isTruncated reports whether TruncateValue would shorten the given string.
func isTruncated(v string, max int) bool {
	if max <= 0 {
		max = defaultMaxLen
	}
	return len(v) > max
}

// Summary returns a human-readable line describing how many values were truncated.
func Summary(r differ.Result, opts Options) string {
	max := opts.MaxLen
	if max <= 0 {
		max = defaultMaxLen
	}
	count := 0
	for _, m := range r.MismatchedValues {
		if isTruncated(m.LeftValue, max) || isTruncated(m.RightValue, max) {
			count++
		}
	}
	if count == 0 {
		return "no values truncated"
	}
	return strings.Join([]string{
		"truncated",
		strconv.Itoa(count),
		"value(s) exceeding",
		strconv.Itoa(max),
		"characters",
	}, " ")
}
