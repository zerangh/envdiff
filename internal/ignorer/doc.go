// Package ignorer implements support for .envdiffignore files.
//
// An ignore file contains one key name per line. Lines beginning with '#'
// and blank lines are skipped. Any key listed in the file is excluded from
// diff results, suppressing both missing-key and mismatched-value findings.
//
// File lookup follows the same precedence as .gitignore: if no explicit path
// is given, [LoadFile] searches the current directory and each parent directory
// up to the filesystem root, stopping at the first file found.
//
// Example .envdiffignore:
//
//	# Keys managed by CI — ignore in local diffs
//	CI_TOKEN
//	DEPLOY_KEY
//
// Usage:
//
//	rules, err := ignorer.LoadFile(".envdiffignore")
//	if err != nil { ... }
//	clean := rules.Apply(diffResult)
//
// If no ignore file is found, [LoadFile] returns an empty [Rules] value and a
// nil error, so callers do not need to special-case the missing-file scenario.
package ignorer
