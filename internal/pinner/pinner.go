// Package pinner provides functionality to pin a set of required keys
// and verify that all pinned keys are present and non-empty across environments.
package pinner

import (
	"fmt"
	"sort"
)

// PinResult holds the outcome of a pin check for a single key.
type PinResult struct {
	Key     string
	Present bool
	Empty   bool
	Env     string
}

// Report summarises the result of checking pinned keys.
type Report struct {
	PinnedKeys []string
	Missing    []PinResult
	Empty      []PinResult
}

// OK returns true when every pinned key is present and non-empty.
func (r Report) OK() bool {
	return len(r.Missing) == 0 && len(r.Empty) == 0
}

// Format returns a human-readable summary of the report.
func (r Report) Format() string {
	if r.OK() {
		return fmt.Sprintf("all %d pinned keys are present and non-empty", len(r.PinnedKeys))
	}
	out := fmt.Sprintf("%d pinned key(s) have issues:\n", len(r.Missing)+len(r.Empty))
	for _, m := range r.Missing {
		out += fmt.Sprintf("  MISSING  [%s] %s\n", m.Env, m.Key)
	}
	for _, e := range r.Empty {
		out += fmt.Sprintf("  EMPTY    [%s] %s\n", e.Env, e.Key)
	}
	return out
}

// CheckEnv verifies that every key in pins exists and is non-empty in env.
// envName is used for labelling results.
func CheckEnv(envName string, env map[string]string, pins []string) Report {
	sorted := make([]string, len(pins))
	copy(sorted, pins)
	sort.Strings(sorted)

	r := Report{PinnedKeys: sorted}
	for _, key := range sorted {
		val, ok := env[key]
		if !ok {
			r.Missing = append(r.Missing, PinResult{Key: key, Present: false, Env: envName})
			continue
		}
		if val == "" {
			r.Empty = append(r.Empty, PinResult{Key: key, Present: true, Empty: true, Env: envName})
		}
	}
	return r
}

// CheckAll runs CheckEnv against every environment in envs and merges results.
func CheckAll(envs map[string]map[string]string, pins []string) Report {
	sorted := make([]string, len(pins))
	copy(sorted, pins)
	sort.Strings(sorted)

	merged := Report{PinnedKeys: sorted}
	names := make([]string, 0, len(envs))
	for name := range envs {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		r := CheckEnv(name, envs[name], pins)
		merged.Missing = append(merged.Missing, r.Missing...)
		merged.Empty = append(merged.Empty, r.Empty...)
	}
	return merged
}
