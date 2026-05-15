package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeSplitTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writeSplitTempEnv: %v", err)
	}
	return path
}

func TestRunSplit_MissingArgs(t *testing.T) {
	err := runSplit([]string{})
	if err == nil {
		t.Fatal("expected error for missing args")
	}
	if !strings.Contains(err.Error(), "usage") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunSplit_InvalidFile(t *testing.T) {
	err := runSplit([]string{"/nonexistent/.env", "/also/missing/.env"})
	if err == nil {
		t.Fatal("expected error for invalid file")
	}
}

func TestRunSplit_NoDifferences(t *testing.T) {
	content := "APP_HOST=localhost\nDB_PORT=5432\n"
	left := writeSplitTempEnv(t, content)
	right := writeSplitTempEnv(t, content)

	err := runSplit([]string{left, right})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunSplit_PrefixStrategy(t *testing.T) {
	left := writeSplitTempEnv(t, "APP_HOST=localhost\nDB_PORT=5432\n")
	right := writeSplitTempEnv(t, "APP_HOST=remotehost\nDB_PORT=5432\n")

	err := runSplit([]string{left, right, "--strategy", "prefix"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunSplit_AlphaStrategy(t *testing.T) {
	left := writeSplitTempEnv(t, "APP_HOST=localhost\nDB_PORT=5432\n")
	right := writeSplitTempEnv(t, "APP_HOST=remotehost\n")

	err := runSplit([]string{left, right, "--strategy", "alpha"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunSplit_JSONFormat(t *testing.T) {
	left := writeSplitTempEnv(t, "APP_ENV=dev\nDB_NAME=mydb\n")
	right := writeSplitTempEnv(t, "APP_ENV=prod\n")

	err := runSplit([]string{left, right, "--format", "json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
