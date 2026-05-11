package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/envdiff/internal/differ"
	"github.com/yourusername/envdiff/internal/snapshot"
)

func makeResult() differ.Result {
	return differ.Result{
		MissingInRight: []string{"DB_HOST"},
		MissingInLeft:  []string{"API_SECRET"},
		Mismatched: []differ.Mismatch{
			{Key: "PORT", LeftValue: "8080", RightValue: "9090"},
		},
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	result := makeResult()
	err := snapshot.Save(path, ".env.dev", ".env.prod", result)
	if err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	rec, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if rec.LeftFile != ".env.dev" {
		t.Errorf("LeftFile = %q, want .env.dev", rec.LeftFile)
	}
	if rec.RightFile != ".env.prod" {
		t.Errorf("RightFile = %q, want .env.prod", rec.RightFile)
	}
	if len(rec.Result.MissingInRight) != 1 || rec.Result.MissingInRight[0] != "DB_HOST" {
		t.Errorf("MissingInRight = %v, want [DB_HOST]", rec.Result.MissingInRight)
	}
	if len(rec.Result.MissingInLeft) != 1 || rec.Result.MissingInLeft[0] != "API_SECRET" {
		t.Errorf("MissingInLeft = %v, want [API_SECRET]", rec.Result.MissingInLeft)
	}
	if len(rec.Result.Mismatched) != 1 || rec.Result.Mismatched[0].Key != "PORT" {
		t.Errorf("Mismatched = %v, want [{PORT 8080 9090}]", rec.Result.Mismatched)
	}
	if rec.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
	if time.Since(rec.CreatedAt) > 5*time.Second {
		t.Error("CreatedAt is unexpectedly old")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Fatal("expected error loading missing file, got nil")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte("not-json"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := snapshot.Load(path)
	if err == nil {
		t.Fatal("expected error on invalid JSON, got nil")
	}
}

func TestSave_UnwritablePath(t *testing.T) {
	err := snapshot.Save("/no/such/directory/snap.json", "a", "b", differ.Result{})
	if err == nil {
		t.Fatal("expected error writing to unwritable path, got nil")
	}
}
