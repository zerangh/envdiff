package anchorer_test

import (
	"testing"

	"github.com/user/envdiff/internal/anchorer"
)

func TestAnchor_EmptyInput(t *testing.T) {
	r := anchorer.Anchor(nil, "")
	if len(r.AnchorKeys) != 0 || len(r.Deviations) != 0 {
		t.Fatal("expected empty result for nil input")
	}
}

func TestAnchor_SingleEnv(t *testing.T) {
	envs := map[string]map[string]string{
		"prod": {"A": "1", "B": "2"},
	}
	r := anchorer.Anchor(envs, "prod")
	if len(r.AnchorKeys) != 2 {
		t.Fatalf("expected 2 anchor keys, got %d", len(r.AnchorKeys))
	}
	if len(r.Deviations) != 0 {
		t.Fatalf("expected no deviations for single env, got %d", len(r.Deviations))
	}
}

func TestAnchor_MissingKeyInOtherEnv(t *testing.T) {
	envs := map[string]map[string]string{
		"prod":    {"A": "1", "B": "2", "C": "3"},
		"staging": {"A": "1"},
	}
	r := anchorer.Anchor(envs, "prod")
	if len(r.AnchorKeys) != 3 {
		t.Fatalf("expected 3 anchor keys, got %d", len(r.AnchorKeys))
	}
	if len(r.Deviations) != 2 {
		t.Fatalf("expected 2 deviations (B and C missing), got %d", len(r.Deviations))
	}
	if r.Deviations[0].Key != "B" || r.Deviations[1].Key != "C" {
		t.Errorf("unexpected deviation keys: %v", r.Deviations)
	}
}

func TestAnchor_ExtraKeyInOtherEnv(t *testing.T) {
	envs := map[string]map[string]string{
		"prod":    {"A": "1"},
		"staging": {"A": "1", "EXTRA": "x"},
	}
	r := anchorer.Anchor(envs, "prod")
	if len(r.Deviations) != 1 {
		t.Fatalf("expected 1 deviation, got %d", len(r.Deviations))
	}
	d := r.Deviations[0]
	if d.Key != "EXTRA" || d.Env != "staging" {
		t.Errorf("unexpected deviation: %+v", d)
	}
}

func TestAnchor_AutoPicksLargestEnv(t *testing.T) {
	envs := map[string]map[string]string{
		"small": {"A": "1"},
		"large": {"A": "1", "B": "2", "C": "3"},
	}
	r := anchorer.Anchor(envs, "")
	if len(r.AnchorKeys) != 3 {
		t.Fatalf("expected large env chosen as anchor (3 keys), got %d", len(r.AnchorKeys))
	}
	if len(r.Deviations) != 2 {
		t.Fatalf("expected 2 deviations for small env, got %d", len(r.Deviations))
	}
}

func TestAnchor_UnknownAnchorEnv(t *testing.T) {
	envs := map[string]map[string]string{
		"prod": {"A": "1"},
	}
	r := anchorer.Anchor(envs, "nonexistent")
	if len(r.AnchorKeys) != 0 {
		t.Fatal("expected empty result for unknown anchor env")
	}
}

func TestAnchor_DeviationsSorted(t *testing.T) {
	envs := map[string]map[string]string{
		"prod":    {"A": "1", "B": "2", "C": "3"},
		"dev":     {},
	}
	r := anchorer.Anchor(envs, "prod")
	for i := 1; i < len(r.Deviations); i++ {
		if r.Deviations[i].Key < r.Deviations[i-1].Key {
			t.Errorf("deviations not sorted at index %d: %v", i, r.Deviations)
		}
	}
}
