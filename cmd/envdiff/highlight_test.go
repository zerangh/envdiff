package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeHighlightTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRunHighlight_MissingArgs(t *testing.T) {
	err := runHighlight([]string{}, "text")
	if err == nil || !strings.Contains(err.Error(), "highlight requires") {
		t.Errorf("expected arg error, got %v", err)
	}
}

func TestRunHighlight_InvalidFile(t *testing.T) {
	err := runHighlight([]string{"/nonexistent/.env", "/also/missing"}, "text")
	if err == nil {
		t.Error("expected error for invalid file")
	}
}

func TestRunHighlight_NoDifferences(t *testing.T) {
	a := writeHighlightTempEnv(t, "APP_ENV=production\nPORT=8080\n")
	b := writeHighlightTempEnv(t, "APP_ENV=production\nPORT=8080\n")
	err := runHighlight([]string{a, b}, "text")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunHighlight_TextOutput(t *testing.T) {
	a := writeHighlightTempEnv(t, "APP_ENV=production\nPORT=8080\n")
	b := writeHighlightTempEnv(t, "APP_ENV=staging\nPORT=8080\n")

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runHighlight([]string{a, b}, "text")
	w.Close()
	os.Stdout = old

	var buf strings.Builder
	io.Copy(&buf, r)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "APP_ENV") {
		t.Errorf("expected APP_ENV in output, got: %s", buf.String())
	}
}

func TestRunHighlight_JSONOutput(t *testing.T) {
	a := writeHighlightTempEnv(t, "DB_HOST=host1\n")
	b := writeHighlightTempEnv(t, "DB_HOST=host2\n")

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runHighlight([]string{a, b}, "json")
	w.Close()
	os.Stdout = old

	var buf strings.Builder
	io.Copy(&buf, r)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "Highlights") {
		t.Errorf("expected JSON with Highlights key, got: %s", buf.String())
	}
}
