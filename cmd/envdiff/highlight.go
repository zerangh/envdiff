package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/highlighter"
	"github.com/user/envdiff/internal/parser"
)

func runHighlight(args []string, format string) error {
	if len(args) < 2 {
		return fmt.Errorf("highlight requires two .env files")
	}

	left, err := parser.ParseFile(args[0])
	if err != nil {
		return fmt.Errorf("reading %s: %w", args[0], err)
	}
	right, err := parser.ParseFile(args[1])
	if err != nil {
		return fmt.Errorf("reading %s: %w", args[1], err)
	}

	result := differ.Diff(left, right)
	hr := highlighter.Compute(result)

	if len(hr.Highlights) == 0 {
		fmt.Fprintln(os.Stdout, "No mismatched values found.")
		return nil
	}

	switch format {
	case "json":
		return printHighlightJSON(hr)
	default:
		printHighlightText(hr)
		return nil
	}
}

func printHighlightText(hr highlighter.HighlightResult) {
	for _, h := range hr.Highlights {
		fmt.Print(highlighter.Format(h))
	}
}

func printHighlightJSON(hr highlighter.HighlightResult) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(hr)
}

func init() {
	registerCommand("highlight", func(args []string, flags map[string]string) error {
		format := flags["format"]
		return runHighlight(args, format)
	})
}
