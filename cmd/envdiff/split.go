package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/your-org/envdiff/internal/differ"
	"github.com/your-org/envdiff/internal/parser"
	"github.com/your-org/envdiff/internal/splitter"
)

func runSplit(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envdiff split <left> <right> [--strategy alpha|prefix] [--format text|json]")
	}

	leftFile := args[0]
	rightFile := args[1]

	strategy := splitter.StrategyPrefix
	format := "text"

	for i := 2; i < len(args)-1; i++ {
		switch args[i] {
		case "--strategy":
			strategy = splitter.Strategy(args[i+1])
			i++
		case "--format":
			format = args[i+1]
			i++
		}
	}

	leftEnv, err := parser.ParseFile(leftFile)
	if err != nil {
		return fmt.Errorf("reading %s: %w", leftFile, err)
	}
	rightEnv, err := parser.ParseFile(rightFile)
	if err != nil {
		return fmt.Errorf("reading %s: %w", rightFile, err)
	}

	result := differ.Diff(leftEnv, rightEnv)
	opts := splitter.Options{Strategy: strategy}
	buckets := splitter.Split(result, opts)

	switch format {
	case "json":
		return printSplitJSON(buckets)
	default:
		printSplitText(buckets, leftFile, rightFile)
		return nil
	}
}

func printSplitText(buckets []splitter.Bucket, left, right string) {
	fmt.Printf("Split diff: %s vs %s\n", left, right)
	if len(buckets) == 0 {
		fmt.Println("  No differences found.")
		return
	}
	for _, b := range buckets {
		fmt.Printf("\n[%s] (%d key(s))\n", b.Name, len(b.Keys))
		for _, k := range b.Result.MissingInRight {
			fmt.Printf("  - missing in right: %s\n", k)
		}
		for _, k := range b.Result.MissingInLeft {
			fmt.Printf("  - missing in left:  %s\n", k)
		}
		for _, m := range b.Result.Mismatched {
			fmt.Printf("  ~ mismatch: %s (%q vs %q)\n", m.Key, m.Left, m.Right)
		}
	}
}

func printSplitJSON(buckets []splitter.Bucket) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(buckets)
}
