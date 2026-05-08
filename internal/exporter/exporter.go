package exporter

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/differ"
)

// Format represents the output format for exporting.
type Format string

const (
	FormatEnv      Format = "env"
	FormatMarkdown Format = "markdown"
	FormatJSON     Format = "json"
)

// Export writes the diff result in the requested format to w.
func Export(w io.Writer, result differ.Result, leftName, rightName string, format Format) error {
	switch format {
	case FormatEnv:
		return exportEnv(w, result)
	case FormatMarkdown:
		return exportMarkdown(w, result, leftName, rightName)
	case FormatJSON:
		return exportJSONFull(w, result, leftName, rightName)
	default:
		return fmt.Errorf("unknown export format: %q", format)
	}
}

// exportEnv writes only keys that are present in both files (no missing)
// as KEY=value lines, sorted alphabetically.
func exportEnv(w io.Writer, result differ.Result) error {
	keys := make([]string, 0, len(result.Mismatched))
	for _, m := range result.Mismatched {
		keys = append(keys, m.Key)
	}
	sort.Strings(keys)
	for _, k := range keys {
		for _, m := range result.Mismatched {
			if m.Key == k {
				_, err := fmt.Fprintf(w, "%s=%s\n", m.Key, m.LeftValue)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// exportMarkdown writes a Markdown table summarising the diff.
func exportMarkdown(w io.Writer, result differ.Result, leftName, rightName string) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# envdiff report: `%s` vs `%s`\n\n", leftName, rightName))

	if len(result.MissingInRight)+len(result.MissingInLeft)+len(result.Mismatched) == 0 {
		sb.WriteString("✅ No differences found.\n")
		_, err := fmt.Fprint(w, sb.String())
		return err
	}

	if len(result.MissingInRight) > 0 {
		sb.WriteString(fmt.Sprintf("## Missing in `%s`\n\n", rightName))
		for _, k := range result.MissingInRight {
			sb.WriteString(fmt.Sprintf("- `%s`\n", k))
		}
		sb.WriteString("\n")
	}

	if len(result.MissingInLeft) > 0 {
		sb.WriteString(fmt.Sprintf("## Missing in `%s`\n\n", leftName))
		for _, k := range result.MissingInLeft {
			sb.WriteString(fmt.Sprintf("- `%s`\n", k))
		}
		sb.WriteString("\n")
	}

	if len(result.Mismatched) > 0 {
		sb.WriteString("## Mismatched values\n\n")
		sb.WriteString(fmt.Sprintf("| Key | `%s` | `%s` |\n", leftName, rightName))
		sb.WriteString("|-----|------|------\n")
		for _, m := range result.Mismatched {
			sb.WriteString(fmt.Sprintf("| `%s` | `%s` | `%s` |\n", m.Key, m.LeftValue, m.RightValue))
		}
		sb.WriteString("\n")
	}

	_, err := fmt.Fprint(w, sb.String())
	return err
}

// exportJSONFull writes a structured JSON object with file names included.
func exportJSONFull(w io.Writer, result differ.Result, leftName, rightName string) error {
	payload := map[string]interface{}{
		"left":           leftName,
		"right":          rightName,
		"missing_right":  result.MissingInRight,
		"missing_left":   result.MissingInLeft,
		"mismatched":     result.Mismatched,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}
