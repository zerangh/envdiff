package linter

import (
	"fmt"
	"strings"

	"github.com/user/envdiff/internal/differ"
)

// Issue represents a linting problem found in an env file.
type Issue struct {
	Key      string
	Severity string // "error" or "warning"
	Message  string
}

// Result holds all linting issues found.
type Result struct {
	Issues []Issue
}

// HasErrors returns true if any issue has severity "error".
func (r Result) HasErrors() bool {
	for _, issue := range r.Issues {
		if issue.Severity == "error" {
			return true
		}
	}
	return false
}

// Lint inspects a parsed env map and a diff result to produce linting issues.
func Lint(env map[string]string, diff differ.Result) Result {
	var issues []Issue

	for key, val := range env {
		// Warn on keys that are not uppercase
		if key != strings.ToUpper(key) {
			issues = append(issues, Issue{
				Key:      key,
				Severity: "warning",
				Message:  fmt.Sprintf("key %q is not uppercase", key),
			})
		}

		// Warn on keys with spaces in value
		if strings.TrimSpace(val) != val {
			issues = append(issues, Issue{
				Key:      key,
				Severity: "warning",
				Message:  fmt.Sprintf("key %q has leading or trailing whitespace in value", key),
			})
		}

		// Error on keys with newlines in value
		if strings.ContainsAny(val, "\n\r") {
			issues = append(issues, Issue{
				Key:      key,
				Severity: "error",
				Message:  fmt.Sprintf("key %q contains newline characters in value", key),
			})
		}
	}

	// Warn on keys missing from the right file
	for _, key := range diff.MissingInRight {
		issues = append(issues, Issue{
			Key:      key,
			Severity: "warning",
			Message:  fmt.Sprintf("key %q is missing in the right env file", key),
		})
	}

	return Result{Issues: issues}
}
