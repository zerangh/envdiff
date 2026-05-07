// Package filter provides post-processing utilities to narrow down the
// results produced by the differ package.
//
// After running a diff between two .env files, callers may want to focus
// on a subset of keys — for example, only keys that share a common prefix
// such as "DB_" or "APP_", or to ignore well-known keys that are expected
// to differ between environments.
//
// Usage:
//
//	opts := filter.Options{
//		Prefix:      "DB_",
//		Exclude:     []string{"DB_PASSWORD"},
//		OnlyMissing: false,
//	}
//	filtered := filter.Apply(diffResult, opts)
//
// The Apply function is non-destructive — it always returns a new
// DiffResult and never modifies the original.
package filter
