package redactor_test

import (
	"regexp"
	"testing"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/redactor"
)

func TestIsSensitive_KnownPatterns(t *testing.T) {
	cases := []struct {
		key  string
		want bool
	}{
		{"DB_PASSWORD", true},
		{"API_SECRET", true},
		{"AUTH_TOKEN", true},
		{"PRIVATE_KEY", true},
		{"APP_NAME", false},
		{"PORT", false},
		{"DATABASE_URL", false},
	}
	for _, tc := range cases {
		got := redactor.IsSensitive(tc.key, nil)
		if got != tc.want {
			t.Errorf("IsSensitive(%q) = %v, want %v", tc.key, got, tc.want)
		}
	}
}

func TestIsSensitive_ExtraPatterns(t *testing.T) {
	extra := []*regexp.Regexp{regexp.MustCompile(`(?i)ssn`)}
	if !redactor.IsSensitive("USER_SSN", extra) {
		t.Error("expected USER_SSN to be sensitive with extra pattern")
	}
	if redactor.IsSensitive("USER_SSN", nil) {
		t.Error("expected USER_SSN NOT sensitive without extra pattern")
	}
}

func TestRedactEnv_MasksSensitiveKeys(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "s3cr3t",
		"APP_NAME":    "myapp",
		"API_TOKEN":   "tok123",
	}
	got := redactor.RedactEnv(env, redactor.Options{})
	if got["DB_PASSWORD"] != "***" {
		t.Errorf("expected DB_PASSWORD masked, got %q", got["DB_PASSWORD"])
	}
	if got["API_TOKEN"] != "***" {
		t.Errorf("expected API_TOKEN masked, got %q", got["API_TOKEN"])
	}
	if got["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME unchanged, got %q", got["APP_NAME"])
	}
}

func TestRedactEnv_CustomMask(t *testing.T) {
	env := map[string]string{"DB_PASSWORD": "hunter2"}
	got := redactor.RedactEnv(env, redactor.Options{Mask: "[REDACTED]"})
	if got["DB_PASSWORD"] != "[REDACTED]" {
		t.Errorf("got %q, want [REDACTED]", got["DB_PASSWORD"])
	}
}

func TestRedactResult_MasksMismatchedSensitiveValues(t *testing.T) {
	r := differ.Result{
		Mismatched: []differ.Mismatch{
			{Key: "DB_PASSWORD", LeftValue: "old", RightValue: "new"},
			{Key: "APP_NAME", LeftValue: "foo", RightValue: "bar"},
		},
	}
	got := redactor.RedactResult(r, redactor.Options{})
	if got.Mismatched[0].LeftValue != "***" {
		t.Errorf("expected LeftValue masked, got %q", got.Mismatched[0].LeftValue)
	}
	if got.Mismatched[1].LeftValue != "foo" {
		t.Errorf("expected APP_NAME LeftValue unchanged, got %q", got.Mismatched[1].LeftValue)
	}
	if got.Mismatched[1].RightValue != "bar" {
		t.Errorf("expected APP_NAME RightValue unchanged, got %q", got.Mismatched[1].RightValue)
	}
}

func TestRedactResult_PreservesOtherFields(t *testing.T) {
	r := differ.Result{
		MissingInLeft:  []string{"FOO"},
		MissingInRight: []string{"BAR"},
		Mismatched:     nil,
	}
	got := redactor.RedactResult(r, redactor.Options{})
	if len(got.MissingInLeft) != 1 || got.MissingInLeft[0] != "FOO" {
		t.Errorf("MissingInLeft not preserved: %v", got.MissingInLeft)
	}
	if len(got.MissingInRight) != 1 || got.MissingInRight[0] != "BAR" {
		t.Errorf("MissingInRight not preserved: %v", got.MissingInRight)
	}
}
