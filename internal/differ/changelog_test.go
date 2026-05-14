package differ

import (
	"testing"
)

func makeChangelogResult() Result {
	return Result{
		LeftFile:  "a.env",
		RightFile: "b.env",
		LeftEnv:   map[string]string{"FOO": "bar", "OLD": "val"},
		RightEnv:  map[string]string{"FOO": "baz", "NEW": "fresh"},
		MissingInRight: []string{"OLD"},
		MissingInLeft:  []string{"NEW"},
		Mismatched: []MismatchedKey{
			{Key: "FOO", LeftValue: "bar", RightValue: "baz"},
		},
	}
}

func TestBuildChangelog_EntryCount(t *testing.T) {
	r := makeChangelogResult()
	cl := BuildChangelog(r)
	if len(cl.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(cl.Entries))
	}
}

func TestBuildChangelog_FileNames(t *testing.T) {
	r := makeChangelogResult()
	cl := BuildChangelog(r)
	if cl.LeftFile != "a.env" {
		t.Errorf("expected left file a.env, got %s", cl.LeftFile)
	}
	if cl.RightFile != "b.env" {
		t.Errorf("expected right file b.env, got %s", cl.RightFile)
	}
}

func TestBuildChangelog_RemovedEntry(t *testing.T) {
	r := makeChangelogResult()
	cl := BuildChangelog(r)
	var found *ChangeEntry
	for i := range cl.Entries {
		if cl.Entries[i].Key == "OLD" {
			found = &cl.Entries[i]
			break
		}
	}
	if found == nil {
		t.Fatal("expected entry for OLD key")
	}
	if found.Kind != ChangeRemoved {
		t.Errorf("expected kind removed, got %s", found.Kind)
	}
	if found.LeftValue != "val" {
		t.Errorf("expected left value 'val', got %s", found.LeftValue)
	}
}

func TestBuildChangelog_AddedEntry(t *testing.T) {
	r := makeChangelogResult()
	cl := BuildChangelog(r)
	var found *ChangeEntry
	for i := range cl.Entries {
		if cl.Entries[i].Key == "NEW" {
			found = &cl.Entries[i]
			break
		}
	}
	if found == nil {
		t.Fatal("expected entry for NEW key")
	}
	if found.Kind != ChangeAdded {
		t.Errorf("expected kind added, got %s", found.Kind)
	}
	if found.RightValue != "fresh" {
		t.Errorf("expected right value 'fresh', got %s", found.RightValue)
	}
}

func TestBuildChangelog_MismatchEntry(t *testing.T) {
	r := makeChangelogResult()
	cl := BuildChangelog(r)
	var found *ChangeEntry
	for i := range cl.Entries {
		if cl.Entries[i].Key == "FOO" {
			found = &cl.Entries[i]
			break
		}
	}
	if found == nil {
		t.Fatal("expected entry for FOO key")
	}
	if found.Kind != ChangeMismatch {
		t.Errorf("expected kind mismatch, got %s", found.Kind)
	}
	if found.LeftValue != "bar" || found.RightValue != "baz" {
		t.Errorf("unexpected values: left=%s right=%s", found.LeftValue, found.RightValue)
	}
}

func TestBuildChangelog_EmptyResult(t *testing.T) {
	r := Result{
		LeftEnv:  map[string]string{},
		RightEnv: map[string]string{},
	}
	cl := BuildChangelog(r)
	if len(cl.Entries) != 0 {
		t.Errorf("expected 0 entries for clean result, got %d", len(cl.Entries))
	}
	if cl.GeneratedAt.IsZero() {
		t.Error("expected GeneratedAt to be set")
	}
}
