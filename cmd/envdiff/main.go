// Command envdiff compares two .env files and reports differences.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/reporter"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		format  = flag.String("format", "text", "output format: text or json")
		leftLabel  = flag.String("left-label", "", "label for the left file (defaults to filename)")
		rightLabel = flag.String("right-label", "", "label for the right file (defaults to filename)")
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: envdiff [flags] <left.env> <right.env>\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		flag.Usage()
		return fmt.Errorf("exactly two .env files are required")
	}

	leftPath, rightPath := args[0], args[1]

	if *leftLabel == "" {
		*leftLabel = leftPath
	}
	if *rightLabel == "" {
		*rightLabel = rightPath
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

	if err := reporter.Report(os.Stdout, result, *format, *leftLabel, *rightLabel); err != nil {
		return fmt.Errorf("reporting: %w", err)
	}

	// Exit with code 1 if differences were found.
	if len(result.MissingInRight) > 0 || len(result.MissingInLeft) > 0 || len(result.Mismatched) > 0 {
		os.Exit(1)
	}
	return nil
}
