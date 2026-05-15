// Package resolver provides functionality for resolving missing keys
// across environments by suggesting values from other available environments.
package resolver

import (
	"sort"

	"github.com/user/envdiff/internal/differ"
)

// Suggestion represents a resolved value candidate for a missing key.
type Suggestion struct {
	Key    string
	Value  string
	Source string // which environment the value came from
}

// Result holds all suggestions produced by the resolver.
type Result struct {
	Suggestions []Suggestion
	Unresolved  []string // keys with no candidate in any env
}

// Resolve inspects missing keys from a diff result and attempts to find
// candidate values from the provided environment map (name -> key/value pairs).
// It returns a Result with suggestions and any keys that remain unresolved.
func Resolve(diff differ.Result, envs map[string]map[string]string) Result {
	missing := collectMissingKeys(diff)

	var suggestions []Suggestion
	var unresolved []string

	for _, key := range missing {
		found := false
		for envName, env := range envs {
			if val, ok := env[key]; ok && val != "" {
				suggestions = append(suggestions, Suggestion{
					Key:    key,
					Value:  val,
					Source: envName,
				})
				found = true
				break
			}
		}
		if !found {
			unresolved = append(unresolved, key)
		}
	}

	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Key < suggestions[j].Key
	})
	sort.Strings(unresolved)

	return Result{
		Suggestions: suggestions,
		Unresolved:  unresolved,
	}
}

// collectMissingKeys gathers all keys missing from either side of the diff.
func collectMissingKeys(diff differ.Result) []string {
	seen := make(map[string]struct{})
	for _, k := range diff.MissingInRight {
		seen[k] = struct{}{}
	}
	for _, k := range diff.MissingInLeft {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
