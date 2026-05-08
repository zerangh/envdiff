package linter

import (
	"testing"

	"github.com/user/envdiff/internal/differ"
)

func emptyDiff() differ.Result {
	return differ.Result{}
}

func TestLint_NoIssues(t *testing.T) {
	env := map[string]string{
		"APP_NAME": "myapp",
		"PORT":     "8080",
	}
	result := Lint(env, emptyDiff())
	if len(result.Issues) != 0 {
		t.Errorf("expected no issues, got %d", len(result.Issues))
	}
}

func TestLint_LowercaseKey(t *testing.T) {
	env := map[string]string{
		"app_name": "myapp",
	}
	result := Lint(env, emptyDiff())
	if len(result.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(result.Issues))
	}
	if result.Issues[0].Severity != "warning" {
		t.Errorf("expected warning, got %s", result.Issues[0].Severity)
	}
}

func TestLint_WhitespaceValue(t *testing.T) {
	env := map[string]string{
		"API_KEY": "  secret  ",
	}
	result := Lint(env, emptyDiff())
	if len(result.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(result.Issues))
	}
	if result.Issues[0].Severity != "warning" {
		t.Errorf("expected warning, got %s", result.Issues[0].Severity)
	}
}

func TestLint_NewlineInValue(t *testing.T) {
	env := map[string]string{
		"SECRET": "line1\nline2",
	}
	result := Lint(env, emptyDiff())
	if len(result.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(result.Issues))
	}
	if result.Issues[0].Severity != "error" {
		t.Errorf("expected error, got %s", result.Issues[0].Severity)
	}
	if !result.HasErrors() {
		t.Error("expected HasErrors() to return true")
	}
}

func TestLint_MissingInRight(t *testing.T) {
	env := map[string]string{}
	diff := differ.Result{
		MissingInRight: []string{"DB_HOST", "DB_PORT"},
	}
	result := Lint(env, diff)
	if len(result.Issues) != 2 {
		t.Fatalf("expected 2 issues, got %d", len(result.Issues))
	}
	for _, issue := range result.Issues {
		if issue.Severity != "warning" {
			t.Errorf("expected warning for missing key, got %s", issue.Severity)
		}
	}
}

func TestLint_HasErrors_False(t *testing.T) {
	result := Result{
		Issues: []Issue{
			{Key: "X", Severity: "warning", Message: "some warning"},
		},
	}
	if result.HasErrors() {
		t.Error("expected HasErrors() to return false for warnings only")
	}
}
