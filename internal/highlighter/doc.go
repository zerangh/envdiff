// Package highlighter provides inline value-level diffing for mismatched
// environment keys.
//
// Given a differ.Result, Compute produces per-key Highlight values that
// break each mismatched value into Segments annotated as added, removed,
// or unchanged. This makes it easy for reporters and formatters to render
// precise, word-level change indicators rather than showing only the raw
// before/after strings.
//
// Example usage:
//
//	result := differ.Diff(left, right)
//	hr := highlighter.Compute(result)
//	for _, h := range hr.Highlights {
//		fmt.Print(highlighter.Format(h))
//	}
package highlighter
