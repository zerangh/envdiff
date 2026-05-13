// Package deduplicator identifies and resolves duplicate or conflicting keys
// across multiple env maps that represent the same environment.
//
// When layering .env files (e.g. .env + .env.local) the same key may appear
// with different values. Deduplicator surfaces those conflicts as Duplicate
// entries and provides a Resolve helper that merges the maps into a single
// canonical map using a configurable Strategy (KeepFirst or KeepLast).
//
// It also integrates with differ.Result via FromResult, converting mismatched
// key/value pairs into Duplicate entries for unified reporting.
package deduplicator
