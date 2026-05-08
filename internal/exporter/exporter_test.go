package exporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/exporter"
)

func makeResult() differ.Result {
	return differ.Result{
		MissingInRight: []string{"DB_HOST"},
		MissingInLeft:  []string{"NEW_KEY"},
		Mismatched: []differ.Mismatch{
			{Key: "APP_ENV", LeftValue: "production", RightValue: "staging"},
		},
	}
}

func TestExport_UnknownFormat(t *testing.T) {
	var buf bytes.Buffer
	err := exporter.Export(&buf, differ.Result{}, "a", "b", "xml")
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestExport_EnvFormat(t *testing.T) {
	var buf bytes.Buffer
	err := exporter.Export(&buf, makeResult(), "a.env", "b.env", exporter.FormatEnv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "APP_ENV=production") {
		t.Errorf("expected APP_ENV=production in env output, got:\n%s", out)
	}
	// missing keys should not appear
	if strings.Contains(out, "DB_HOST") {
		t.Errorf("missing key DB_HOST should not appear in env output")
	}
}

func TestExport_MarkdownFormat_WithDiffs(t *testing.T) {
	var buf bytes.Buffer
	err := exporter.Export(&buf, makeResult(), "prod.env", "staging.env", exporter.FormatMarkdown)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"prod.env", "staging.env", "DB_HOST", "NEW_KEY", "APP_ENV", "production", "staging"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in markdown output, got:\n%s", want, out)
		}
	}
}

func TestExport_MarkdownFormat_Clean(t *testing.T) {
	var buf bytes.Buffer
	err := exporter.Export(&buf, differ.Result{}, "a.env", "b.env", exporter.FormatMarkdown)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No differences found") {
		t.Errorf("expected clean message, got: %s", buf.String())
	}
}

func TestExport_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	err := exporter.Export(&buf, makeResult(), "prod.env", "staging.env", exporter.FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"prod.env", "staging.env", "DB_HOST", "APP_ENV"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in JSON output, got:\n%s", want, out)
		}
	}
}
