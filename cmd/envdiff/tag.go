package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/your-org/envdiff/internal/differ"
	"github.com/your-org/envdiff/internal/parser"
	"github.com/your-org/envdiff/internal/tagger"
)

func runTag(args []string, format string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envdiff tag <left> <right> [--format text|json]")
	}

	left, err := parser.ParseFile(args[0])
	if err != nil {
		return fmt.Errorf("parsing %s: %w", args[0], err)
	}
	right, err := parser.ParseFile(args[1])
	if err != nil {
		return fmt.Errorf("parsing %s: %w", args[1], err)
	}

	result := differ.Diff(left, right)
	tagResult := tagger.TagResult(result)

	switch format {
	case "json":
		return printTagJSON(tagResult)
	default:
		return printTagText(tagResult)
	}
}

func printTagText(r tagger.Result) error {
	if len(r.Tags) == 0 {
		fmt.Println("No keys to tag.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "KEY\tTAGS")
	fmt.Fprintln(w, "---\t----")

	keys := make([]string, 0, len(r.Tags))
	for k := range r.Tags {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		tagStrs := make([]string, 0, len(r.Tags[k]))
		for _, t := range r.Tags[k] {
			tagStrs = append(tagStrs, string(t))
		}
		fmt.Fprintf(w, "%s\t%v\n", k, tagStrs)
	}
	return w.Flush()
}

func printTagJSON(r tagger.Result) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
