package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/redactor"
)

// runRedact prints a redacted view of a single .env file to stdout.
// Usage: envdiff redact <file> [--mask=<mask>] [--format=text|json]
func runRedact(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envdiff redact <file> [--mask=<mask>] [--format=text|json]")
	}

	filePath := args[0]
	mask := "***"
	format := "text"

	for _, arg := range args[1:] {
		switch {
		case len(arg) > 7 && arg[:7] == "--mask=":
			mask = arg[7:]
		case len(arg) > 9 && arg[:9] == "--format=":
			format = arg[9:]
		}
	}

	env, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("parse %s: %w", filePath, err)
	}

	opts := redactor.Options{Mask: mask}
	safe := redactor.RedactEnv(env, opts)

	switch format {
	case "json":
		return printRedactJSON(safe)
	default:
		printRedactText(safe)
		return nil
	}
}

func printRedactText(env map[string]string) {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(os.Stdout, "%s=%s\n", k, env[k])
	}
}

func printRedactJSON(env map[string]string) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(env); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}
