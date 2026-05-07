package filter_test

import (
	"testing"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/filter"
)

func baseResult() differ.DiffResult {
	return differ.DiffResult{
		MissingInRight: []string{"APP_HOST", "DB_HOST"},
		MissingInLeft:  []string{"REDIS_URL"},
		Mismatched: []differ.Mismatch{
			{Key: "APP_PORT", LeftVal: "8080", RightVal: "9090"},
			{Key: "DB_PORT", LeftVal: "5432", RightVal: "5433"},
		},
	}
}

func TestApply_NoOptions(t *testing.T) {
	r := filter.Apply(baseResult(), filter.Options{})
	if len(r.MissingInRight) != 2 {
		t.Errorf("expected 2 MissingInRight, got %d", len(r.MissingInRight))
	}
	if len(r.MissingInLeft) != 1 {
		t.Errorf("expected 1 MissingInLeft, got %d", len(r.MissingInLeft))
	}
	if len(r.Mismatched) != 2 {
		t.Errorf("expected 2 Mismatched, got %d", len(r.Mismatched))
	}
}

func TestApply_PrefixFilter(t *testing.T) {
	r := filter.Apply(baseResult(), filter.Options{Prefix: "DB_"})
	if len(r.MissingInRight) != 1 || r.MissingInRight[0] != "DB_HOST" {
		t.Errorf("unexpected MissingInRight: %v", r.MissingInRight)
	}
	if len(r.MissingInLeft) != 0 {
		t.Errorf("expected no MissingInLeft, got %v", r.MissingInLeft)
	}
	if len(r.Mismatched) != 1 || r.Mismatched[0].Key != "DB_PORT" {
		t.Errorf("unexpected Mismatched: %v", r.Mismatched)
	}
}

func TestApply_ExcludeKeys(t *testing.T) {
	opts := filter.Options{Exclude: []string{"APP_HOST", "APP_PORT"}}
	r := filter.Apply(baseResult(), opts)
	for _, k := range r.MissingInRight {
		if k == "APP_HOST" {
			t.Error("APP_HOST should have been excluded from MissingInRight")
		}
	}
	for _, m := range r.Mismatched {
		if m.Key == "APP_PORT" {
			t.Error("APP_PORT should have been excluded from Mismatched")
		}
	}
}

func TestApply_OnlyMissing(t *testing.T) {
	r := filter.Apply(baseResult(), filter.Options{OnlyMissing: true})
	if len(r.Mismatched) != 0 {
		t.Errorf("expected no Mismatched when OnlyMissing=true, got %d", len(r.Mismatched))
	}
	if len(r.MissingInRight)+len(r.MissingInLeft) == 0 {
		t.Error("expected missing keys to be present")
	}
}

func TestApply_OnlyMismatched(t *testing.T) {
	r := filter.Apply(baseResult(), filter.Options{OnlyMismatched: true})
	if len(r.MissingInRight) != 0 || len(r.MissingInLeft) != 0 {
		t.Error("expected no missing keys when OnlyMismatched=true")
	}
	if len(r.Mismatched) != 2 {
		t.Errorf("expected 2 mismatched entries, got %d", len(r.Mismatched))
	}
}
