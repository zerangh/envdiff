// Package differ compares two parsed .env maps and produces structured
// difference results.
//
// # Core Diff
//
// Use [Diff] to compare two env maps directly, or [DiffFiles] to parse and
// compare two files in one step. Both return a [Result] describing keys that
// are missing on either side and keys whose values differ.
//
// # Baseline Comparison
//
// [CompareToBaseline] and [CompareFilesToBaseline] compare a current env
// against a previously recorded baseline snapshot. The returned [BaselineResult]
// separates changes into Added, Removed, Changed, and Unchanged categories,
// making it easy to understand drift from a known-good state.
//
// # Changelog
//
// [BuildChangelog] produces a human-readable ordered list of changes between
// two env files, suitable for inclusion in release notes or audit logs.
package differ
