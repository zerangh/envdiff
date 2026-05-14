package formatter_test

import (
	"testing"

	"github.com/your-org/envdiff/internal/differ"
	"github.com/your-org/envdiff/internal/formatter"
)

func TestFormatEnv_PlainSorted(t *testing.T) {
	env := map[string]string{"ZEBRA": "1", "ALPHA": "2", "MIDDLE": "3"}
	opts := formatter.DefaultOptions()
	lines := formatter.FormatEnv(env, opts)
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if lines[0] != "ALPHA=2" {
		t.Errorf("expected ALPHA=2, got %s", lines[0])
	}
	if lines[2] != "ZEBRA=1" {
		t.Errorf("expected ZEBRA=1, got %s", lines[2])
	}
}

func TestFormatEnv_ExportStyle(t *testing.T) {
	env := map[string]string{"FOO": "bar"}
	opts := formatter.DefaultOptions()
	opts.Style = formatter.StyleExport
	lines := formatter.FormatEnv(env, opts)
	if lines[0] != "export FOO=bar" {
		t.Errorf("unexpected line: %s", lines[0])
	}
}

func TestFormatEnv_OmitEmpty(t *testing.T) {
	env := map[string]string{"PRESENT": "yes", "EMPTY": ""}
	opts := formatter.DefaultOptions()
	opts.OmitEmpty = true
	lines := formatter.FormatEnv(env, opts)
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
	if lines[0] != "PRESENT=yes" {
		t.Errorf("unexpected line: %s", lines[0])
	}
}

func TestFormatEnv_QuoteValues(t *testing.T) {
	env := map[string]string{"MSG": "hello world"}
	opts := formatter.DefaultOptions()
	opts.QuoteValues = true
	lines := formatter.FormatEnv(env, opts)
	if lines[0] != `MSG="hello world"` {
		t.Errorf("unexpected line: %s", lines[0])
	}
}

func TestFormatEnv_NoQuoteWhenNotNeeded(t *testing.T) {
	env := map[string]string{"KEY": "simple"}
	opts := formatter.DefaultOptions()
	opts.QuoteValues = true
	lines := formatter.FormatEnv(env, opts)
	if lines[0] != "KEY=simple" {
		t.Errorf("unexpected line: %s", lines[0])
	}
}

func TestFormatResult_UsesMismatchedLeftValues(t *testing.T) {
	result := differ.Result{
		Mismatched: []differ.Mismatch{
			{Key: "DB_HOST", LeftValue: "localhost", RightValue: "prod.db"},
		},
		MissingInRight: []string{"ORPHAN"},
	}
	opts := formatter.DefaultOptions()
	lines := formatter.FormatResult(result, opts)
	found := map[string]bool{}
	for _, l := range lines {
		found[l] = true
	}
	if !found["DB_HOST=localhost"] {
		t.Error("expected DB_HOST=localhost in output")
	}
	if !found["ORPHAN="] {
		t.Error("expected ORPHAN= in output")
	}
}
