// Package encoder converts env maps to formatted string representations.
//
// Three output formats are supported:
//
//   - dotenv   — standard KEY=VALUE format compatible with most .env parsers
//   - exports  — shell-compatible format prefixing each line with "export"
//   - docker   — Docker --env-file compatible format (no quoting)
//
// Example usage:
//
//	opts := encoder.DefaultOptions()
//	text, err := encoder.Encode(env, opts)
//
Options allow controlling whether empty values are omitted, whether all values
are quoted, and whether keys are sorted before output.
package encoder
