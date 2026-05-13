package grouper_test

import (
	"testing"

	"github.com/yourusername/envdiff/internal/differ"
	"github.com/yourusername/envdiff/internal/grouper"
)

func makeResult(missingRight, missingLeft []string, mismatched []differ.Mismatch) differ.Result {
	return differ.Result{
		MissingInRight: missingRight,
		MissingInLeft:  missingLeft,
		Mismatched:     mismatched,
	}
}

func TestByPrefix_EmptyResult(t *testing.T) {
	d := makeResult(nil, nil, nil)
	r := grouper.ByPrefix(d, "_")
	if len(r.Groups) != 0 {
		t.Errorf("expected 0 groups, got %d", len(r.Groups))
	}
}

func TestByPrefix_SinglePrefix(t *testing.T) {
	d := makeResult([]string{"DB_HOST", "DB_PORT"}, nil, nil)
	r := grouper.ByPrefix(d, "_")
	if len(r.Groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(r.Groups))
	}
	if r.Groups[0].Prefix != "DB" {
		t.Errorf("expected prefix DB, got %s", r.Groups[0].Prefix)
	}
	if len(r.Groups[0].Missing) != 2 {
		t.Errorf("expected 2 missing keys in group, got %d", len(r.Groups[0].Missing))
	}
}

func TestByPrefix_MultiplePrefix(t *testing.T) {
	d := makeResult(
		[]string{"DB_HOST"},
		[]string{"REDIS_URL"},
		[]differ.Mismatch{{Key: "AWS_REGION", LeftVal: "us-east-1", RightVal: "eu-west-1"}},
	)
	r := grouper.ByPrefix(d, "_")
	if len(r.Groups) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(r.Groups))
	}
	prefixes := map[string]bool{}
	for _, g := range r.Groups {
		prefixes[g.Prefix] = true
	}
	for _, want := range []string{"AWS", "DB", "REDIS"} {
		if !prefixes[want] {
			t.Errorf("expected prefix %s in groups", want)
		}
	}
}

func TestByPrefix_MismatchAnnotated(t *testing.T) {
	d := makeResult(nil, nil, []differ.Mismatch{
		{Key: "APP_SECRET", LeftVal: "abc", RightVal: "xyz"},
	})
	r := grouper.ByPrefix(d, "_")
	if len(r.Groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(r.Groups))
	}
	if len(r.Groups[0].Mismatch) != 1 || r.Groups[0].Mismatch[0] != "APP_SECRET" {
		t.Errorf("expected APP_SECRET in Mismatch, got %v", r.Groups[0].Mismatch)
	}
}

func TestByPrefix_UngroupedKeys(t *testing.T) {
	d := makeResult([]string{"NOPREFIXKEY"}, nil, nil)
	r := grouper.ByPrefix(d, "_")
	if len(r.Groups) != 0 {
		t.Errorf("expected 0 named groups, got %d", len(r.Groups))
	}
	if r.Ungrouped.Prefix != "(none)" {
		t.Errorf("expected ungrouped prefix '(none)', got %q", r.Ungrouped.Prefix)
	}
}

func TestByPrefix_DefaultDelimiter(t *testing.T) {
	d := makeResult([]string{"DB_HOST"}, nil, nil)
	r := grouper.ByPrefix(d, "")
	if len(r.Groups) != 1 || r.Groups[0].Prefix != "DB" {
		t.Errorf("expected group DB with default delimiter, got %+v", r.Groups)
	}
}
