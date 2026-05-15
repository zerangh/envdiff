// Package comparer provides flexible value comparison for env key-value pairs.
//
// It supports multiple comparison modes:
//
//   - Exact: byte-for-byte equality
//   - TrimSpace: ignores leading/trailing whitespace
//   - CaseInsensitive: case-folded comparison
//   - NormalizeBools: treats "true"/"1"/"yes" and "false"/"0"/"no" as equivalent
//
// The FilterMismatches function applies comparison options to a differ.Result,
// removing entries from the Mismatched list that are considered equivalent
// under the selected options. This is useful for reducing noise when comparing
// environments that use slightly different value conventions.
//
// Example:
//
//	opts := comparer.Options{TrimSpace: true, NormalizeBools: true}
//	cleaned := comparer.FilterMismatches(result, opts)
package comparer
