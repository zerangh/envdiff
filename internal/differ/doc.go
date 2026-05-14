// Package differ provides utilities for comparing two parsed .env maps.
//
// Core diff functionality:
//
//	result := differ.Diff(leftMap, rightMap)
//
// File-based diffing:
//
//	result, err := differ.DiffFiles(".env.staging", ".env.production")
//
// Changelog generation:
//
// The Changelog type converts a diff Result into an ordered list of
// ChangeEntry records, each tagged with a ChangeKind (added, removed,
// mismatch) and a timestamp. This is useful for audit logs and export.
//
//	cl := differ.BuildChangelog(result)
//	for _, entry := range cl.Entries {
//		fmt.Printf("%s %s\n", entry.Kind, entry.Key)
//	}
package differ
