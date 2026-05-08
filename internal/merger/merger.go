// Package merger provides functionality to merge multiple .env files
// into a single unified map, with configurable conflict resolution strategies.
package merger

import "fmt"

// Strategy defines how key conflicts are resolved when merging.
type Strategy int

const (
	// StrategyFirst keeps the value from the first file that defines the key.
	StrategyFirst Strategy = iota
	// StrategyLast overwrites with the value from the last file that defines the key.
	StrategyLast
)

// Result holds the merged environment map and metadata about conflicts.
type Result struct {
	// Merged is the final key-value map after merging.
	Merged map[string]string
	// Conflicts maps each conflicting key to the list of values seen across files.
	Conflicts map[string][]string
	// Sources maps each key to the file path it was taken from.
	Sources map[string]string
}

// Merge combines multiple parsed env maps into a single Result.
// The files parameter is a slice of (path, map) pairs in order.
// Strategy controls how conflicts are resolved.
func Merge(files []NamedMap, strategy Strategy) (*Result, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("merger: at least one file is required")
	}

	result := &Result{
		Merged:    make(map[string]string),
		Conflicts: make(map[string][]string),
		Sources:   make(map[string]string),
	}

	for _, nmap := range files {
		for key, val := range nmap.Env {
			existing, exists := result.Merged[key]
			if !exists {
				result.Merged[key] = val
				result.Sources[key] = nmap.Path
				continue
			}
			// Record conflict regardless of strategy.
			if len(result.Conflicts[key]) == 0 {
				result.Conflicts[key] = []string{existing}
			}
			result.Conflicts[key] = append(result.Conflicts[key], val)

			if strategy == StrategyLast {
				result.Merged[key] = val
				result.Sources[key] = nmap.Path
			}
		}
	}

	return result, nil
}

// NamedMap associates a file path with its parsed env map.
type NamedMap struct {
	Path string
	Env  map[string]string
}
