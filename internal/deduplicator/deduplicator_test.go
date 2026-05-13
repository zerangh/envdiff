package deduplicator_test

import (
	"testing"

	"github.com/user/envdiff/internal/deduplicator"
	"github.com/user/envdiff/internal/differ"
)

func TestDetect_NoDuplicates(t *testing.T) {
	envs := []map[string]string{
		{"FOO": "bar", "BAZ": "qux"},
		{"FOO": "bar", "NEW": "val"},
	}
	r := deduplicator.Detect(envs)
	if len(r.Duplicates) != 0 {
		t.Fatalf("expected no duplicates, got %d", len(r.Duplicates))
	}
}

func TestDetect_FindsDuplicates(t *testing.T) {
	envs := []map[string]string{
		{"FOO": "original"},
		{"FOO": "changed"},
	}
	r := deduplicator.Detect(envs)
	if len(r.Duplicates) != 1 {
		t.Fatalf("expected 1 duplicate, got %d", len(r.Duplicates))
	}
	if r.Duplicates[0].Key != "FOO" {
		t.Errorf("expected key FOO, got %s", r.Duplicates[0].Key)
	}
	if len(r.Duplicates[0].Values) != 2 {
		t.Errorf("expected 2 values, got %d", len(r.Duplicates[0].Values))
	}
}

func TestDetect_SortedOutput(t *testing.T) {
	envs := []map[string]string{
		{"ZEBRA": "a", "ALPHA": "x"},
		{"ZEBRA": "b", "ALPHA": "y"},
	}
	r := deduplicator.Detect(envs)
	if len(r.Duplicates) != 2 {
		t.Fatalf("expected 2 duplicates, got %d", len(r.Duplicates))
	}
	if r.Duplicates[0].Key != "ALPHA" {
		t.Errorf("expected sorted first key ALPHA, got %s", r.Duplicates[0].Key)
	}
}

func TestResolve_KeepFirst(t *testing.T) {
	envs := []map[string]string{
		{"FOO": "first"},
		{"FOO": "second"},
	}
	out := deduplicator.Resolve(envs, deduplicator.StrategyKeepFirst)
	if out["FOO"] != "first" {
		t.Errorf("expected 'first', got %s", out["FOO"])
	}
}

func TestResolve_KeepLast(t *testing.T) {
	envs := []map[string]string{
		{"FOO": "first"},
		{"FOO": "second"},
	}
	out := deduplicator.Resolve(envs, deduplicator.StrategyKeepLast)
	if out["FOO"] != "second" {
		t.Errorf("expected 'second', got %s", out["FOO"])
	}
}

func TestFromResult_MismatchedBecomeDuplicates(t *testing.T) {
	result := differ.Result{
		Mismatched: []differ.Mismatch{
			{Key: "DB_URL", Left: "postgres://dev", Right: "postgres://prod"},
			{Key: "API_KEY", Left: "abc", Right: "xyz"},
		},
	}
	dups := deduplicator.FromResult(result)
	if len(dups) != 2 {
		t.Fatalf("expected 2 duplicates, got %d", len(dups))
	}
	if dups[0].Key != "API_KEY" {
		t.Errorf("expected sorted first key API_KEY, got %s", dups[0].Key)
	}
}

func TestFromResult_NoMismatches(t *testing.T) {
	result := differ.Result{}
	dups := deduplicator.FromResult(result)
	if len(dups) != 0 {
		t.Errorf("expected no duplicates, got %d", len(dups))
	}
}
