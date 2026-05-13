// Package scoper partitions a differ.Result into named Scopes based on
// key-name prefixes. Each unique first segment of an underscore-delimited
// key becomes its own Scope, allowing callers to analyse or report on
// logical sub-sections of an environment file independently.
//
// Example usage:
//
//	result := differ.Diff(left, right)
//	scopes := scoper.Extract(result, scoper.DefaultOptions())
//	for _, s := range scopes {
//		fmt.Printf("Scope %s: %d keys\n", s.Name, len(s.Keys))
//	}
package scoper
