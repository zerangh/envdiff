package trimmer_test

import (
	"testing"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/trimmer"
)

func makeResult() differ.Result {
	return differ.Result{
		MissingInRight: []string{"CI_TOKEN", "APP_KEY", "DEBUG"},
		MissingInLeft:  []string{"DEPLOY_HOST", "PORT"},
		Mismatched: []differ.Mismatch{
			{Key: "CI_BUILD_ID", Left: "1", Right: "2"},
			{Key: "DATABASE_URL", Left: "a", Right: "b"},
		},
	}
}

func TestTrimEnv_ExactKeys(t *testing.T) {
	env := map[string]string{"FOO": "1", "BAR": "2", "BAZ": "3"}
	out := trimmer.TrimEnv(env, trimmer.Options{Keys: []string{"BAR"}})
	if _, ok := out["BAR"]; ok {
		t.Error("expected BAR to be removed")
	}
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
}

func TestTrimEnv_Prefix(t *testing.T) {
	env := map[string]string{"CI_TOKEN": "x", "CI_JOB": "y", "APP_KEY": "z"}
	out := trimmer.TrimEnv(env, trimmer.Options{Prefixes: []string{"CI_"}})
	if len(out) != 1 {
		t.Errorf("expected 1 key, got %d", len(out))
	}
	if _, ok := out["APP_KEY"]; !ok {
		t.Error("expected APP_KEY to remain")
	}
}

func TestTrimEnv_NoOptions(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	out := trimmer.TrimEnv(env, trimmer.Options{})
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
}

func TestTrimResult_ExactKey(t *testing.T) {
	r := makeResult()
	out := trimmer.TrimResult(r, trimmer.Options{Keys: []string{"APP_KEY", "PORT"}})
	for _, k := range out.MissingInRight {
		if k == "APP_KEY" {
			t.Error("APP_KEY should have been trimmed from MissingInRight")
		}
	}
	for _, k := range out.MissingInLeft {
		if k == "PORT" {
			t.Error("PORT should have been trimmed from MissingInLeft")
		}
	}
}

func TestTrimResult_PrefixRemovesMismatched(t *testing.T) {
	r := makeResult()
	out := trimmer.TrimResult(r, trimmer.Options{Prefixes: []string{"CI_"}})
	for _, m := range out.Mismatched {
		if m.Key == "CI_BUILD_ID" {
			t.Error("CI_BUILD_ID should have been trimmed from Mismatched")
		}
	}
	if len(out.Mismatched) != 1 || out.Mismatched[0].Key != "DATABASE_URL" {
		t.Errorf("expected only DATABASE_URL in Mismatched, got %v", out.Mismatched)
	}
}

func TestTrimResult_DoesNotMutateOriginal(t *testing.T) {
	r := makeResult()
	orig := len(r.MissingInRight)
	_ = trimmer.TrimResult(r, trimmer.Options{Keys: []string{"CI_TOKEN"}})
	if len(r.MissingInRight) != orig {
		t.Error("original result was mutated")
	}
}
