// Package canonicalizer normalizes .env key names to a canonical form,
// helping detect keys that differ only by casing or separator style.
package canonicalizer

import (
	"strings"

	"github.com/user/envdiff/internal/differ"
)

// Options controls how canonicalization is applied.
type Options struct {
	// NormalizeCase converts all keys to UPPER_CASE when true.
	NormalizeCase bool
	// NormalizeSeparators replaces hyphens and dots in keys with underscores.
	NormalizeSeparators bool
}

// DefaultOptions returns sensible defaults for canonicalization.
func DefaultOptions() Options {
	return Options{
		NormalizeCase:       true,
		NormalizeSeparators: true,
	}
}

// Canonicalize returns a normalized version of the given key according to opts.
func Canonicalize(key string, opts Options) string {
	if opts.NormalizeSeparators {
		key = strings.NewReplacer("-", "_", ".", "_").Replace(key)
	}
	if opts.NormalizeCase {
		key = strings.ToUpper(key)
	}
	return key
}

// NormalizeEnv returns a new map whose keys have been canonicalized.
// If two keys collapse to the same canonical form, the last one wins.
func NormalizeEnv(env map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[Canonicalize(k, opts)] = v
	}
	return out
}

// NormalizeResult returns a new differ.Result whose key names inside
// MissingInLeft, MissingInRight, and Mismatched have been canonicalized.
func NormalizeResult(result differ.Result, opts Options) differ.Result {
	normLeft := make([]string, len(result.MissingInLeft))
	for i, k := range result.MissingInLeft {
		normLeft[i] = Canonicalize(k, opts)
	}

	normRight := make([]string, len(result.MissingInRight))
	for i, k := range result.MissingInRight {
		normRight[i] = Canonicalize(k, opts)
	}

	normMismatched := make([]differ.Mismatch, len(result.Mismatched))
	for i, m := range result.Mismatched {
		normMismatched[i] = differ.Mismatch{
			Key:        Canonicalize(m.Key, opts),
			LeftValue:  m.LeftValue,
			RightValue: m.RightValue,
		}
	}

	return differ.Result{
		MissingInLeft:  normLeft,
		MissingInRight: normRight,
		Mismatched:     normMismatched,
	}
}
