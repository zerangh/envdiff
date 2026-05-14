package normalizer_test

import (
	"testing"

	"github.com/your-org/envdiff/internal/differ"
	"github.com/your-org/envdiff/internal/normalizer"
)

func TestNormalizeValue_TrimSpace(t *testing.T) {
	opts := normalizer.Options{TrimSpace: true}
	got := normalizer.NormalizeValue("  hello  ", opts)
	if got != "hello" {
		t.Errorf("expected 'hello', got %q", got)
	}
}

func TestNormalizeValue_NormalizeBools(t *testing.T) {
	opts := normalizer.Options{NormalizeBools: true}
	cases := map[string]string{
		"yes": "true", "YES": "true", "1": "true", "on": "true",
		"no": "false", "NO": "false", "0": "false", "off": "false",
		"maybe": "maybe",
	}
	for input, want := range cases {
		got := normalizer.NormalizeValue(input, opts)
		if got != want {
			t.Errorf("NormalizeValue(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestNormalizeValue_LowercaseValues(t *testing.T) {
	opts := normalizer.Options{LowercaseValues: true}
	got := normalizer.NormalizeValue("Hello World", opts)
	if got != "hello world" {
		t.Errorf("expected 'hello world', got %q", got)
	}
}

func TestNormalizeEnv_AllValues(t *testing.T) {
	env := map[string]string{
		"FOO": "  bar  ",
		"ENABLED": "yes",
	}
	opts := normalizer.DefaultOptions()
	out := normalizer.NormalizeEnv(env, opts)
	if out["FOO"] != "bar" {
		t.Errorf("FOO: expected 'bar', got %q", out["FOO"])
	}
	if out["ENABLED"] != "true" {
		t.Errorf("ENABLED: expected 'true', got %q", out["ENABLED"])
	}
}

func TestNormalizeResult_RemovesEquivalentMismatches(t *testing.T) {
	r := differ.Result{
		Mismatched: []differ.Mismatch{
			{Key: "FLAG", LeftValue: "yes", RightValue: "true"},
			{Key: "NAME", LeftValue: "alice", RightValue: "bob"},
		},
	}
	out := normalizer.NormalizeResult(r, normalizer.DefaultOptions())
	if len(out.Mismatched) != 1 {
		t.Fatalf("expected 1 mismatch after normalization, got %d", len(out.Mismatched))
	}
	if out.Mismatched[0].Key != "NAME" {
		t.Errorf("expected remaining mismatch to be NAME, got %q", out.Mismatched[0].Key)
	}
}

func TestNormalizeResult_PreservesMissingKeys(t *testing.T) {
	r := differ.Result{
		MissingInLeft:  []string{"A"},
		MissingInRight: []string{"B"},
	}
	out := normalizer.NormalizeResult(r, normalizer.DefaultOptions())
	if len(out.MissingInLeft) != 1 || out.MissingInLeft[0] != "A" {
		t.Errorf("MissingInLeft not preserved: %v", out.MissingInLeft)
	}
	if len(out.MissingInRight) != 1 || out.MissingInRight[0] != "B" {
		t.Errorf("MissingInRight not preserved: %v", out.MissingInRight)
	}
}

func TestNormalizeResult_NoMismatches(t *testing.T) {
	r := differ.Result{}
	out := normalizer.NormalizeResult(r, normalizer.DefaultOptions())
	if len(out.Mismatched) != 0 {
		t.Errorf("expected no mismatches, got %d", len(out.Mismatched))
	}
}
