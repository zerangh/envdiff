// Package sorter provides utilities for sorting and ordering .env key-value
// maps and diff results by various strategies such as alphabetical, by prefix
// group, or by change type (missing, mismatched, clean).
package sorter

import (
	"sort"
	"strings"

	"github.com/user/envdiff/internal/differ"
)

// Strategy defines how keys should be ordered.
type Strategy string

const (
	Alpha      Strategy = "alpha"       // alphabetical by key name
	ChangeType Strategy = "change-type" // missing first, then mismatched, then clean
	Prefix     Strategy = "prefix"      // grouped by key prefix (e.g. DB_, AWS_)
)

// SortedEnv returns a slice of keys from the given env map sorted by the
// requested strategy. For strategies other than Alpha, ties are broken
// alphabetically.
func SortedEnv(env map[string]string, s Strategy) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	if s == Prefix {
		sort.SliceStable(keys, func(i, j int) bool {
			return extractPrefix(keys[i]) < extractPrefix(keys[j])
		})
	}
	return keys
}

// SortResult returns a new differ.Result whose slices are ordered by the
// requested strategy. The original result is not modified.
func SortResult(r differ.Result, s Strategy) differ.Result {
	out := differ.Result{
		MissingInRight: append([]string(nil), r.MissingInRight...),
		MissingInLeft:  append([]string(nil), r.MissingInLeft...),
		Mismatched:     append([]differ.Mismatch(nil), r.Mismatched...),
	}

	switch s {
	case ChangeType:
		// already grouped by type; just sort each bucket alphabetically
		sort.Strings(out.MissingInRight)
		sort.Strings(out.MissingInLeft)
		sort.Slice(out.Mismatched, func(i, j int) bool {
			return out.Mismatched[i].Key < out.Mismatched[j].Key
		})
	case Prefix:
		sort.SliceStable(out.MissingInRight, func(i, j int) bool {
			pi, pj := extractPrefix(out.MissingInRight[i]), extractPrefix(out.MissingInRight[j])
			if pi != pj {
				return pi < pj
			}
			return out.MissingInRight[i] < out.MissingInRight[j]
		})
		sort.SliceStable(out.MissingInLeft, func(i, j int) bool {
			pi, pj := extractPrefix(out.MissingInLeft[i]), extractPrefix(out.MissingInLeft[j])
			if pi != pj {
				return pi < pj
			}
			return out.MissingInLeft[i] < out.MissingInLeft[j]
		})
		sort.SliceStable(out.Mismatched, func(i, j int) bool {
			pi, pj := extractPrefix(out.Mismatched[i].Key), extractPrefix(out.Mismatched[j].Key)
			if pi != pj {
				return pi < pj
			}
			return out.Mismatched[i].Key < out.Mismatched[j].Key
		})
	default: // Alpha
		sort.Strings(out.MissingInRight)
		sort.Strings(out.MissingInLeft)
		sort.Slice(out.Mismatched, func(i, j int) bool {
			return out.Mismatched[i].Key < out.Mismatched[j].Key
		})
	}
	return out
}

// extractPrefix returns the portion of a key up to and including the first
// underscore, or the full key if no underscore is present.
func extractPrefix(key string) string {
	if idx := strings.Index(key, "_"); idx >= 0 {
		return key[:idx+1]
	}
	return key
}
