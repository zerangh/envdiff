package main

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return path
}

func TestRun_NoDifferences(t *testing.T) {
	left := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	right := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")

	os.Args = []string{"envdiff", left, right}
	if err := run(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestRun_MissingArgs(t *testing.T) {
	os.Args = []string{"envdiff"}
	if err := run(); err == nil {
		t.Fatal("expected error for missing args, got nil")
	}
}

func TestRun_InvalidFile(t *testing.T) {
	os.Args = []string{"envdiff", "nonexistent_left.env", "nonexistent_right.env"}
	if err := run(); err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestRun_JSONFormat(t *testing.T) {
	left := writeTempEnv(t, "FOO=bar\n")
	right := writeTempEnv(t, "FOO=bar\n")

	os.Args = []string{"envdiff", "-format", "json", left, right}
	if err := run(); err != nil {
		t.Fatalf("expected no error with json format, got: %v", err)
	}
}

func TestRun_CustomLabels(t *testing.T) {
	left := writeTempEnv(t, "KEY=value\n")
	right := writeTempEnv(t, "KEY=value\n")

	os.Args = []string{"envdiff", "-left-label", "staging", "-right-label", "production", left, right}
	if err := run(); err != nil {
		t.Fatalf("expected no error with custom labels, got: %v", err)
	}
}

func TestRun_InvalidFormat(t *testing.T) {
	left := writeTempEnv(t, "FOO=bar\n")
	right := writeTempEnv(t, "FOO=bar\n")

	os.Args = []string{"envdiff", "-format", "xml", left, right}
	if err := run(); err == nil {
		t.Fatal("expected error for invalid format, got nil")
	}
}
