package watcher_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/envdiff/internal/watcher"
)

func writeTempEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return path
}

func TestWatch_InvalidReferenceFile(t *testing.T) {
	dir := t.TempDir()
	watched := writeTempEnv(t, dir, "watched.env", "KEY=val\n")

	_, _, err := watcher.Watch("/nonexistent/ref.env", watched, nil)
	if err == nil {
		t.Fatal("expected error for missing reference file, got nil")
	}
}

func TestWatch_InvalidWatchedFile(t *testing.T) {
	dir := t.TempDir()
	ref := writeTempEnv(t, dir, "ref.env", "KEY=val\n")

	_, _, err := watcher.Watch(ref, "/nonexistent/watched.env", nil)
	if err == nil {
		t.Fatal("expected error for missing watched file, got nil")
	}
}

func TestWatch_DetectsChange(t *testing.T) {
	dir := t.TempDir()
	ref := writeTempEnv(t, dir, "ref.env", "KEY=value\nOTHER=foo\n")
	watched := writeTempEnv(t, dir, "watched.env", "KEY=value\nOTHER=foo\n")

	opts := &watcher.Options{PollInterval: 50 * time.Millisecond}
	events, stop, err := watcher.Watch(ref, watched, opts)
	if err != nil {
		t.Fatalf("Watch: %v", err)
	}
	defer stop()

	// Modify the watched file after a short delay.
	time.Sleep(80 * time.Millisecond)
	if err := os.WriteFile(watched, []byte("KEY=changed\nOTHER=foo\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	select {
	case ev, ok := <-events:
		if !ok {
			t.Fatal("events channel closed unexpectedly")
		}
		if ev.File != watched {
			t.Errorf("Event.File = %q, want %q", ev.File, watched)
		}
		if len(ev.Result.Mismatched) != 1 {
			t.Errorf("expected 1 mismatch, got %d", len(ev.Result.Mismatched))
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for change event")
	}
}

func TestWatch_NoEventWhenUnchanged(t *testing.T) {
	dir := t.TempDir()
	ref := writeTempEnv(t, dir, "ref.env", "KEY=value\n")
	watched := writeTempEnv(t, dir, "watched.env", "KEY=value\n")

	opts := &watcher.Options{PollInterval: 50 * time.Millisecond}
	events, stop, err := watcher.Watch(ref, watched, opts)
	if err != nil {
		t.Fatalf("Watch: %v", err)
	}

	time.Sleep(200 * time.Millisecond)
	stop()

	select {
	case ev, ok := <-events:
		if ok {
			t.Errorf("unexpected event received: %+v", ev)
		}
	default:
		// channel drained with no events — expected
	}
}
