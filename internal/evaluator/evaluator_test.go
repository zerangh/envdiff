package evaluator_test

import (
	"strings"
	"testing"

	"github.com/your/envdiff/internal/differ"
	"github.com/your/envdiff/internal/evaluator"
)

func makeResult(missingRight, missingLeft []string, mismatched []differ.Mismatch) differ.Result {
	return differ.Result{
		MissingInRight: missingRight,
		MissingInLeft:  missingLeft,
		Mismatched:     mismatched,
	}
}

func TestEvaluate_Clean(t *testing.T) {
	result := makeResult(nil, nil, nil)
	report := evaluator.Evaluate(result)
	if len(report.Findings) != 0 {
		t.Errorf("expected no findings, got %d", len(report.Findings))
	}
	if report.HasErrors || report.HasWarnings {
		t.Error("expected no errors or warnings")
	}
}

func TestEvaluate_MissingInRight_IsError(t *testing.T) {
	result := makeResult([]string{"DB_HOST"}, nil, nil)
	report := evaluator.Evaluate(result)
	if !report.HasErrors {
		t.Error("expected HasErrors to be true")
	}
	if len(report.Findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(report.Findings))
	}
	if report.Findings[0].Key != "DB_HOST" {
		t.Errorf("unexpected key: %s", report.Findings[0].Key)
	}
	if report.Findings[0].Rule.Severity != "error" {
		t.Errorf("expected error severity, got %s", report.Findings[0].Rule.Severity)
	}
}

func TestEvaluate_MissingInLeft_IsWarn(t *testing.T) {
	result := makeResult(nil, []string{"API_KEY"}, nil)
	report := evaluator.Evaluate(result)
	if !report.HasWarnings {
		t.Error("expected HasWarnings to be true")
	}
	if report.HasErrors {
		t.Error("expected no errors")
	}
}

func TestEvaluate_Mismatch_IsWarn(t *testing.T) {
	mm := []differ.Mismatch{{Key: "PORT", Left: "8080", Right: "9090"}}
	result := makeResult(nil, nil, mm)
	report := evaluator.Evaluate(result)
	if !report.HasWarnings {
		t.Error("expected HasWarnings to be true")
	}
	if len(report.Findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(report.Findings))
	}
	if !strings.Contains(report.Findings[0].Detail, "8080") {
		t.Errorf("expected detail to contain left value, got: %s", report.Findings[0].Detail)
	}
}

func TestFinding_String_WithDetail(t *testing.T) {
	f := evaluator.Finding{
		Rule:   evaluator.Rule{Name: "value-mismatch", Severity: "warn"},
		Key:    "PORT",
		Detail: "\"8080\" vs \"9090\"",
	}
	s := f.String()
	if !strings.Contains(s, "[WARN]") {
		t.Errorf("expected [WARN] in output, got: %s", s)
	}
	if !strings.Contains(s, "PORT") {
		t.Errorf("expected key in output, got: %s", s)
	}
}

func TestFinding_String_NoDetail(t *testing.T) {
	f := evaluator.Finding{
		Rule: evaluator.Rule{Name: "missing-in-right", Severity: "error"},
		Key:  "SECRET",
	}
	s := f.String()
	if !strings.Contains(s, "[ERROR]") {
		t.Errorf("expected [ERROR] in output, got: %s", s)
	}
}
