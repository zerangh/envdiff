// Package validator provides value validation for .env file entries,
// checking for common issues such as empty values, suspicious whitespace,
// or keys that appear in one environment but have no value set.
package validator

import (
	"fmt"
	"strings"

	"github.com/user/envdiff/internal/differ"
)

// Warning represents a non-fatal issue found during validation.
type Warning struct {
	Key     string
	Message string
}

// Result holds all warnings produced by a validation pass.
type Result struct {
	Warnings []Warning
}

// HasWarnings returns true if any warnings were recorded.
func (r Result) HasWarnings() bool {
	return len(r.Warnings) > 0
}

// Validate inspects a differ.Result and the raw parsed env maps for
// common value-level problems, returning a validator.Result with warnings.
func Validate(diff differ.Result, left, right map[string]string) Result {
	var warnings []Warning

	// Warn on keys present in both sides but with an empty value on either side.
	for _, m := range diff.Mismatched {
		lv := left[m.Key]
		rv := right[m.Key]
		if lv == "" {
			warnings = append(warnings, Warning{
				Key:     m.Key,
				Message: "left value is empty",
			})
		}
		if rv == "" {
			warnings = append(warnings, Warning{
				Key:     m.Key,
				Message: "right value is empty",
			})
		}
		if hasLeadingOrTrailingSpace(lv) || hasLeadingOrTrailingSpace(rv) {
			warnings = append(warnings, Warning{
				Key:     m.Key,
				Message: "value contains leading or trailing whitespace",
			})
		}
	}

	// Warn on keys missing from right that have an empty value on the left.
	for _, key := range diff.MissingInRight {
		if left[key] == "" {
			warnings = append(warnings, Warning{
				Key:     key,
				Message: "key is missing in right and left value is also empty",
			})
		}
	}

	// Warn on keys missing from left that have an empty value on the right.
	for _, key := range diff.MissingInLeft {
		if right[key] == "" {
			warnings = append(warnings, Warning{
				Key:     key,
				Message: "key is missing in left and right value is also empty",
			})
		}
	}

	return Result{Warnings: warnings}
}

// Format returns a human-readable summary of all warnings.
func (r Result) Format() string {
	if !r.HasWarnings() {
		return "no validation warnings\n"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d validation warning(s):\n", len(r.Warnings)))
	for _, w := range r.Warnings {
		sb.WriteString(fmt.Sprintf("  [%s] %s\n", w.Key, w.Message))
	}
	return sb.String()
}

func hasLeadingOrTrailingSpace(s string) bool {
	return len(s) > 0 && (s != strings.TrimSpace(s))
}
