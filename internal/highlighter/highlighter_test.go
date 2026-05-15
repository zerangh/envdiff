package highlighter_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/highlighter"
)

func makeResult(mismatched []differ.Mismatch) differ.Result {
	return differ.Result{Mismatched: mismatched}
}

func TestCompute_Empty(t *testing.T) {
	r := makeResult(nil)
	hr := highlighter.Compute(r)
	if len(hr.Highlights) != 0 {
		t.Fatalf("expected 0 highlights, got %d", len(hr.Highlights))
	}
}

func TestCompute_SingleMismatch(t *testing.T) {
	r := makeResult([]differ.Mismatch{
		{Key: "APP_ENV", Left: "production", Right: "staging"},
	})
	hr := highlighter.Compute(r)
	if len(hr.Highlights) != 1 {
		t.Fatalf("expected 1 highlight, got %d", len(hr.Highlights))
	}
	h := hr.Highlights[0]
	if h.Key != "APP_ENV" {
		t.Errorf("expected key APP_ENV, got %s", h.Key)
	}
	if h.Left != "production" {
		t.Errorf("unexpected Left: %s", h.Left)
	}
	if h.Right != "staging" {
		t.Errorf("unexpected Right: %s", h.Right)
	}
}

func TestCompute_SegmentsMarkRemovedAndAdded(t *testing.T) {
	r := makeResult([]differ.Mismatch{
		{Key: "DB_URL", Left: "host1 port", Right: "host2 port"},
	})
	hr := highlighter.Compute(r)
	segs := hr.Highlights[0].Diff

	var removed, added []string
	for _, s := range segs {
		if s.Removed {
			removed = append(removed, s.Text)
		}
		if s.Added {
			added = append(added, s.Text)
		}
	}
	if len(removed) != 1 || removed[0] != "host1" {
		t.Errorf("expected removed=[host1], got %v", removed)
	}
	if len(added) != 1 || added[0] != "host2" {
		t.Errorf("expected added=[host2], got %v", added)
	}
}

func TestCompute_MultipleMismatches(t *testing.T) {
	r := makeResult([]differ.Mismatch{
		{Key: "A", Left: "x", Right: "y"},
		{Key: "B", Left: "1", Right: "2"},
	})
	hr := highlighter.Compute(r)
	if len(hr.Highlights) != 2 {
		t.Fatalf("expected 2 highlights, got %d", len(hr.Highlights))
	}
}

func TestFormat_ContainsKey(t *testing.T) {
	h := highlighter.Highlight{
		Key:   "SECRET_KEY",
		Left:  "abc",
		Right: "xyz",
		Diff: []highlighter.Segment{
			{Text: "abc", Removed: true},
			{Text: "xyz", Added: true},
		},
	}
	out := highlighter.Format(h)
	if !strings.Contains(out, "SECRET_KEY") {
		t.Errorf("expected key in output, got: %s", out)
	}
	if !strings.Contains(out, "[+]") {
		t.Errorf("expected added marker in output")
	}
	if !strings.Contains(out, "[-]") {
		t.Errorf("expected removed marker in output")
	}
}

func TestFormat_UnchangedSegment(t *testing.T) {
	h := highlighter.Highlight{
		Key:  "PORT",
		Diff: []highlighter.Segment{{Text: "8080"}},
	}
	out := highlighter.Format(h)
	if !strings.Contains(out, "[ ]") {
		t.Errorf("expected neutral marker, got: %s", out)
	}
}
