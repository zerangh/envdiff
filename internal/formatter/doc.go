// Package formatter renders parsed env maps and diff results into
// formatted text lines suitable for writing back to .env files or
// passing to other tools.
//
// Three output styles are supported:
//
//   - StylePlain  — bare KEY=VALUE lines (default)
//   - StyleExport — lines prefixed with "export " for shell sourcing
//   - StyleDocker — alias for StylePlain, intended for docker --env-file
//
// Additional options control key ordering, omission of empty values,
// and automatic quoting of values that contain whitespace or special
// characters.
package formatter
