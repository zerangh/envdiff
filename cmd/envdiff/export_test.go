package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeExportTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "*.env")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestRunExport_MissingArgs(t *testing.T) {
	err := runExport([]string{"only-one.env"})
	if err == nil || !strings.Contains(err.Error(), "two .env") {
		t.Errorf("expected missing-args error, got: %v", err)
	}
}

func TestRunExport_InvalidFile(t *testing.T) {
	err := runExport([]string{"/nonexistent/a.env", "/nonexistent/b.env"})
	if err == nil {
		t.Fatal("expected error for invalid file")
	}
}

func TestRunExport_MarkdownToStdout(t *testing.T) {
	left := writeExportTempEnv(t, "APP_ENV=production\nDB_HOST=localhost\n")
	right := writeExportTempEnv(t, "APP_ENV=staging\n")

	err := runExport([]string{left, right, "--format=markdown"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunExport_JSONToFile(t *testing.T) {
	left := writeExportTempEnv(t, "KEY=val\nONLY_LEFT=x\n")
	right := writeExportTempEnv(t, "KEY=other\nONLY_RIGHT=y\n")

	outFile := filepath.Join(t.TempDir(), "out.json")
	err := runExport([]string{left, right, "--format=json", "--output=" + outFile})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("reading output file: %v", err)
	}
	if !strings.Contains(string(data), "KEY") {
		t.Errorf("expected KEY in JSON output, got: %s", string(data))
	}
}

func TestRunExport_EnvFormat(t *testing.T) {
	left := writeExportTempEnv(t, "APP_ENV=production\nSECRET=abc\n")
	right := writeExportTempEnv(t, "APP_ENV=staging\nSECRET=xyz\n")

	err := runExport([]string{left, right, "--format=env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
