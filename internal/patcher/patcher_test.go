package patcher_test

import (
	"os"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/patcher"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func makeDiff(missingRight, missingLeft []string, mismatched []differ.Mismatch) differ.Result {
	return differ.Result{
		MissingInRight: missingRight,
		MissingInLeft:  missingLeft,
		Mismatched:     mismatched,
	}
}

func TestPatch_AddsMissingKeys(t *testing.T) {
	path := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	patch := map[string]string{"NEW_KEY": "hello"}
	diff := makeDiff([]string{"NEW_KEY"}, nil, nil)

	_, res, err := patcher.Patch(path, patch, diff, patcher.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Added) != 1 || res.Added[0] != "NEW_KEY" {
		t.Errorf("expected Added=[NEW_KEY], got %v", res.Added)
	}
	got, _ := os.ReadFile(path)
	if !strings.Contains(string(got), "NEW_KEY=hello") {
		t.Errorf("file missing NEW_KEY=hello, got:\n%s", got)
	}
}

func TestPatch_UpdatesMismatchedWhenEnabled(t *testing.T) {
	path := writeTempEnv(t, "FOO=old\n")
	patch := map[string]string{"FOO": "new"}
	diff := makeDiff(nil, nil, []differ.Mismatch{{Key: "FOO", Left: "new", Right: "old"}})

	_, res, err := patcher.Patch(path, patch, diff, patcher.Options{UpdateMismatched: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Updated) != 1 || res.Updated[0] != "FOO" {
		t.Errorf("expected Updated=[FOO], got %v", res.Updated)
	}
	got, _ := os.ReadFile(path)
	if !strings.Contains(string(got), "FOO=new") {
		t.Errorf("expected FOO=new in file, got:\n%s", got)
	}
}

func TestPatch_SkipsMismatchedWhenDisabled(t *testing.T) {
	path := writeTempEnv(t, "FOO=old\n")
	patch := map[string]string{"FOO": "new"}
	diff := makeDiff(nil, nil, []differ.Mismatch{{Key: "FOO", Left: "new", Right: "old"}})

	_, res, err := patcher.Patch(path, patch, diff, patcher.Options{UpdateMismatched: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "FOO" {
		t.Errorf("expected Skipped=[FOO], got %v", res.Skipped)
	}
	got, _ := os.ReadFile(path)
	if !strings.Contains(string(got), "FOO=old") {
		t.Errorf("expected original FOO=old preserved, got:\n%s", got)
	}
}

func TestPatch_DryRunDoesNotWriteFile(t *testing.T) {
	original := "FOO=bar\n"
	path := writeTempEnv(t, original)
	patch := map[string]string{"NEW": "val"}
	diff := makeDiff([]string{"NEW"}, nil, nil)

	output, _, err := patcher.Patch(path, patch, diff, patcher.Options{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(output, "NEW=val") {
		t.Errorf("expected output to contain NEW=val, got: %s", output)
	}
	got, _ := os.ReadFile(path)
	if string(got) != original {
		t.Errorf("dry-run should not modify file; got:\n%s", got)
	}
}
