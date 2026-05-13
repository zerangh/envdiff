package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/profiler"
)

// runProfile parses two or more .env files supplied as positional arguments
// and prints a cross-environment key profile to stdout.
//
// Usage: envdiff profile [--json] file1.env file2.env [file3.env ...]
func runProfile(args []string) error {
	useJSON := false
	files := []string{}
	for _, a := range args {
		if a == "--json" {
			useJSON = true
		} else {
			files = append(files, a)
		}
	}

	if len(files) < 2 {
		return fmt.Errorf("profile requires at least 2 .env files")
	}

	envs := make(map[string]map[string]string, len(files))
	for _, f := range files {
		m, err := parser.ParseFile(f)
		if err != nil {
			return fmt.Errorf("parsing %s: %w", f, err)
		}
		envs[f] = m
	}

	report := profiler.Profile(envs)

	if useJSON {
		return printProfileJSON(report)
	}
	return printProfileText(report)
}

func printProfileText(r profiler.Report) error {
	fmt.Fprintf(os.Stdout, "Environments : %s\n", strings.Join(r.EnvNames, ", "))
	fmt.Fprintf(os.Stdout, "Total keys   : %d\n", r.TotalKeys)
	fmt.Fprintf(os.Stdout, "Consistent   : %d\n", r.Consistent)
	fmt.Fprintf(os.Stdout, "Inconsistent : %d\n\n", r.Inconsistent)

	for _, p := range r.Profiles {
		status := "OK"
		if len(p.MissingFrom) > 0 || p.UniqueValues > 1 {
			status = "WARN"
		}
		fmt.Fprintf(os.Stdout, "[%s] %s", status, p.Key)
		if len(p.MissingFrom) > 0 {
			fmt.Fprintf(os.Stdout, " (missing: %s)", strings.Join(p.MissingFrom, ", "))
		}
		if p.UniqueValues > 1 {
			fmt.Fprintf(os.Stdout, " (%d unique values)", p.UniqueValues)
		}
		fmt.Fprintln(os.Stdout)
	}
	return nil
}

func printProfileJSON(r profiler.Report) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
