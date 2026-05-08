// Package exporter provides functionality to export envdiff results
// to various output formats beyond the default terminal reporter.
//
// Supported formats:
//
//   - env      — plain KEY=value lines for mismatched keys (left-side values)
//   - markdown — a Markdown report with tables and sections
//   - json     — a structured JSON object including file names
//
// Usage:
//
//	exporter.Export(os.Stdout, result, "prod.env", "staging.env", exporter.FormatMarkdown)
package exporter
