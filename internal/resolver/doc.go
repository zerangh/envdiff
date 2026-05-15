// Package resolver provides value resolution for missing environment keys.
//
// When comparing two .env files, some keys may be absent from one side.
// The resolver searches a set of reference environments (e.g. staging, prod)
// to find candidate values for those missing keys, helping teams quickly
// populate gaps without manual lookup.
//
// Usage:
//
//	envs := map[string]map[string]string{
//		"staging": stagingEnv,
//		"prod":    prodEnv,
//	}
//	res := resolver.Resolve(diffResult, envs)
//	for _, s := range res.Suggestions {
//		fmt.Printf("%s=%s  (from %s)\n", s.Key, s.Value, s.Source)
//	}
//	for _, k := range res.Unresolved {
//		fmt.Printf("%s has no candidate value\n", k)
//	}
package resolver
