// Package sorter provides flexible ordering strategies for env key maps and
// differ.Result values produced by the envdiff toolchain.
//
// Three strategies are available:
//
//   - Alpha: strict alphabetical order by key name (default).
//   - ChangeType: keys are grouped into buckets — missing-in-right,
//     missing-in-left, mismatched — and sorted alphabetically within each
//     bucket. Useful when rendering reports that highlight change severity.
//   - Prefix: keys are grouped by the portion of the name before the first
//     underscore (e.g. DB_, AWS_, APP_), then sorted alphabetically within
//     each prefix group. Useful for large .env files with many logical
//     subsystems.
//
// All operations return new values and do not mutate their inputs.
package sorter
