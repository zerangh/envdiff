package reporter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/reporter"
)

func makeResult(missingLeft, missingRight []string, mismatched []differ.Mismatch) differ.Result {
	return differ.Result{
		MissingInLeft:  missingLeft,
		MissingInRight: missingRight,
		Mismatched:     mismatched,
	}
}

func TestReportText_Clean(t *testing.T) {
	var buf bytes.Buffer
	result := makeResult(nil, nil, nil)
	if err := reporter.Report(&buf, result, "left", "right", reporter.FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected clean message, got: %q", buf.String())
	}
}

func TestReportText_MissingInRight(t *testing.T) {
	var buf bytes.Buffer
	result := makeResult(nil, []string{"FOO", "BAR"}, nil)
	_ = reporter.Report(&buf, result, "left", "right", reporter.FormatText)
	out := buf.String()
	if !strings.Contains(out, "- FOO") || !strings.Contains(out, "- BAR") {
		t.Errorf("expected missing keys in output, got: %q", out)
	}
}

func TestReportText_Mismatched(t *testing.T) {
	var buf bytes.Buffer
	result := makeResult(nil, nil, []differ.Mismatch{
		{Key: "DB_HOST", LeftValue: "localhost", RightValue: "prod.db"},
	})
	_ = reporter.Report(&buf, result, "dev", "prod", reporter.FormatText)
	out := buf.String()
	if !strings.Contains(out, "~ DB_HOST") {
		t.Errorf("expected mismatch marker, got: %q", out)
	}
	if !strings.Contains(out, "localhost") || !strings.Contains(out, "prod.db") {
		t.Errorf("expected both values in output, got: %q", out)
	}
}

func TestReportJSON_Clean(t *testing.T) {
	var buf bytes.Buffer
	result := makeResult(nil, nil, nil)
	if err := reporter.Report(&buf, result, "left", "right", reporter.FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if clean, ok := out["clean"].(bool); !ok || !clean {
		t.Errorf("expected clean=true in JSON output")
	}
}

func TestReportJSON_Mismatched(t *testing.T) {
	var buf bytes.Buffer
	result := makeResult([]string{"ONLY_RIGHT"}, nil, []differ.Mismatch{
		{Key: "PORT", LeftValue: "3000", RightValue: "8080"},
	})
	_ = reporter.Report(&buf, result, "a", "b", reporter.FormatJSON)
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	mismatched, ok := out["mismatched"].([]interface{})
	if !ok || len(mismatched) != 1 {
		t.Errorf("expected 1 mismatch in JSON, got: %v", out["mismatched"])
	}
}
