// Package summarizer provides a high-level overview of a diff result,
// combining key counts, validation warnings, and a health score into a
// single Summary value that can be formatted for display.
//
// Usage:
//
//	result := differ.Diff(left, right)
//	summary := summarizer.Summarize(result)
//	fmt.Print(summary.Format())
//
// The Healthy field is true when the computed score is 80 or above,
// giving a quick pass/fail signal suitable for CI pipelines.
package summarizer
