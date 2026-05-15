package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/your/envdiff/internal/differ"
	"github.com/your/envdiff/internal/evaluator"
	"github.com/your/envdiff/internal/parser"
)

func runEvaluate(args []string, out io.Writer) error {
	fs := flag.NewFlagSet("evaluate", flag.ContinueOnError)
	format := fs.String("format", "text", "output format: text or json")
	fs.SetOutput(out)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() < 2 {
		return fmt.Errorf("usage: envdiff evaluate [--format=text|json] <left.env> <right.env>")
	}

	leftPath := fs.Arg(0)
	rightPath := fs.Arg(1)

	left, err := parser.ParseFile(leftPath)
	if err != nil {
		return fmt.Errorf("reading %s: %w", leftPath, err)
	}

	right, err := parser.ParseFile(rightPath)
	if err != nil {
		return fmt.Errorf("reading %s: %w", rightPath, err)
	}

	result := differ.Diff(left, right)
	report := evaluator.Evaluate(result)

	switch *format {
	case "json":
		return printEvaluateJSON(out, report)
	default:
		return printEvaluateText(out, report)
	}
}

func printEvaluateText(out io.Writer, report evaluator.Report) error {
	if len(report.Findings) == 0 {
		fmt.Fprintln(out, "OK: no issues found")
		return nil
	}
	for _, f := range report.Findings {
		fmt.Fprintln(out, f.String())
	}
	if report.HasErrors {
		return fmt.Errorf("evaluation failed: errors found")
	}
	return nil
}

func printEvaluateJSON(out io.Writer, report evaluator.Report) error {
	type jsonFinding struct {
		Severity string `json:"severity"`
		Rule     string `json:"rule"`
		Key      string `json:"key"`
		Detail   string `json:"detail,omitempty"`
	}
	type jsonReport struct {
		Findings    []jsonFinding `json:"findings"`
		HasErrors   bool          `json:"has_errors"`
		HasWarnings bool          `json:"has_warnings"`
	}
	jr := jsonReport{HasErrors: report.HasErrors, HasWarnings: report.HasWarnings}
	for _, f := range report.Findings {
		jr.Findings = append(jr.Findings, jsonFinding{
			Severity: f.Rule.Severity,
			Rule:     f.Rule.Name,
			Key:      f.Key,
			Detail:   f.Detail,
		})
	}
	if jr.Findings == nil {
		jr.Findings = []jsonFinding{}
	}
	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	if err := enc.Encode(jr); err != nil {
		return err
	}
	if report.HasErrors {
		return fmt.Errorf("evaluation failed: errors found")
	}
	return nil
}

func init() {
	_ = os.Stderr
}
