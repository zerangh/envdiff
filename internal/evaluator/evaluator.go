package evaluator

import (
	"fmt"
	"strings"

	"github.com/your/envdiff/internal/differ"
)

// Rule represents a single evaluation rule applied to a diff result.
type Rule struct {
	Name    string
	Severity string // "error", "warn", "info"
	Message string
}

// Finding is a rule violation found during evaluation.
type Finding struct {
	Rule     Rule
	Key      string
	Detail   string
}

// Report holds all findings from an evaluation run.
type Report struct {
	Findings []Finding
	HasErrors bool
	HasWarnings bool
}

// String formats a single finding for display.
func (f Finding) String() string {
	if f.Detail != "" {
		return fmt.Sprintf("[%s] %s: %s (%s)", strings.ToUpper(f.Rule.Severity), f.Rule.Name, f.Key, f.Detail)
	}
	return fmt.Sprintf("[%s] %s: %s", strings.ToUpper(f.Rule.Severity), f.Rule.Name, f.Key)
}

// Evaluate runs all built-in rules against a diff result and returns a Report.
func Evaluate(result differ.Result) Report {
	var findings []Finding

	for _, key := range result.MissingInRight {
		findings = append(findings, Finding{
			Rule:   Rule{Name: "missing-in-right", Severity: "error", Message: "Key present in left but missing in right"},
			Key:    key,
			Detail: "not found in right file",
		})
	}

	for _, key := range result.MissingInLeft {
		findings = append(findings, Finding{
			Rule:   Rule{Name: "missing-in-left", Severity: "warn", Message: "Key present in right but missing in left"},
			Key:    key,
			Detail: "not found in left file",
		})
	}

	for _, mm := range result.Mismatched {
		findings = append(findings, Finding{
			Rule:   Rule{Name: "value-mismatch", Severity: "warn", Message: "Key exists in both files but values differ"},
			Key:    mm.Key,
			Detail: fmt.Sprintf("%q vs %q", mm.Left, mm.Right),
		})
	}

	report := Report{Findings: findings}
	for _, f := range findings {
		switch f.Rule.Severity {
		case "error":
			report.HasErrors = true
		case "warn":
			report.HasWarnings = true
		}
	}
	return report
}
