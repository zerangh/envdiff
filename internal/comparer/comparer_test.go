package comparer_test

import (
	"testing"

	"github.com/user/envdiff/internal/comparer"
	"github.com/user/envdiff/internal/differ"
)

func TestEqual_ExactMatch(t *testing.T) {
	opts := comparer.DefaultOptions()
	if !comparer.Equal("hello", "hello", opts) {
		t.Error("expected equal")
	}
}

func TestEqual_TrimSpace(t *testing.T) {
	opts := comparer.DefaultOptions() // TrimSpace: true
	if !comparer.Equal("  hello  ", "hello", opts) {
		t.Error("expected equal after trim")
	}
}

func TestEqual_NoTrimSpace(t *testing.T) {
	opts := comparer.Options{TrimSpace: false}
	if comparer.Equal("  hello", "hello", opts) {
		t.Error("expected not equal without trim")
	}
}

func TestEqual_CaseInsensitive(t *testing.T) {
	opts := comparer.Options{CaseInsensitive: true, TrimSpace: false}
	if !comparer.Equal("Hello", "hello", opts) {
		t.Error("expected equal case-insensitive")
	}
}

func TestEqual_CaseSensitive(t *testing.T) {
	opts := comparer.Options{CaseInsensitive: false, TrimSpace: false}
	if comparer.Equal("Hello", "hello", opts) {
		t.Error("expected not equal case-sensitive")
	}
}

func TestEqual_NormalizeBools(t *testing.T) {
	opts := comparer.Options{NormalizeBools: true, TrimSpace: true}
	cases := [][2]string{
		{"true", "1"}, {"yes", "true"}, {"false", "0"}, {"no", "false"}, {"1", "yes"},
	}
	for _, c := range cases {
		if !comparer.Equal(c[0], c[1], opts) {
			t.Errorf("expected %q == %q with bool normalization", c[0], c[1])
		}
	}
}

func TestFilterMismatches_RemovesEquivalent(t *testing.T) {
	result := differ.Result{
		Mismatched: []differ.Mismatch{
			{Key: "A", LeftValue: " hello ", RightValue: "hello"},
			{Key: "B", LeftValue: "world", RightValue: "earth"},
		},
	}
	opts := comparer.DefaultOptions() // TrimSpace: true
	out := comparer.FilterMismatches(result, opts)
	if len(out.Mismatched) != 1 {
		t.Fatalf("expected 1 mismatch, got %d", len(out.Mismatched))
	}
	if out.Mismatched[0].Key != "B" {
		t.Errorf("expected key B, got %s", out.Mismatched[0].Key)
	}
}

func TestFilterMismatches_PreservesMissingKeys(t *testing.T) {
	result := differ.Result{
		MissingInLeft:  []string{"X"},
		MissingInRight: []string{"Y"},
		Mismatched:     []differ.Mismatch{},
	}
	opts := comparer.DefaultOptions()
	out := comparer.FilterMismatches(result, opts)
	if len(out.MissingInLeft) != 1 || out.MissingInLeft[0] != "X" {
		t.Error("expected MissingInLeft to be preserved")
	}
	if len(out.MissingInRight) != 1 || out.MissingInRight[0] != "Y" {
		t.Error("expected MissingInRight to be preserved")
	}
}
