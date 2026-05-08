// Package linter provides static analysis for .env file contents.
//
// It inspects parsed environment variable maps and differ results to
// surface issues such as non-uppercase keys, values with surrounding
// whitespace, values containing newline characters, and keys missing
// from one side of a diff.
//
// Issues are classified by severity:
//   - "warning": non-fatal style or consistency problems
//   - "error": structural problems that may cause runtime failures
//
// Example usage:
//
//	env, _ := parser.ParseFile(".env")
//	diff := differ.Diff(left, right)
//	result := linter.Lint(env, diff)
//	if result.HasErrors() {
//		log.Fatal("linting errors found")
//	}
package linter
