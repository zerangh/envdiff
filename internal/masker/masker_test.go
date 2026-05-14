package masker_test

import (
	"testing"

	"github.com/your-org/envdiff/internal/differ"
	"github.com/your-org/envdiff/internal/masker"
)

func TestIsSensitive_KnownPatterns(t *testing.T) {
	opts := masker.DefaultOptions()
	cases := []struct {
		key  string
		want bool
	}{
		{"DB_PASSWORD", true},
		{"API_TOKEN", true},
		{"SECRET_KEY", true},
		{"PRIVATE_KEY", true},
		{"APP_NAME", false},
		{"PORT", false},
	}
	for _, tc := range cases {
		got := masker.IsSensitive(tc.key, opts.Patterns)
		if got != tc.want {
			t.Errorf("IsSensitive(%q) = %v; want %v", tc.key, got, tc.want)
		}
	}
}

func TestMaskEnv_MasksSensitiveKeys(t *testing.T) {
	env := map[string]string{
		"APP_NAME":    "myapp",
		"DB_PASSWORD": "supersecret",
		"PORT":        "8080",
		"API_TOKEN":   "tok_abc123",
	}
	opts := masker.DefaultOptions()
	result := masker.MaskEnv(env, opts)

	if result["APP_NAME"] != "myapp" {
		t.Errorf("APP_NAME should be unmasked, got %q", result["APP_NAME"])
	}
	if result["PORT"] != "8080" {
		t.Errorf("PORT should be unmasked, got %q", result["PORT"])
	}
	if result["DB_PASSWORD"] == "supersecret" {
		t.Error("DB_PASSWORD should be masked")
	}
	if result["API_TOKEN"] == "tok_abc123" {
		t.Error("API_TOKEN should be masked")
	}
}

func TestMaskEnv_RevealSuffix(t *testing.T) {
	env := map[string]string{
		"API_TOKEN": "tok_abc123",
	}
	opts := masker.DefaultOptions()
	opts.RevealSuffix = 3
	result := masker.MaskEnv(env, opts)

	val := result["API_TOKEN"]
	if len(val) == 0 {
		t.Fatal("masked value should not be empty")
	}
	if val[len(val)-3:] != "123" {
		t.Errorf("expected suffix '123', got %q", val[len(val)-3:])
	}
}

func TestMaskResult_MasksMismatchedValues(t *testing.T) {
	result := differ.Result{
		MissingInLeft:  []string{"ONLY_RIGHT"},
		MissingInRight: []string{"ONLY_LEFT"},
		Mismatched: []differ.Mismatch{
			{Key: "DB_PASSWORD", LeftValue: "pass1", RightValue: "pass2"},
			{Key: "APP_NAME", LeftValue: "foo", RightValue: "bar"},
		},
	}
	opts := masker.DefaultOptions()
	masked := masker.MaskResult(result, opts)

	if masked.Mismatched[0].LeftValue == "pass1" {
		t.Error("DB_PASSWORD left value should be masked")
	}
	if masked.Mismatched[0].RightValue == "pass2" {
		t.Error("DB_PASSWORD right value should be masked")
	}
	if masked.Mismatched[1].LeftValue != "foo" {
		t.Errorf("APP_NAME left value should be unmasked, got %q", masked.Mismatched[1].LeftValue)
	}
	if len(masked.MissingInLeft) != 1 || masked.MissingInLeft[0] != "ONLY_RIGHT" {
		t.Error("MissingInLeft should be preserved unchanged")
	}
}

func TestMaskEnv_EmptyValue(t *testing.T) {
	env := map[string]string{
		"API_TOKEN": "",
	}
	opts := masker.DefaultOptions()
	result := masker.MaskEnv(env, opts)
	if result["API_TOKEN"] != "" {
		t.Errorf("empty sensitive value should remain empty, got %q", result["API_TOKEN"])
	}
}
