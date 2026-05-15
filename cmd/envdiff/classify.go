package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/classifier"
	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/parser"
)

func runClassify(args []string, format string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envdiff classify <file1> <file2>")
	}

	left, err := parser.ParseFile(args[0])
	if err != nil {
		return fmt.Errorf("reading %s: %w", args[0], err)
	}
	right, err := parser.ParseFile(args[1])
	if err != nil {
		return fmt.Errorf("reading %s: %w", args[1], err)
	}

	diffResult := differ.Diff(left, right)
	classResult := classifier.Classify(diffResult)

	if len(classResult.Classes) == 0 {
		fmt.Println("No differences to classify.")
		return nil
	}

	switch format {
	case "json":
		return printClassifyJSON(classResult)
	default:
		printClassifyText(classResult)
		return nil
	}
}

func printClassifyText(res classifier.Result) {
	current := classifier.Category("")
	for _, c := range res.Classes {
		if c.Category != current {
			current = c.Category
			fmt.Fprintf(os.Stdout, "\n[%s]\n", current)
		}
		fmt.Fprintf(os.Stdout, "  %s\n", c.Key)
	}
}

func printClassifyJSON(res classifier.Result) error {
	type jsonOut struct {
		Classes    []classifier.KeyClass       `json:"classes"`
		Categories map[classifier.Category][]string `json:"categories"`
	}
	out := jsonOut{Classes: res.Classes, Categories: res.Categories}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
