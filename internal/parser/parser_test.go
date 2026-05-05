package parser

import (
	"os"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "envdiff-*.env")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestParseFile_Basic(t *testing.T) {
	path := writeTempEnv(t, "KEY1=value1\nKEY2=value2\n")
	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["KEY1"] != "value1" {
		t.Errorf("expected KEY1=value1, got %q", env["KEY1"])
	}
	if env["KEY2"] != "value2" {
		t.Errorf("expected KEY2=value2, got %q", env["KEY2"])
	}
}

func TestParseFile_SkipsCommentsAndBlanks(t *testing.T) {
	path := writeTempEnv(t, "# comment\n\nKEY=val\n")
	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env) != 1 {
		t.Errorf("expected 1 key, got %d", len(env))
	}
}

func TestParseFile_QuotedValues(t *testing.T) {
	path := writeTempEnv(t, `KEY1="hello world"` + "\n" + `KEY2='single'` + "\n")
	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["KEY1"] != "hello world" {
		t.Errorf("expected 'hello world', got %q", env["KEY1"])
	}
	if env["KEY2"] != "single" {
		t.Errorf("expected 'single', got %q", env["KEY2"])
	}
}

func TestParseFile_MalformedLine(t *testing.T) {
	path := writeTempEnv(t, "BADLINE\n")
	_, err := ParseFile(path)
	if err == nil {
		t.Error("expected error for malformed line, got nil")
	}
}

func TestParseFile_NotFound(t *testing.T) {
	_, err := ParseFile("/nonexistent/path/.env")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestParseFile_EmptyKey(t *testing.T) {
	path := writeTempEnv(t, "=value\n")
	_, err := ParseFile(path)
	if err == nil {
		t.Error("expected error for empty key, got nil")
	}
}

func TestParseFile_EmptyFile(t *testing.T) {
	path := writeTempEnv(t, "")
	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error for empty file: %v", err)
	}
	if len(env) != 0 {
		t.Errorf("expected 0 keys for empty file, got %d", len(env))
	}
}
