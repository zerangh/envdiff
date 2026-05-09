package inspector_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/inspector"
)

func TestInspect_Clean(t *testing.T) {
	env := map[string]string{
		"APP_NAME": "myapp",
		"APP_PORT": "8080",
	}
	keys := []string{"APP_NAME", "APP_PORT"}

	r := inspector.Inspect(env, keys)

	if r.TotalKeys != 2 {
		t.Errorf("expected 2 total keys, got %d", r.TotalKeys)
	}
	if len(r.EmptyValues) != 0 {
		t.Errorf("expected no empty values, got %v", r.EmptyValues)
	}
	if len(r.DuplicateKeys) != 0 {
		t.Errorf("expected no duplicates, got %v", r.DuplicateKeys)
	}
	if len(r.SpecialCharKeys) != 0 {
		t.Errorf("expected no special char keys, got %v", r.SpecialCharKeys)
	}
}

func TestInspect_EmptyValues(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "",
		"DB_PORT": "5432",
	}
	keys := []string{"DB_HOST", "DB_PORT"}

	r := inspector.Inspect(env, keys)

	if len(r.EmptyValues) != 1 || r.EmptyValues[0] != "DB_HOST" {
		t.Errorf("expected [DB_HOST] as empty value key, got %v", r.EmptyValues)
	}
}

func TestInspect_DuplicateKeys(t *testing.T) {
	env := map[string]string{
		"API_KEY": "abc123",
	}
	keys := []string{"API_KEY", "API_KEY", "OTHER"}

	r := inspector.Inspect(env, keys)

	if len(r.DuplicateKeys) != 1 || r.DuplicateKeys[0] != "API_KEY" {
		t.Errorf("expected [API_KEY] as duplicate, got %v", r.DuplicateKeys)
	}
}

func TestInspect_SpecialCharKeys(t *testing.T) {
	env := map[string]string{
		"VALID_KEY":   "value",
		"INVALID-KEY": "value",
	}
	keys := []string{"VALID_KEY", "INVALID-KEY"}

	r := inspector.Inspect(env, keys)

	if len(r.SpecialCharKeys) != 1 || r.SpecialCharKeys[0] != "INVALID-KEY" {
		t.Errorf("expected [INVALID-KEY] as special char key, got %v", r.SpecialCharKeys)
	}
}

func TestReport_Format_Clean(t *testing.T) {
	r := inspector.Report{TotalKeys: 3}
	out := r.Format()

	if !strings.Contains(out, "Total keys: 3") {
		t.Errorf("expected total keys in output, got: %s", out)
	}
	if !strings.Contains(out, "Empty values: none") {
		t.Errorf("expected 'Empty values: none' in output, got: %s", out)
	}
	if !strings.Contains(out, "Duplicate keys: none") {
		t.Errorf("expected 'Duplicate keys: none' in output, got: %s", out)
	}
}

func TestReport_Format_WithIssues(t *testing.T) {
	r := inspector.Report{
		TotalKeys:       2,
		EmptyValues:     []string{"SECRET"},
		DuplicateKeys:   []string{"HOST"},
		SpecialCharKeys: []string{"BAD-KEY"},
	}
	out := r.Format()

	if !strings.Contains(out, "SECRET") {
		t.Errorf("expected SECRET in output, got: %s", out)
	}
	if !strings.Contains(out, "HOST") {
		t.Errorf("expected HOST in output, got: %s", out)
	}
	if !strings.Contains(out, "BAD-KEY") {
		t.Errorf("expected BAD-KEY in output, got: %s", out)
	}
}
