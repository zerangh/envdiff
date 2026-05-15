// Package evaluator applies rule-based evaluation to diff results,
// producing structured findings with severity levels.
//
// Each finding maps to a named rule (e.g. "missing-in-right", "value-mismatch")
// and carries a severity of "error", "warn", or "info".
//
// Usage:
//
//	result := differ.Diff(left, right)
//	report := evaluator.Evaluate(result)
//	for _, f := range report.Findings {
//		fmt.Println(f)
//	}
//	if report.HasErrors {
//		os.Exit(1)
//	}
package evaluator
