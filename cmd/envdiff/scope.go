package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/your-org/envdiff/internal/differ"
	"github.com/your-org/envdiff/internal/parser"
	"github.com/your-org/envdiff/internal/scoper"
)

// runScope implements the `envdiff scope` sub-command.
// It parses two .env files, diffs them, and prints per-scope summaries.
func runScope(args []string, format string) error {
	if len(args) < 2 {
		return fmt.Errorf("scope requires two .env file paths")
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
	scopes := scoper.Extract(result, scoper.DefaultOptions())

	if format == "json" {
		return printScopeJSON(scopes)
	}
	return printScopeText(scopes)
}

func printScopeText(scopes []scoper.Scope) error {
	if len(scopes) == 0 {
		fmt.Println("No differences found across any scope.")
		return nil
	}
	for _, s := range scopes {
		missingR := len(s.Result.MissingInRight)
		missingL := len(s.Result.MissingInLeft)
		mismatch := len(s.Result.Mismatched)
		fmt.Printf("[%s] keys=%d  missing_right=%d  missing_left=%d  mismatched=%d\n",
			s.Name, len(s.Keys), missingR, missingL, mismatch)
	}
	return nil
}

func printScopeJSON(scopes []scoper.Scope) error {
	type scopeJSON struct {
		Name         string   `json:"name"`
		Keys         []string `json:"keys"`
		MissingRight int      `json:"missing_in_right"`
		MissingLeft  int      `json:"missing_in_left"`
		Mismatched   int      `json:"mismatched"`
	}
	out := make([]scopeJSON, 0, len(scopes))
	for _, s := range scopes {
		out = append(out, scopeJSON{
			Name:         s.Name,
			Keys:         s.Keys,
			MissingRight: len(s.Result.MissingInRight),
			MissingLeft:  len(s.Result.MissingInLeft),
			Mismatched:   len(s.Result.Mismatched),
		})
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
