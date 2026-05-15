// Package stripper removes keys from env maps or diff results based on
// exact matches, prefix patterns, or suffix patterns.
package stripper

import (
	"strings"

	"github.com/user/envdiff/internal/differ"
)

// Options controls which keys are stripped.
type Options struct {
	// Keys is a list of exact key names to remove.
	Keys []string
	// Prefixes removes any key that starts with one of these strings.
	Prefixes []string
	// Suffixes removes any key that ends with one of these strings.
	Suffixes []string
}

// StripEnv removes matching keys from an env map and returns a new map.
func StripEnv(env map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if !shouldStrip(k, opts) {
			out[k] = v
		}
	}
	return out
}

// StripResult removes matching keys from a differ.Result and returns a new one.
func StripResult(r differ.Result, opts Options) differ.Result {
	out := differ.Result{}

	for _, k := range r.MissingInRight {
		if !shouldStrip(k, opts) {
			out.MissingInRight = append(out.MissingInRight, k)
		}
	}

	for _, k := range r.MissingInLeft {
		if !shouldStrip(k, opts) {
			out.MissingInLeft = append(out.MissingInLeft, k)
		}
	}

	for _, m := range r.Mismatched {
		if !shouldStrip(m.Key, opts) {
			out.Mismatched = append(out.Mismatched, m)
		}
	}

	return out
}

func shouldStrip(key string, opts Options) bool {
	for _, k := range opts.Keys {
		if key == k {
			return true
		}
	}
	for _, p := range opts.Prefixes {
		if strings.HasPrefix(key, p) {
			return true
		}
	}
	for _, s := range opts.Suffixes {
		if strings.HasSuffix(key, s) {
			return true
		}
	}
	return false
}
