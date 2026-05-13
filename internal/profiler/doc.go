// Package profiler provides cross-environment key profiling for .env files.
//
// Given a set of named environment maps, Profile produces a Report that
// summarises which keys are present, missing, or have diverging values
// across the supplied environments.
//
// Typical usage:
//
//	envs := map[string]map[string]string{
//		"dev":  devMap,
//		"prod": prodMap,
//	}
//	report := profiler.Profile(envs)
//	for _, p := range report.Profiles {
//		fmt.Println(p.Key, p.UniqueValues, p.MissingFrom)
//	}
package profiler
