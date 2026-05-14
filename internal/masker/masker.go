// Package masker provides utilities for masking sensitive values
// in .env maps before display or export, supporting configurable
// patterns and partial reveal modes.
package masker

import (
	"strings"

	"github.com/your-org/envdiff/internal/differ"
)

const defaultMaskChar = "*"
const defaultRevealChars = 0

// Options controls masking behaviour.
type Options struct {
	// Patterns is a list of substrings; keys containing any of them are masked.
	Patterns []string
	// MaskChar is the character used to replace masked values. Defaults to "*".
	MaskChar string
	// RevealSuffix is the number of trailing characters to leave visible (0 = hide all).
	RevealSuffix int
}

// DefaultOptions returns sensible masking defaults.
func DefaultOptions() Options {
	return Options{
		Patterns:     []string{"SECRET", "PASSWORD", "TOKEN", "KEY", "PRIVATE", "CREDENTIAL"},
		MaskChar:     defaultMaskChar,
		RevealSuffix: defaultRevealChars,
	}
}

// IsSensitive reports whether key matches any of the given patterns (case-insensitive).
func IsSensitive(key string, patterns []string) bool {
	upper := strings.ToUpper(key)
	for _, p := range patterns {
		if strings.Contains(upper, strings.ToUpper(p)) {
			return true
		}
	}
	return false
}

// maskValue replaces a value with a mask string, optionally revealing a suffix.
func maskValue(value, maskChar string, revealSuffix int) string {
	if len(value) == 0 {
		return value
	}
	if revealSuffix <= 0 || revealSuffix >= len(value) {
		return strings.Repeat(maskChar, 8)
	}
	hidden := len(value) - revealSuffix
	return strings.Repeat(maskChar, hidden) + value[hidden:]
}

// MaskEnv returns a copy of env with sensitive values replaced.
func MaskEnv(env map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if IsSensitive(k, opts.Patterns) {
			out[k] = maskValue(v, opts.MaskChar, opts.RevealSuffix)
		} else {
			out[k] = v
		}
	}
	return out
}

// MaskResult returns a copy of result with sensitive mismatched values masked.
func MaskResult(result differ.Result, opts Options) differ.Result {
	masked := make([]differ.Mismatch, len(result.Mismatched))
	for i, m := range result.Mismatched {
		if IsSensitive(m.Key, opts.Patterns) {
			masked[i] = differ.Mismatch{
				Key:        m.Key,
				LeftValue:  maskValue(m.LeftValue, opts.MaskChar, opts.RevealSuffix),
				RightValue: maskValue(m.RightValue, opts.MaskChar, opts.RevealSuffix),
			}
		} else {
			masked[i] = m
		}
	}
	return differ.Result{
		MissingInLeft:  result.MissingInLeft,
		MissingInRight: result.MissingInRight,
		Mismatched:     masked,
	}
}
