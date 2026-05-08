package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/exporter"
	"github.com/user/envdiff/internal/parser"
)

// runExport handles the `envdiff export` sub-command.
// Usage: envdiff export <left> <right> --format=<env|markdown|json> [--output=<file>]
func runExport(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("export requires two .env file arguments")
	}

	leftPath := args[0]
	rightPath := args[1]

	format := exporter.FormatMarkdown
	outputPath := ""

	for _, arg := range args[2:] {
		switch {
		case strings.HasPrefix(arg, "--format="):
			format = exporter.Format(strings.TrimPrefix(arg, "--format="))
		case strings.HasPrefix(arg, "--output="):
			outputPath = strings.TrimPrefix(arg, "--output=")
		}
	}

	leftEnv, err := parser.ParseFile(leftPath)
	if err != nil {
		return fmt.Errorf("parsing %s: %w", leftPath, err)
	}

	rightEnv, err := parser.ParseFile(rightPath)
	if err != nil {
		return fmt.Errorf("parsing %s: %w", rightPath, err)
	}

	result := differ.Diff(leftEnv, rightEnv)

	w := os.Stdout
	if outputPath != "" {
		f, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("creating output file: %w", err)
		}
		defer f.Close()
		w = f
	}

	leftName := filepath.Base(leftPath)
	rightName := filepath.Base(rightPath)

	if err := exporter.Export(w, result, leftName, rightName, format); err != nil {
		return fmt.Errorf("exporting: %w", err)
	}

	if outputPath != "" {
		fmt.Fprintf(os.Stderr, "exported to %s\n", outputPath)
	}
	return nil
}
