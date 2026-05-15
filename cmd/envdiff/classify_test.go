package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeClassifyTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRunClassify_MissingArgs(t *testing.T) {
	err := runClassify([]string{}, "text")
	if err == nil || !strings.Contains(err.Error(), "usage") {
		t.Errorf("expected usage error, got %v", err)
	}
}

func TestRunClassify_InvalidFile(t *testing.T) {
	err := runClassify([]string{"/nonexistent/.env", "/also/none"}, "text")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestRunClassify_NoDifferences(t *testing.T) {
	f1 := writeClassifyTempEnv(t, "APP_PORT=8080\n")
	f2 := writeClassifyTempEnv(t, "APP_PORT=8080\n")
	err := runClassify([]string{f1, f2}, "text")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunClassify_TextOutput(t *testing.T) {
	f1 := writeClassifyTempEnv(t, "DB_HOST=localhost\nJWT_SECRET=abc\n")
	f2 := writeClassifyTempEnv(t, "DB_HOST=prod\n")

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runClassify([]string{f1, f2}, "text")
	w.Close()
	os.Stdout = old

	var buf strings.Builder
	b := make([]byte, 4096)
	for {
		n, e := r.Read(b)
		buf.Write(b[:n])
		if e != nil {
			break
		}
	}

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "auth") && !strings.Contains(out, "database") {
		t.Errorf("expected category headers in output, got: %s", out)
	}
}

func TestRunClassify_JSONOutput(t *testing.T) {
	f1 := writeClassifyTempEnv(t, "DB_HOST=localhost\n")
	f2 := writeClassifyTempEnv(t, "DB_HOST=prod\n")

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runClassify([]string{f1, f2}, "json")
	w.Close()
	os.Stdout = old

	var buf strings.Builder
	b := make([]byte, 4096)
	for {
		n, e := r.Read(b)
		buf.Write(b[:n])
		if e != nil {
			break
		}
	}

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "classes") || !strings.Contains(out, "categories") {
		t.Errorf("expected JSON keys in output, got: %s", out)
	}
}
