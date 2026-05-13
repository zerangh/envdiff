// Package ignorer implements support for .envdiffignore files.
//
// An ignore file contains one key name per line. Lines beginning with '#'
// and blank lines are skipped. Any key listed in the file is excluded from
// diff results, suppressing both missing-key and mismatched-value findings.
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
package ignorer
