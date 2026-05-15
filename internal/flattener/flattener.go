// Package flattener merges multiple env maps into a single flat map,
// with configurable conflict resolution and key tracking.
package flattener

import "sort"

// Strategy controls how conflicting keys are resolved.
type Strategy int

const (
	// StrategyFirst keeps the value from the first env that defines the key.
	StrategyFirst Strategy = iota
	// StrategyLast keeps the value from the last env that defines the key.
	StrategyLast
)

// Origin records which environment a key's value came from.
type Origin struct {
	Key    string
	EnvName string
	Value  string
}

// Result holds the flattened env map and provenance information.
type Result struct {
	Env     map[string]string
	Origins []Origin
	Conflicts []string // keys that appeared in more than one env
}

// Options configures the Flatten operation.
type Options struct {
	Strategy Strategy
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{Strategy: StrategyFirst}
}

// Flatten merges the provided named env maps into a single flat map.
// envs is a slice of (name, map) pairs represented as named inputs.
func Flatten(envs []NamedEnv, opts Options) Result {
	result := Result{
		Env: make(map[string]string),
	}

	seen := make(map[string]string) // key -> first env name
	conflictSet := make(map[string]bool)

	for _, ne := range envs {
		for k, v := range ne.Env {
			if existing, ok := seen[k]; ok {
				conflictSet[k] = true
				if opts.Strategy == StrategyLast {
					result.Env[k] = v
					updateOrigin(&result.Origins, k, ne.Name, v)
				}
				_ = existing
			} else {
				seen[k] = ne.Name
				result.Env[k] = v
				result.Origins = append(result.Origins, Origin{
					Key:     k,
					EnvName: ne.Name,
					Value:   v,
				})
			}
		}
	}

	for k := range conflictSet {
		result.Conflicts = append(result.Conflicts, k)
	}
	sort.Strings(result.Conflicts)

	sort.Slice(result.Origins, func(i, j int) bool {
		return result.Origins[i].Key < result.Origins[j].Key
	})

	return result
}

// NamedEnv pairs an environment name with its key-value map.
type NamedEnv struct {
	Name string
	Env  map[string]string
}

func updateOrigin(origins *[]Origin, key, envName, value string) {
	for i, o := range *origins {
		if o.Key == key {
			(*origins)[i].EnvName = envName
			(*origins)[i].Value = value
			return
		}
	}
}
