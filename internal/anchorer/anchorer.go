// Package anchorer identifies a "canonical" reference key set from a
// collection of parsed env maps and flags any key that deviates from it.
package anchorer

import (
	"fmt"
	"sort"
	"strings"
)

// Deviation describes a single key that differs from the anchor set.
type Deviation struct {
	Key     string
	Env     string
	Reason  string
}

// Result holds the anchor key set and all deviations found.
type Result struct {
	AnchorKeys []string
	Deviations []Deviation
}

// Anchor takes a named map of env variable maps (envName -> key/value pairs)
// and a designated anchor environment name. It returns every key present in
// the anchor that is absent in another env, and every key present in another
// env that is absent from the anchor.
//
// If anchorEnv is empty the function picks the env with the most keys.
func Anchor(envs map[string]map[string]string, anchorEnv string) Result {
	if len(envs) == 0 {
		return Result{}
	}

	if anchorEnv == "" {
		anchorEnv = pickLargest(envs)
	}

	anchor, ok := envs[anchorEnv]
	if !ok {
		return Result{}
	}

	anchorKeys := sortedKeys(anchor)
	var deviations []Deviation

	for envName, env := range envs {
		if envName == anchorEnv {
			continue
		}
		// Keys in anchor missing from this env.
		for _, k := range anchorKeys {
			if _, found := env[k]; !found {
				deviations = append(deviations, Deviation{
					Key:    k,
					Env:    envName,
					Reason: fmt.Sprintf("missing (present in anchor %q)", anchorEnv),
				})
			}
		}
		// Keys in this env not in anchor.
		for _, k := range sortedKeys(env) {
			if _, found := anchor[k]; !found {
				deviations = append(deviations, Deviation{
					Key:    k,
					Env:    envName,
					Reason: fmt.Sprintf("extra (not in anchor %q)", anchorEnv),
				})
			}
		}
	}

	sort.Slice(deviations, func(i, j int) bool {
		if deviations[i].Key != deviations[j].Key {
			return deviations[i].Key < deviations[j].Key
		}
		return strings.Compare(deviations[i].Env, deviations[j].Env) < 0
	})

	return Result{
		AnchorKeys: anchorKeys,
		Deviations: deviations,
	}
}

func pickLargest(envs map[string]map[string]string) string {
	best, bestLen := "", -1
	for name, m := range envs {
		if len(m) > bestLen || (len(m) == bestLen && name < best) {
			best, bestLen = name, len(m)
		}
	}
	return best
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
