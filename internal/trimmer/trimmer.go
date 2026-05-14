// Package trimmer removes keys from an env map or diff result
// that match a given set of key names or prefixes, useful for
// stripping deployment-specific or CI-only variables before comparison.
package trimmer

import (
	"strings"

	"github.com/user/envdiff/internal/differ"
)

// Options controls which keys are removed.
type Options struct {
	// Keys is an exact list of key names to remove.
	Keys []string
	// Prefixes removes any key whose name starts with one of these strings.
	Prefixes []string
}

// TrimEnv returns a copy of env with matching keys removed.
func TrimEnv(env map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if shouldTrim(k, opts) {
			continue
		}
		out[k] = v
	}
	return out
}

// TrimResult returns a copy of the differ.Result with matching keys
// removed from MissingInRight, MissingInLeft, and Mismatched slices.
func TrimResult(r differ.Result, opts Options) differ.Result {
	return differ.Result{
		MissingInRight: trimStrings(r.MissingInRight, opts),
		MissingInLeft:  trimStrings(r.MissingInLeft, opts),
		Mismatched:     trimMismatched(r.Mismatched, opts),
	}
}

func shouldTrim(key string, opts Options) bool {
	for _, k := range opts.Keys {
		if k == key {
			return true
		}
	}
	for _, p := range opts.Prefixes {
		if strings.HasPrefix(key, p) {
			return true
		}
	}
	return false
}

func trimStrings(keys []string, opts Options) []string {
	var out []string
	for _, k := range keys {
		if !shouldTrim(k, opts) {
			out = append(out, k)
		}
	}
	return out
}

func trimMismatched(pairs []differ.Mismatch, opts Options) []differ.Mismatch {
	var out []differ.Mismatch
	for _, m := range pairs {
		if !shouldTrim(m.Key, opts) {
			out = append(out, m)
		}
	}
	return out
}
