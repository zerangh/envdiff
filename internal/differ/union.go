package differ

import (
	"sort"

	"github.com/your-org/envdiff/internal/parser"
)

// UnionResult holds the merged superset of keys across all environments,
// along with per-key presence and value information.
type UnionResult struct {
	// AllKeys is the sorted union of every key found across all envs.
	AllKeys []string

	// Values maps key -> (envName -> value). Missing entries mean the key
	// was absent in that environment.
	Values map[string]map[string]string

	// Envs is the sorted list of environment names included in the union.
	Envs []string
}

// Union computes the superset of keys across all provided named environments.
// Each map entry is envName -> parsed key/value map.
func Union(envs map[string]map[string]string) UnionResult {
	keySet := map[string]struct{}{}
	envNames := make([]string, 0, len(envs))

	for name, kv := range envs {
		envNames = append(envNames, name)
		for k := range kv {
			keySet[k] = struct{}{}
		}
	}

	sort.Strings(envNames)

	allKeys := make([]string, 0, len(keySet))
	for k := range keySet {
		allKeys = append(allKeys, k)
	}
	sort.Strings(allKeys)

	values := make(map[string]map[string]string, len(allKeys))
	for _, k := range allKeys {
		values[k] = make(map[string]string, len(envNames))
		for _, name := range envNames {
			if v, ok := envs[name][k]; ok {
				values[k][name] = v
			}
		}
	}

	return UnionResult{
		AllKeys: allKeys,
		Values:  values,
		Envs:    envNames,
	}
}

// UnionFiles parses multiple .env files by path and delegates to Union.
// The map key is the file path used as the environment name.
func UnionFiles(paths []string) (UnionResult, error) {
	envs := make(map[string]map[string]string, len(paths))
	for _, p := range paths {
		kv, err := parser.ParseFile(p)
		if err != nil {
			return UnionResult{}, err
		}
		envs[p] = kv
	}
	return Union(envs), nil
}
