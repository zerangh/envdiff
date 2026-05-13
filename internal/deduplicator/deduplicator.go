package deduplicator

import (
	"sort"

	"github.com/user/envdiff/internal/differ"
)

// Strategy controls how duplicate keys are resolved when merging envs.
type Strategy int

const (
	// StrategyKeepFirst retains the first occurrence of a duplicate key.
	StrategyKeepFirst Strategy = iota
	// StrategyKeepLast retains the last occurrence of a duplicate key.
	StrategyKeepLast
)

// Duplicate describes a key that appears more than once in a single env map.
type Duplicate struct {
	Key    string
	Values []string
}

// Report holds the result of a deduplication pass.
type Report struct {
	Duplicates []Duplicate
	Clean      map[string]string
}

// Detect finds keys that have conflicting values across multiple env maps
// representing the same environment (e.g. layered .env files).
func Detect(envs []map[string]string) Report {
	type entry struct {
		values []string
		seen   map[string]struct{}
	}

	agg := make(map[string]*entry)

	for _, env := range envs {
		for k, v := range env {
			if _, ok := agg[k]; !ok {
				agg[k] = &entry{seen: make(map[string]struct{})}
			}
			if _, seen := agg[k].seen[v]; !seen {
				agg[k].values = append(agg[k].values, v)
				agg[k].seen[v] = struct{}{}
			}
		}
	}

	var dups []Duplicate
	for k, e := range agg {
		if len(e.values) > 1 {
			dups = append(dups, Duplicate{Key: k, Values: e.values})
		}
	}
	sort.Slice(dups, func(i, j int) bool { return dups[i].Key < dups[j].Key })

	return Report{Duplicates: dups}
}

// Resolve merges multiple env maps into one, resolving duplicates using the
// given strategy. It also populates Report.Clean with the resolved map.
func Resolve(envs []map[string]string, strategy Strategy) map[string]string {
	resolved := make(map[string]string)
	for _, env := range envs {
		for k, v := range env {
			if _, exists := resolved[k]; exists && strategy == StrategyKeepFirst {
				continue
			}
			resolved[k] = v
		}
	}
	return resolved
}

// FromResult extracts all unique values referenced in a differ.Result and
// returns any keys whose left/right values are both non-empty but differ,
// formatted as Duplicate entries for reporting.
func FromResult(result differ.Result) []Duplicate {
	var dups []Duplicate
	for _, m := range result.Mismatched {
		dups = append(dups, Duplicate{
			Key:    m.Key,
			Values: []string{m.Left, m.Right},
		})
	}
	sort.Slice(dups, func(i, j int) bool { return dups[i].Key < dups[j].Key })
	return dups
}
