package annotator_test

import (
	"strings"
	"testing"

	"github.com/your/envdiff/internal/annotator"
	"github.com/your/envdiff/internal/differ"
)

func makeResult(missingRight, missingLeft []string, mismatched []differ.Mismatch) differ.Result {
	return differ.Result{
		MissingInRight: missingRight,
		MissingInLeft:  missingLeft,
		Mismatched:     mismatched,
	}
}

func TestAnnotate_Clean(t *testing.T) {
	r := makeResult(nil, nil, nil)
	out := annotator.Annotate(r)
	if len(out.Annotations) != 0 {
		t.Errorf("expected 0 annotations, got %d", len(out.Annotations))
	}
	if out.Format() != "No annotations." {
		t.Errorf("unexpected format output: %s", out.Format())
	}
}

func TestAnnotate_MissingInRight(t *testing.T) {
	r := makeResult([]string{"DB_HOST"}, nil, nil)
	out := annotator.Annotate(r)
	if len(out.Annotations) != 1 {
		t.Fatalf("expected 1 annotation, got %d", len(out.Annotations))
	}
	a := out.Annotations[0]
	if a.Key != "DB_HOST" {
		t.Errorf("expected key DB_HOST, got %s", a.Key)
	}
	if a.Severity != annotator.SeverityWarning {
		t.Errorf("expected warning severity, got %s", a.Severity)
	}
}

func TestAnnotate_MissingInLeft(t *testing.T) {
	r := makeResult(nil, []string{"NEW_FEATURE_FLAG"}, nil)
	out := annotator.Annotate(r)
	if len(out.Annotations) != 1 {
		t.Fatalf("expected 1 annotation, got %d", len(out.Annotations))
	}
	if out.Annotations[0].Severity != annotator.SeverityInfo {
		t.Errorf("expected info severity")
	}
}

func TestAnnotate_MismatchedValues(t *testing.T) {
	mm := []differ.Mismatch{{Key: "API_URL", Left: "http://local", Right: "https://prod"}}
	r := makeResult(nil, nil, mm)
	out := annotator.Annotate(r)
	if len(out.Annotations) != 1 {
		t.Fatalf("expected 1 annotation, got %d", len(out.Annotations))
	}
	if out.Annotations[0].Severity != annotator.SeverityWarning {
		t.Errorf("expected warning severity")
	}
}

func TestAnnotate_EmptyValueMismatch_IsCritical(t *testing.T) {
	mm := []differ.Mismatch{{Key: "SECRET", Left: "", Right: "abc123"}}
	r := makeResult(nil, nil, mm)
	out := annotator.Annotate(r)
	if out.Annotations[0].Severity != annotator.SeverityCritical {
		t.Errorf("expected critical severity for empty-value mismatch")
	}
}

func TestAnnotate_Format_MultipleEntries(t *testing.T) {
	r := makeResult(
		[]string{"DB_HOST"},
		[]string{"NEW_KEY"},
		[]differ.Mismatch{{Key: "PORT", Left: "5432", Right: "3306"}},
	)
	out := annotator.Annotate(r)
	formatted := out.Format()
	if !strings.Contains(formatted, "DB_HOST") {
		t.Error("expected DB_HOST in output")
	}
	if !strings.Contains(formatted, "NEW_KEY") {
		t.Error("expected NEW_KEY in output")
	}
	if !strings.Contains(formatted, "PORT") {
		t.Error("expected PORT in output")
	}
}
