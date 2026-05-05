// Package differ provides functionality for comparing two parsed .env file
// maps and identifying discrepancies between them.
//
// The primary entry point is Diff, which accepts two map[string]string values
// (as produced by the parser package) and returns a Result describing:
//
//   - Keys present in the left file but missing from the right.
//   - Keys present in the right file but missing from the left.
//   - Keys present in both files whose values differ.
//
// Example usage:
//
//	left, _ := parser.ParseFile(".env.development")
//	right, _ := parser.ParseFile(".env.production")
//	result := differ.Diff(left, right)
//	if result.HasDiff() {
//		// handle differences
//	}
package differ
