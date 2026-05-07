package reporter

import (
	"fmt"
	"io"
	"os"

	"github.com/user/envdiff/internal/differ"
)

// colorCode maps semantic meaning to ANSI escape codes.
var colorCode = struct {
	Red    string
	Yellow string
	Green  string
	Reset  string
}{
	Red:    "\033[31m",
	Yellow: "\033[33m",
	Green:  "\033[32m",
	Reset:  "\033[0m",
}

// isTerminal reports whether w is a terminal that supports ANSI color.
func isTerminal(w io.Writer) bool {
	if f, ok := w.(*os.File); ok {
		stat, err := f.Stat()
		if err == nil && (stat.Mode()&os.ModeCharDevice) != 0 {
			return true
		}
	}
	return false
}

// reportColor writes a colorized, human-readable diff report to w.
// Missing keys are shown in red, mismatched values in yellow, and a
// clean result in green.
func reportColor(w io.Writer, result differ.Result) {
	if len(result.MissingInRight) == 0 &&
		len(result.MissingInLeft) == 0 &&
		len(result.Mismatched) == 0 {
		fmt.Fprintf(w, "%s✔ No differences found.%s\n",
			colorCode.Green, colorCode.Reset)
		return
	}

	for _, key := range result.MissingInRight {
		fmt.Fprintf(w, "%s- MISSING IN RIGHT: %s%s\n",
			colorCode.Red, key, colorCode.Reset)
	}
	for _, key := range result.MissingInLeft {
		fmt.Fprintf(w, "%s- MISSING IN LEFT:  %s%s\n",
			colorCode.Red, key, colorCode.Reset)
	}
	for _, m := range result.Mismatched {
		fmt.Fprintf(w, "%s~ MISMATCH: %s\n  left:  %q\n  right: %q%s\n",
			colorCode.Yellow, m.Key, m.LeftValue, m.RightValue, colorCode.Reset)
	}
}
