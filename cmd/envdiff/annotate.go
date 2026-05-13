package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/your/envdiff/internal/annotator"
	"github.com/your/envdiff/internal/differ"
	"github.com/your/envdiff/internal/parser"
)

func runAnnotate(args []string, out io.Writer) error {
	fs := flag.NewFlagSet("annotate", flag.ContinueOnError)
	format := fs.String("format", "text", "Output format: text or json")
	fs.SetOutput(out)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() < 2 {
		return fmt.Errorf("usage: envdiff annotate [--format=text|json] <left.env> <right.env>")
	}

	leftPath := fs.Arg(0)
	rightPath := fs.Arg(1)

	leftEnv, err := parser.ParseFile(leftPath)
	if err != nil {
		return fmt.Errorf("parsing %s: %w", leftPath, err)
	}

	rightEnv, err := parser.ParseFile(rightPath)
	if err != nil {
		return fmt.Errorf("parsing %s: %w", rightPath, err)
	}

	diffResult := differ.Diff(leftEnv, rightEnv)
	annotations := annotator.Annotate(diffResult)

	switch *format {
	case "json":
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		return enc.Encode(annotations)
	default:
		fmt.Fprintln(out, annotations.Format())
		return nil
	}
}

func init() {
	_ = os.Stderr // ensure os is used
}
