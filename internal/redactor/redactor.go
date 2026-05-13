package redactor

import (
	"regexp"
	"strings"

	"github.com/user/envdiff/internal/differ"
)

// sensitivePatterns lists key substrings that indicate sensitive values.
var sensitivePatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(password|passwd|pwd)`),
	regexp.MustCompile(`(?i)(secret|token|apikey|api_key)`),
	regexp.MustCompile(`(?i)(private|credentials|auth)`),
}

// Options configures redaction behaviour.
type Options struct {
	// Mask is the string used to replace sensitive values. Defaults to "***".
	Mask string
	// ExtraPatterns are additional regexp patterns to treat as sensitive.
	ExtraPatterns []*regexp.Regexp
}

func defaultMask(opts Options) string {
	if opts.Mask != "" {
		return opts.Mask
	}
	return "***"
}

// IsSensitive reports whether a key name looks like it holds a sensitive value.
func IsSensitive(key string, extra []*regexp.Regexp) bool {
	for _, re := range sensitivePatterns {
		if re.MatchString(key) {
			return true
		}
	}
	for _, re := range extra {
		if re != nil && re.MatchString(key) {
			return true
		}
	}
	return false
}

// RedactEnv returns a copy of the env map with sensitive values replaced by mask.
func RedactEnv(env map[string]string, opts Options) map[string]string {
	mask := defaultMask(opts)
	out := make(map[string]string, len(env))
	for k, v := range env {
		if IsSensitive(k, opts.ExtraPatterns) {
			out[k] = mask
		} else {
			out[k] = v
		}
	}
	return out
}

// RedactResult returns a copy of the differ.Result with sensitive mismatched
// values replaced by mask so they are safe to display or log.
func RedactResult(r differ.Result, opts Options) differ.Result {
	mask := defaultMask(opts)
	redacted := make([]differ.Mismatch, len(r.Mismatched))
	for i, m := range r.Mismatched {
		if IsSensitive(m.Key, opts.ExtraPatterns) {
			m.LeftValue = mask
			m.RightValue = strings.Repeat("*", len(m.RightValue))
			if m.RightValue == "" {
				m.RightValue = mask
			}
		}
		redacted[i] = m
	}
	return differ.Result{
		MissingInLeft:  r.MissingInLeft,
		MissingInRight: r.MissingInRight,
		Mismatched:     redacted,
	}
}
