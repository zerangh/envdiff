// Package grouper groups environment keys by a common prefix delimiter,
// allowing callers to analyse which logical "namespaces" differ across envs.
package grouper

import (
	"sort"
	"strings"

	"github.com/yourusername/envdiff/internal/differ"
)

// Group represents a single prefix namespace and the keys that belong to it.
type Group struct {
	Prefix   string
	Keys     []string
	Missing  []string // keys missing in right env
	Mismatch []string // keys with differing values
}

// Result holds all groups produced from a diff.
type Result struct {
	Groups    []Group
	Ungrouped Group // keys that have no delimiter-based prefix
}

// ByPrefix partitions the keys in a differ.Result into groups based on the
// first segment when the key is split by delimiter (default "_").
func ByPrefix(d differ.Result, delimiter string) Result {
	if delimiter == "" {
		delimiter = "_"
	}

	groupMap := map[string]*Group{}

	allKeys := uniqueKeys(d)

	for _, key := range allKeys {
		prefix, isGrouped := extractPrefix(key, delimiter)
		g := ensureGroup(groupMap, prefix)
		g.Keys = append(g.Keys, key)
		if !isGrouped {
			// will be collected into Ungrouped later
			groupMap[""].Keys = append(groupMap[""].Keys, key)
		}
		_ = isGrouped
	}

	// annotate missing / mismatched
	for _, key := range d.MissingInRight {
		prefix, _ := extractPrefix(key, delimiter)
		if g, ok := groupMap[prefix]; ok {
			g.Missing = append(g.Missing, key)
		}
	}
	for _, m := range d.Mismatched {
		prefix, _ := extractPrefix(m.Key, delimiter)
		if g, ok := groupMap[prefix]; ok {
			g.Mismatch = append(g.Mismatch, m.Key)
		}
	}

	var groups []Group
	var ungrouped Group

	for prefix, g := range groupMap {
		if prefix == "" {
			ungrouped = *g
			ungrouped.Prefix = "(none)"
			continue
		}
		groups = append(groups, *g)
	}

	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Prefix < groups[j].Prefix
	})

	return Result{Groups: groups, Ungrouped: ungrouped}
}

func extractPrefix(key, delimiter string) (string, bool) {
	idx := strings.Index(key, delimiter)
	if idx <= 0 {
		return "", false
	}
	return key[:idx], true
}

func ensureGroup(m map[string]*Group, prefix string) *Group {
	if g, ok := m[prefix]; ok {
		return g
	}
	g := &Group{Prefix: prefix}
	m[prefix] = g
	return g
}

func uniqueKeys(d differ.Result) []string {
	seen := map[string]struct{}{}
	var keys []string
	for _, k := range d.MissingInRight {
		if _, ok := seen[k]; !ok {
			seen[k] = struct{}{}
			keys = append(keys, k)
		}
	}
	for _, k := range d.MissingInLeft {
		if _, ok := seen[k]; !ok {
			seen[k] = struct{}{}
			keys = append(keys, k)
		}
	}
	for _, m := range d.Mismatched {
		if _, ok := seen[m.Key]; !ok {
			seen[m.Key] = struct{}{}
			keys = append(keys, m.Key)
		}
	}
	sort.Strings(keys)
	return keys
}
