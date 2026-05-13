// Package aliaser maps environment variable keys across different naming
// conventions (e.g. LEGACY_DB_HOST -> DATABASE_HOST) using a user-supplied
// alias map, and rewrites a parsed env map or diff result accordingly.
package aliaser

import (
	"sort"
	"strings"

	"github.com/user/envdiff/internal/differ"
)

// AliasMap maps old key names to their canonical replacements.
type AliasMap map[string]string

// ApplyToEnv returns a new env map with aliased keys renamed.
// Keys not present in the alias map are passed through unchanged.
// If both an old key and its alias exist in the env map, the alias wins.
func ApplyToEnv(env map[string]string, aliases AliasMap) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if canonical, ok := aliases[k]; ok {
			out[canonical] = v
		} else {
			out[k] = v
		}
	}
	return out
}

// ApplyToResult returns a copy of the differ.Result with aliased keys renamed
// in MissingInLeft, MissingInRight, and Mismatched slices.
func ApplyToResult(r differ.Result, aliases AliasMap) differ.Result {
	out := differ.Result{
		MissingInLeft:  renameKeys(r.MissingInLeft, aliases),
		MissingInRight: renameKeys(r.MissingInRight, aliases),
		Mismatched:     renameMismatched(r.Mismatched, aliases),
	}
	return out
}

// LoadAliasMap parses a flat list of "OLD=NEW" strings into an AliasMap.
// Lines that are blank or start with '#' are ignored.
func LoadAliasMap(lines []string) AliasMap {
	am := make(AliasMap)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		old := strings.TrimSpace(parts[0])
		new := strings.TrimSpace(parts[1])
		if old != "" && new != "" {
			am[old] = new
		}
	}
	return am
}

// Invert returns a new AliasMap with keys and values swapped.
func (am AliasMap) Invert() AliasMap {
	inv := make(AliasMap, len(am))
	for k, v := range am {
		inv[v] = k
	}
	return inv
}

func renameKeys(keys []string, aliases AliasMap) []string {
	out := make([]string, len(keys))
	for i, k := range keys {
		if canonical, ok := aliases[k]; ok {
			out[i] = canonical
		} else {
			out[i] = k
		}
	}
	sort.Strings(out)
	return out
}

func renameMismatched(pairs []differ.Mismatch, aliases AliasMap) []differ.Mismatch {
	out := make([]differ.Mismatch, len(pairs))
	for i, m := range pairs {
		if canonical, ok := aliases[m.Key]; ok {
			m.Key = canonical
		}
		out[i] = m
	}
	return out
}
