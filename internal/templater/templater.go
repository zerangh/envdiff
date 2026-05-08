// Package templater generates a .env.example file from one or more parsed
// environment maps, replacing values with empty strings or placeholder hints.
package templater

import (
	"fmt"
	"sort"
	"strings"
)

// Options controls how the template is generated.
type Options struct {
	// Placeholder is written as the value for every key.
	// If empty, values are left blank (KEY=).
	Placeholder string

	// IncludeValues preserves original values as comments above each key.
	IncludeValues bool
}

// Generate produces a .env.example-style string from the provided env map.
// Keys are sorted alphabetically for deterministic output.
func Generate(env map[string]string, opts Options) string {
	if len(env) == 0 {
		return ""
	}

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for i, k := range keys {
		if i > 0 {
			sb.WriteByte('\n')
		}
		if opts.IncludeValues {
			fmt.Fprintf(&sb, "# %s=%s\n", k, env[k])
		}
		placeholder := opts.Placeholder
		fmt.Fprintf(&sb, "%s=%s\n", k, placeholder)
	}
	return sb.String()
}

// Merge combines multiple env maps into one before generating the template.
// Keys present in later maps do not overwrite earlier ones.
func Merge(envs []map[string]string) map[string]string {
	result := make(map[string]string)
	for _, env := range envs {
		for k, v := range env {
			if _, exists := result[k]; !exists {
				result[k] = v
			}
		}
	}
	return result
}
