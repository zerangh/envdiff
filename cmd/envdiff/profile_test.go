package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeProfileTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeProfileTempEnv: %v", err)
	}
	return p
}

func TestRunProfile_MissingArgs(t *testing.T) {
	err := runProfile([]string{})
	if err == nil || !strings.Contains(err.Error(), "at least 2") {
		t.Errorf("expected at-least-2 error, got %v", err)
	}
}

func TestRunProfile_SingleFile(t *testing.T) {
	f := writeProfileTempEnv(t, "KEY=val\n")
	err := runProfile([]string{f})
	if err == nil {
		t.Error("expected error for single file")
	}
}

func TestRunProfile_InvalidFile(t *testing.T) {
	err := runProfile([]string{"/nonexistent/.env", "/also/missing/.env"})
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestRunProfile_TextOutput(t *testing.T) {
	f1 := writeProfileTempEnv(t, "APP=myapp\nDB=localhost\n")
	f2 := writeProfileTempEnv(t, "APP=myapp\nDB=remote\nSECRET=x\n")

	// Redirect stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runProfile([]string{f1, f2})

	w.Close()
	os.Stdout = old

	var buf strings.Builder
	tmp := make([]byte, 1024)
	for {
		n, e := r.Read(tmp)
		buf.Write(tmp[:n])
		if e != nil {
			break
		}
	}

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Total keys") {
		t.Errorf("expected 'Total keys' in output, got:\n%s", out)
	}
	if !strings.Contains(out, "DB") {
		t.Errorf("expected DB key in output")
	}
}

func TestRunProfile_JSONFlag(t *testing.T) {
	f1 := writeProfileTempEnv(t, "APP=myapp\n")
	f2 := writeProfileTempEnv(t, "APP=myapp\n")

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runProfile([]string{"--json", f1, f2})

	w.Close()
	os.Stdout = old

	var buf strings.Builder
	tmp := make([]byte, 1024)
	for {
		n, e := r.Read(tmp)
		buf.Write(tmp[:n])
		if e != nil {
			break
		}
	}

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.HasPrefix(strings.TrimSpace(out), "{") {
		t.Errorf("expected JSON output, got:\n%s", out)
	}
}
