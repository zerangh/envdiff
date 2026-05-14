package formatter

import (
	"fmt"
	"strings"

	"github.com/your-org/envdiff/internal/differ"
)

// Style controls how key=value pairs are rendered.
type Style int

const (
	// StyleExport prefixes every line with "export ".
	StyleExport Style = iota
	// StylePlain emits bare KEY=VALUE lines.
	StylePlain
	// StyleDocker emits lines suitable for docker --env-file (no export, no quotes).
	StyleDocker
)

// Options configures the formatting behaviour.
type Options struct {
	Style        Style
	SortKeys     bool
	OmitEmpty    bool
	QuoteValues  bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Style:       StylePlain,
		SortKeys:    true,
		OmitEmpty:   false,
		QuoteValues: false,
	}
}

// FormatEnv renders a parsed env map into a slice of formatted lines.
func FormatEnv(env map[string]string, opts Options) []string {
	keys := sortedKeys(env, opts.SortKeys)
	lines := make([]string, 0, len(keys))
	for _, k := range keys {
		v := env[k]
		if opts.OmitEmpty && v == "" {
			continue
		}
		lines = append(lines, formatLine(k, v, opts))
	}
	return lines
}

// FormatResult renders the left-hand env from a diff result.
func FormatResult(result differ.Result, opts Options) []string {
	env := make(map[string]string)
	for _, m := range result.Mismatched {
		env[m.Key] = m.LeftValue
	}
	for _, k := range result.MissingInRight {
		env[k] = ""
	}
	return FormatEnv(env, opts)
}

func formatLine(key, value string, opts Options) string {
	if opts.QuoteValues && strings.ContainsAny(value, " \t\n#") {
		value = fmt.Sprintf("%q", value)
	}
	pair := key + "=" + value
	if opts.Style == StyleExport {
		return "export " + pair
	}
	return pair
}

func sortedKeys(env map[string]string, doSort bool) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	if doSort {
		for i := 1; i < len(keys); i++ {
			for j := i; j > 0 && keys[j] < keys[j-1]; j-- {
				keys[j], keys[j-1] = keys[j-1], keys[j]
			}
		}
	}
	return keys
}
