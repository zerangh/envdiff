package profiler

import (
	"sort"

	"github.com/user/envdiff/internal/differ"
)

// KeyProfile holds statistics about a key across multiple environments.
type KeyProfile struct {
	Key          string
	PresentIn    []string
	MissingFrom  []string
	UniqueValues int
	Values       map[string]string // env name -> value
}

// Report is the result of profiling keys across environments.
type Report struct {
	Profiles    []KeyProfile
	EnvNames    []string
	TotalKeys   int
	Consistent  int
	Inconsistent int
}

// Profile analyses a set of named env maps and produces a cross-environment
// key profile report.
func Profile(envs map[string]map[string]string) Report {
	keySet := map[string]struct{}{}
	for _, env := range envs {
		for k := range env {
			keySet[k] = struct{}{}
		}
	}

	envNames := make([]string, 0, len(envs))
	for name := range envs {
		envNames = append(envNames, name)
	}
	sort.Strings(envNames)

	profiles := make([]KeyProfile, 0, len(keySet))
	for key := range keySet {
		kp := KeyProfile{
			Key:    key,
			Values: make(map[string]string),
		}
		valueSet := map[string]struct{}{}
		for _, name := range envNames {
			env := envs[name]
			if val, ok := env[key]; ok {
				kp.PresentIn = append(kp.PresentIn, name)
				kp.Values[name] = val
				valueSet[val] = struct{}{}
			} else {
				kp.MissingFrom = append(kp.MissingFrom, name)
			}
		}
		kp.UniqueValues = len(valueSet)
		profiles = append(profiles, kp)
	}

	sort.Slice(profiles, func(i, j int) bool {
		return profiles[i].Key < profiles[j].Key
	})

	consistent := 0
	for _, p := range profiles {
		if len(p.MissingFrom) == 0 && p.UniqueValues <= 1 {
			consistent++
		}
	}

	return Report{
		Profiles:     profiles,
		EnvNames:     envNames,
		TotalKeys:    len(profiles),
		Consistent:   consistent,
		Inconsistent: len(profiles) - consistent,
	}
}

// FromDiffs builds a multi-env map from a slice of named differ.Result values.
func FromDiffs(left map[string]string, rights map[string]map[string]string) map[string]map[string]string {
	_ = differ.Result{} // ensure import used via type reference in callers
	result := map[string]map[string]string{"left": left}
	for name, env := range rights {
		result[name] = env
	}
	return result
}
