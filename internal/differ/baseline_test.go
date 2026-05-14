package differ

import (
	"testing"
)

func TestCompareToBaseline_NoChanges(t *testing.T) {
	base := map[string]string{"FOO": "bar", "BAZ": "qux"}
	curr := map[string]string{"FOO": "bar", "BAZ": "qux"}

	r := CompareToBaseline(base, curr)

	if len(r.Added) != 0 {
		t.Errorf("expected no added, got %v", r.Added)
	}
	if len(r.Removed) != 0 {
		t.Errorf("expected no removed, got %v", r.Removed)
	}
	if len(r.Changed) != 0 {
		t.Errorf("expected no changed, got %v", r.Changed)
	}
	if r.Unchanged != 2 {
		t.Errorf("expected unchanged=2, got %d", r.Unchanged)
	}
}

func TestCompareToBaseline_AddedKeys(t *testing.T) {
	base := map[string]string{"FOO": "bar"}
	curr := map[string]string{"FOO": "bar", "NEW_KEY": "value"}

	r := CompareToBaseline(base, curr)

	if len(r.Added) != 1 || r.Added[0] != "NEW_KEY" {
		t.Errorf("expected Added=[NEW_KEY], got %v", r.Added)
	}
	if len(r.Removed) != 0 {
		t.Errorf("expected no removed, got %v", r.Removed)
	}
}

func TestCompareToBaseline_RemovedKeys(t *testing.T) {
	base := map[string]string{"FOO": "bar", "OLD_KEY": "gone"}
	curr := map[string]string{"FOO": "bar"}

	r := CompareToBaseline(base, curr)

	if len(r.Removed) != 1 || r.Removed[0] != "OLD_KEY" {
		t.Errorf("expected Removed=[OLD_KEY], got %v", r.Removed)
	}
	if len(r.Added) != 0 {
		t.Errorf("expected no added, got %v", r.Added)
	}
}

func TestCompareToBaseline_ChangedValues(t *testing.T) {
	base := map[string]string{"FOO": "old", "BAR": "same"}
	curr := map[string]string{"FOO": "new", "BAR": "same"}

	r := CompareToBaseline(base, curr)

	if len(r.Changed) != 1 {
		t.Fatalf("expected 1 changed, got %d", len(r.Changed))
	}
	if r.Changed[0].Key != "FOO" {
		t.Errorf("expected changed key FOO, got %s", r.Changed[0].Key)
	}
	if r.Changed[0].LeftValue != "old" || r.Changed[0].RightValue != "new" {
		t.Errorf("unexpected changed values: %+v", r.Changed[0])
	}
	if r.Unchanged != 1 {
		t.Errorf("expected unchanged=1, got %d", r.Unchanged)
	}
}

func TestCompareToBaseline_SortedOutput(t *testing.T) {
	base := map[string]string{}
	curr := map[string]string{"ZEBRA": "z", "ALPHA": "a", "MIDDLE": "m"}

	r := CompareToBaseline(base, curr)

	if len(r.Added) != 3 {
		t.Fatalf("expected 3 added, got %d", len(r.Added))
	}
	if r.Added[0] != "ALPHA" || r.Added[1] != "MIDDLE" || r.Added[2] != "ZEBRA" {
		t.Errorf("expected sorted added keys, got %v", r.Added)
	}
}
