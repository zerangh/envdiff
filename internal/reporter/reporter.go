package reporter

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/envdiff/internal/differ"
)

// Format represents the output format for the report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Report writes a diff result to the given writer in the specified format.
func Report(w io.Writer, result differ.Result, leftName, rightName string, format Format) error {
	switch format {
	case FormatJSON:
		return reportJSON(w, result, leftName, rightName)
	default:
		return reportText(w, result, leftName, rightName)
	}
}

func reportText(w io.Writer, result differ.Result, leftName, rightName string) error {
	if result.IsClean() {
		_, err := fmt.Fprintln(w, "✓ No differences found.")
		return err
	}

	var sb strings.Builder

	if len(result.MissingInRight) > 0 {
		sb.WriteString(fmt.Sprintf("Keys present in %s but missing in %s:\n", leftName, rightName))
		for _, k := range result.MissingInRight {
			sb.WriteString(fmt.Sprintf("  - %s\n", k))
		}
	}

	if len(result.MissingInLeft) > 0 {
		sb.WriteString(fmt.Sprintf("Keys present in %s but missing in %s:\n", rightName, leftName))
		for _, k := range result.MissingInLeft {
			sb.WriteString(fmt.Sprintf("  + %s\n", k))
		}
	}

	if len(result.Mismatched) > 0 {
		sb.WriteString("Mismatched values:\n")
		for _, m := range result.Mismatched {
			sb.WriteString(fmt.Sprintf("  ~ %s\n", m.Key))
			sb.WriteString(fmt.Sprintf("      %s: %q\n", leftName, m.LeftValue))
			sb.WriteString(fmt.Sprintf("      %s: %q\n", rightName, m.RightValue))
		}
	}

	_, err := fmt.Fprint(w, sb.String())
	return err
}
