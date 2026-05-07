// Package reporter provides formatted output for envdiff results.
//
// It supports multiple output formats for presenting the differences
// found between two .env files:
//
//   - text: human-readable plain text output with symbols to indicate
//     missing (-), extra (+), and mismatched (~) keys.
//
//   - json: machine-readable JSON output suitable for integration with
//     other tools or CI pipelines.
//
// Usage:
//
//	result := differ.Diff(left, right)
//	err := reporter.Report(os.Stdout, result, ".env.dev", ".env.prod", reporter.FormatText)
//
package reporter
