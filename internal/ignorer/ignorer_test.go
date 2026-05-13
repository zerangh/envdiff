package ignorer_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/ignorer"
)

func writeIgnoreFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".envdiffignore")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write ignore file: %v", err)
	}
	return p
}

func baseResult() differ.Result {
	return differ.Result{
		MissingInRight: []string{"SECRET", "HOST"},
		MissingInLeft:  []string{"DEBUG"},
		Mismatched: []differ.Mismatch{
			{Key: "PORT", Left: "8080", Right: "9090"},
			{Key: "DB_URL", Left: "a", Right: "b"},
		},
	}
}

func TestLoadFile_Basic(t *testing.T) {
	p := writeIgnoreFile(t, "# comment\nSECRET\nPORT\n")
	r, err := ignorer.LoadFile(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !r.Has("SECRET") {
		t.Error("expected SECRET to be ignored")
	}
	if !r.Has("PORT") {
		t.Error("expected PORT to be ignored")
	}
	if r.Has("HOST") {
		t.Error("HOST should not be ignored")
	}
}

func TestLoadFile_Missing(t *testing.T) {
	_, err := ignorer.LoadFile("/nonexistent/.envdiffignore")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestNewRules_Has(t *testing.T) {
	r := ignorer.NewRules([]string{"FOO", "BAR"})
	if !r.Has("FOO") || !r.Has("BAR") {
		t.Error("expected FOO and BAR to be present")
	}
	if r.Has("BAZ") {
		t.Error("BAZ should not be present")
	}
}

func TestApply_RemovesMissingInRight(t *testing.T) {
	r := ignorer.NewRules([]string{"SECRET"})
	out := r.Apply(baseResult())
	for _, k := range out.MissingInRight {
		if k == "SECRET" {
			t.Error("SECRET should have been removed from MissingInRight")
		}
	}
}

func TestApply_RemovesMismatched(t *testing.T) {
	r := ignorer.NewRules([]string{"PORT"})
	out := r.Apply(baseResult())
	for _, m := range out.Mismatched {
		if m.Key == "PORT" {
			t.Error("PORT should have been removed from Mismatched")
		}
	}
	if len(out.Mismatched) != 1 {
		t.Errorf("expected 1 mismatch, got %d", len(out.Mismatched))
	}
}

func TestApply_NoRules(t *testing.T) {
	r := ignorer.NewRules(nil)
	out := r.Apply(baseResult())
	if len(out.MissingInRight) != 2 || len(out.MissingInLeft) != 1 || len(out.Mismatched) != 2 {
		t.Error("expected result to be unchanged with no ignore rules")
	}
}
