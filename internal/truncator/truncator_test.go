package truncator_test

import (
	"strings"
	"testing"

	"github.com/envdiff/envdiff/internal/differ"
	"github.com/envdiff/envdiff/internal/truncator"
)

func TestTruncateValue_ShortValue(t *testing.T) {
	opts := truncator.DefaultOptions()
	v := "short"
	got := truncator.TruncateValue(v, opts)
	if got != v {
		t.Errorf("expected %q, got %q", v, got)
	}
}

func TestTruncateValue_ExactLength(t *testing.T) {
	opts := truncator.Options{MaxLen: 5}
	v := "hello"
	got := truncator.TruncateValue(v, opts)
	if got != v {
		t.Errorf("expected value unchanged at exact length, got %q", got)
	}
}

func TestTruncateValue_LongValue(t *testing.T) {
	opts := truncator.Options{MaxLen: 10}
	v := "this is a very long environment variable value"
	got := truncator.TruncateValue(v, opts)
	if len(got) != 13 { // 10 + len("...")
		t.Errorf("expected length 13, got %d: %q", len(got), got)
	}
	if !strings.HasSuffix(got, "...") {
		t.Errorf("expected ellipsis suffix, got %q", got)
	}
}

func TestTruncateValue_ZeroMaxUsesDefault(t *testing.T) {
	opts := truncator.Options{MaxLen: 0}
	long := strings.Repeat("x", 100)
	got := truncator.TruncateValue(long, opts)
	if !strings.HasSuffix(got, "...") {
		t.Errorf("expected truncation with zero MaxLen, got %q", got)
	}
}

func TestTruncateEnv_AllValuesTruncated(t *testing.T) {
	opts := truncator.Options{MaxLen: 5}
	env := map[string]string{
		"SHORT": "hi",
		"LONG":  "this value is definitely too long",
	}
	out := truncator.TruncateEnv(env, opts)
	if out["SHORT"] != "hi" {
		t.Errorf("SHORT should be unchanged, got %q", out["SHORT"])
	}
	if !strings.HasSuffix(out["LONG"], "...") {
		t.Errorf("LONG should be truncated, got %q", out["LONG"])
	}
}

func TestTruncateEnv_DoesNotMutateOriginal(t *testing.T) {
	opts := truncator.Options{MaxLen: 3}
	env := map[string]string{"KEY": "longvalue"}
	_ = truncator.TruncateEnv(env, opts)
	if env["KEY"] != "longvalue" {
		t.Error("original env map was mutated")
	}
}

func TestTruncateResult_MismatchedValues(t *testing.T) {
	opts := truncator.Options{MaxLen: 8}
	r := differ.Result{
		MissingInLeft:  []string{"A"},
		MissingInRight: []string{"B"},
		MismatchedValues: []differ.Mismatch{
			{Key: "K", LeftValue: "short", RightValue: "a very long right value here"},
		},
	}
	out := truncator.TruncateResult(r, opts)
	if out.MismatchedValues[0].LeftValue != "short" {
		t.Errorf("short left value should be unchanged, got %q", out.MismatchedValues[0].LeftValue)
	}
	if !strings.HasSuffix(out.MismatchedValues[0].RightValue, "...") {
		t.Errorf("long right value should be truncated, got %q", out.MismatchedValues[0].RightValue)
	}
	// original untouched
	if r.MismatchedValues[0].RightValue != "a very long right value here" {
		t.Error("original result was mutated")
	}
}

func TestTruncateResult_PreservesMissingSlices(t *testing.T) {
	opts := truncator.DefaultOptions()
	r := differ.Result{
		MissingInLeft:  []string{"X", "Y"},
		MissingInRight: []string{"Z"},
	}
	out := truncator.TruncateResult(r, opts)
	if len(out.MissingInLeft) != 2 || len(out.MissingInRight) != 1 {
		t.Error("missing key slices not preserved correctly")
	}
}
