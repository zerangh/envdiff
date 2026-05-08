package main

import (
	"fmt"
	"io"
	"os"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/linter"
	"github.com/user/envdiff/internal/parser"
)

// runLint parses the left env file, diffs it against the right, then
// runs the linter and prints any issues found.
func runLint(args []string, out io.Writer) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envdiff lint <left> <right>")
	}

	left, err := parser.ParseFile(args[0])
	if err != nil {
		return fmt.Errorf("parsing %s: %w", args[0], err)
	}

	right, err := parser.ParseFile(args[1])
	if err != nil {
		return fmt.Errorf("parsing %s: %w", args[1], err)
	}

	diff := differ.Diff(left, right)
	result := linter.Lint(left, diff)

	if len(result.Issues) == 0 {
		fmt.Fprintln(out, "No linting issues found.")
		return nil
	}

	for _, issue := range result.Issues {
		fmt.Fprintf(out, "[%s] %s\n", issue.Severity, issue.Message)
	}

	if result.HasErrors() {
		os.Exit(1)
	}
	return nil
}
