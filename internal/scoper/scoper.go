package scoper

import (
	"sort"
	"strings"

	"github.com/your-org/envdiff/internal/differ"
)

// Scope represents a named subset of environment keys filtered by a scope tag.
type Scope struct {
	Name string
	Keys []string
	Result differ.Result
}

// Options controls how scopes are extracted.
type Options struct {
	// ScopePrefix is the prefix used to identify scope annotations in key names.
	// e.g. "APP_" or "DB_". Defaults to empty (all keys belong to a single scope).
	ScopePrefix string

	// MinKeys is the minimum number of keys a scope must have to be included.
	MinKeys int
}

// DefaultOptions returns sensible defaults for scoping.
func DefaultOptions() Options {
	return Options{
		ScopePrefix: "",
		MinKeys:     1,
	}
}

// Extract partitions a differ.Result into named Scopes based on key prefixes.
// Each unique first-segment prefix (split by "_") becomes its own scope.
// Keys with no underscore are grouped under "default".
func Extract(result differ.Result, opts Options) []Scope {
	scopes := map[string]*Scope{}

	allKeys := uniqueKeys(result)

	for _, key := range allKeys {
		name := scopeName(key, opts.ScopePrefix)
		if _, ok := scopes[name]; !ok {
			scopes[name] = &Scope{Name: name}
		}
		scopes[name].Keys = append(scopes[name].Keys, key)
	}

	for name, sc := range scopes {
		sc.Result = filterResult(result, sc.Keys)
		scopes[name] = sc
	}

	var out []Scope
	for _, sc := range scopes {
		if len(sc.Keys) >= opts.MinKeys {
			out = append(out, *sc)
		}
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Name < out[j].Name
	})
	return out
}

func scopeName(key, prefix string) string {
	if prefix != "" && !strings.HasPrefix(key, prefix) {
		return "other"
	}
	parts := strings.SplitN(key, "_", 2)
	if len(parts) < 2 {
		return "default"
	}
	return strings.ToLower(parts[0])
}

func uniqueKeys(result differ.Result) []string {
	seen := map[string]struct{}{}
	for _, k := range result.MissingInRight {
		seen[k] = struct{}{}
	}
	for _, k := range result.MissingInLeft {
		seen[k] = struct{}{}
	}
	for _, m := range result.Mismatched {
		seen[m.Key] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func filterResult(result differ.Result, keys []string) differ.Result {
	set := map[string]struct{}{}
	for _, k := range keys {
		set[k] = struct{}{}
	}
	out := differ.Result{}
	for _, k := range result.MissingInRight {
		if _, ok := set[k]; ok {
			out.MissingInRight = append(out.MissingInRight, k)
		}
	}
	for _, k := range result.MissingInLeft {
		if _, ok := set[k]; ok {
			out.MissingInLeft = append(out.MissingInLeft, k)
		}
	}
	for _, m := range result.Mismatched {
		if _, ok := set[m.Key]; ok {
			out.Mismatched = append(out.Mismatched, m)
		}
	}
	return out
}
