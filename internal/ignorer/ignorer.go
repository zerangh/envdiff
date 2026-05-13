// Package ignorer provides support for loading and applying .envdiffignore
// files, allowing users to suppress specific keys from diff results.
package ignorer

import (
	"bufio"
	"os"
	"strings"

	"github.com/user/envdiff/internal/differ"
)

// Rules holds a set of key patterns to ignore during diffing.
type Rules struct {
	keys map[string]struct{}
}

// LoadFile reads an ignore file from the given path. Each non-blank,
// non-comment line is treated as a key name to ignore.
func LoadFile(path string) (*Rules, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rules := &Rules{keys: make(map[string]struct{})}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		rules.keys[line] = struct{}{}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return rules, nil
}

// NewRules constructs a Rules set from a slice of key names.
func NewRules(keys []string) *Rules {
	r := &Rules{keys: make(map[string]struct{}, len(keys))}
	for _, k := range keys {
		r.keys[k] = struct{}{}
	}
	return r
}

// Has reports whether the given key is in the ignore list.
func (r *Rules) Has(key string) bool {
	_, ok := r.keys[key]
	return ok
}

// Apply removes ignored keys from a differ.Result and returns a new Result.
func (r *Rules) Apply(result differ.Result) differ.Result {
	filtered := differ.Result{}

	for _, k := range result.MissingInRight {
		if !r.Has(k) {
			filtered.MissingInRight = append(filtered.MissingInRight, k)
		}
	}
	for _, k := range result.MissingInLeft {
		if !r.Has(k) {
			filtered.MissingInLeft = append(filtered.MissingInLeft, k)
		}
	}
	for _, m := range result.Mismatched {
		if !r.Has(m.Key) {
			filtered.Mismatched = append(filtered.Mismatched, m)
		}
	}
	return filtered
}
