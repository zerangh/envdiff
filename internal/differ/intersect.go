package differ

import (
	"sort"

	"github.com/your-org/envdiff/internal/parser"
)

// Intersection holds keys that exist in all provided environments
// along with any value disagreements across them.
type Intersection struct {
	// CommonKeys lists keys present in every env map.
	CommonKeys []string

	// ValueMap maps each common key to the set of distinct values seen
	// across all envs (env-name → value).
	ValueMap map[string]map[string]string

	// Consistent contains keys whose value is identical in every env.
	Consistent []string

	// Divergent contains keys that have at least two distinct values.
	Divergent []string
}

// Intersect computes the set of keys common to all named envs and
// classifies them as consistent or divergent based on their values.
func Intersect(envs map[string]parser.Env) Intersection {
	if len(envs) == 0 {
		return Intersection{}
	}

	// Build a frequency map: key → count of envs containing it.
	freq := make(map[string]int)
	for _, env := range envs {
		for k := range env {
			freq[k]++
		}
	}

	total := len(envs)
	valueMap := make(map[string]map[string]string)

	for k, count := range freq {
		if count == total {
			valueMap[k] = make(map[string]string)
			for name, env := range envs {
				valueMap[k][name] = env[k]
			}
		}
	}

	var common, consistent, divergent []string
	for k, vals := range valueMap {
		common = append(common, k)
		if allSame(vals) {
			consistent = append(consistent, k)
		} else {
			divergent = append(divergent, k)
		}
	}

	sort.Strings(common)
	sort.Strings(consistent)
	sort.Strings(divergent)

	return Intersection{
		CommonKeys: common,
		ValueMap:   valueMap,
		Consistent: consistent,
		Divergent:  divergent,
	}
}

func allSame(vals map[string]string) bool {
	var ref string
	first := true
	for _, v := range vals {
		if first {
			ref = v
			first = false
			continue
		}
		if v != ref {
			return false
		}
	}
	return true
}
