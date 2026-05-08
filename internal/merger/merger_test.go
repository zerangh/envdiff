package merger_test

import (
	"testing"

	"github.com/user/envdiff/internal/merger"
)

func nm(path string, env map[string]string) merger.NamedMap {
	return merger.NamedMap{Path: path, Env: env}
}

func TestMerge_EmptyInput(t *testing.T) {
	_, err := merger.Merge(nil, merger.StrategyFirst)
	if err == nil {
		t.Fatal("expected error for empty input, got nil")
	}
}

func TestMerge_SingleFile(t *testing.T) {
	files := []merger.NamedMap{
		nm(".env", map[string]string{"FOO": "bar", "BAZ": "qux"}),
	}
	res, err := merger.Merge(files, merger.StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Merged["FOO"] != "bar" || res.Merged["BAZ"] != "qux" {
		t.Errorf("unexpected merged values: %v", res.Merged)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got: %v", res.Conflicts)
	}
}

func TestMerge_StrategyFirst_KeepsFirst(t *testing.T) {
	files := []merger.NamedMap{
		nm(".env.base", map[string]string{"KEY": "first"}),
		nm(".env.local", map[string]string{"KEY": "second"}),
	}
	res, err := merger.Merge(files, merger.StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Merged["KEY"] != "first" {
		t.Errorf("expected 'first', got %q", res.Merged["KEY"])
	}
	if len(res.Conflicts["KEY"]) != 2 {
		t.Errorf("expected 2 conflict entries for KEY, got %v", res.Conflicts["KEY"])
	}
}

func TestMerge_StrategyLast_KeepsLast(t *testing.T) {
	files := []merger.NamedMap{
		nm(".env.base", map[string]string{"KEY": "first"}),
		nm(".env.local", map[string]string{"KEY": "second"}),
	}
	res, err := merger.Merge(files, merger.StrategyLast)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Merged["KEY"] != "second" {
		t.Errorf("expected 'second', got %q", res.Merged["KEY"])
	}
	if res.Sources["KEY"] != ".env.local" {
		t.Errorf("expected source '.env.local', got %q", res.Sources["KEY"])
	}
}

func TestMerge_NoConflict_UniqueKeys(t *testing.T) {
	files := []merger.NamedMap{
		nm(".env.a", map[string]string{"A": "1"}),
		nm(".env.b", map[string]string{"B": "2"}),
	}
	res, err := merger.Merge(files, merger.StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got: %v", res.Conflicts)
	}
	if res.Merged["A"] != "1" || res.Merged["B"] != "2" {
		t.Errorf("unexpected merged values: %v", res.Merged)
	}
}
