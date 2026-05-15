package stripper_test

import (
	"testing"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/stripper"
)

func makeResult() differ.Result {
	return differ.Result{
		MissingInRight: []string{"DB_HOST", "INTERNAL_TOKEN", "APP_PORT"},
		MissingInLeft:  []string{"REDIS_URL", "SECRET_KEY"},
		Mismatched: []differ.Mismatch{
			{Key: "DB_PASS", Left: "abc", Right: "xyz"},
			{Key: "LOG_LEVEL", Left: "info", Right: "debug"},
		},
	}
}

func TestStripEnv_ExactKey(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2", "C": "3"}
	out := stripper.StripEnv(env, stripper.Options{Keys: []string{"B"}})
	if _, ok := out["B"]; ok {
		t.Error("expected B to be stripped")
	}
	if out["A"] != "1" || out["C"] != "3" {
		t.Error("expected A and C to remain")
	}
}

func TestStripEnv_Prefix(t *testing.T) {
	env := map[string]string{"DB_HOST": "h", "DB_PASS": "p", "APP_PORT": "80"}
	out := stripper.StripEnv(env, stripper.Options{Prefixes: []string{"DB_"}})
	if _, ok := out["DB_HOST"]; ok {
		t.Error("expected DB_HOST stripped")
	}
	if _, ok := out["DB_PASS"]; ok {
		t.Error("expected DB_PASS stripped")
	}
	if out["APP_PORT"] != "80" {
		t.Error("expected APP_PORT to remain")
	}
}

func TestStripEnv_Suffix(t *testing.T) {
	env := map[string]string{"API_TOKEN": "t", "SECRET_KEY": "s", "HOST": "h"}
	out := stripper.StripEnv(env, stripper.Options{Suffixes: []string{"_TOKEN", "_KEY"}})
	if _, ok := out["API_TOKEN"]; ok {
		t.Error("expected API_TOKEN stripped")
	}
	if _, ok := out["SECRET_KEY"]; ok {
		t.Error("expected SECRET_KEY stripped")
	}
	if out["HOST"] != "h" {
		t.Error("expected HOST to remain")
	}
}

func TestStripResult_ExactKey(t *testing.T) {
	r := makeResult()
	out := stripper.StripResult(r, stripper.Options{Keys: []string{"DB_HOST", "SECRET_KEY", "DB_PASS"}})
	for _, k := range out.MissingInRight {
		if k == "DB_HOST" {
			t.Error("DB_HOST should be stripped from MissingInRight")
		}
	}
	for _, k := range out.MissingInLeft {
		if k == "SECRET_KEY" {
			t.Error("SECRET_KEY should be stripped from MissingInLeft")
		}
	}
	for _, m := range out.Mismatched {
		if m.Key == "DB_PASS" {
			t.Error("DB_PASS should be stripped from Mismatched")
		}
	}
}

func TestStripResult_Prefix(t *testing.T) {
	r := makeResult()
	out := stripper.StripResult(r, stripper.Options{Prefixes: []string{"INTERNAL_"}})
	for _, k := range out.MissingInRight {
		if k == "INTERNAL_TOKEN" {
			t.Error("INTERNAL_TOKEN should be stripped")
		}
	}
	if len(out.MissingInRight) != 2 {
		t.Errorf("expected 2 remaining MissingInRight, got %d", len(out.MissingInRight))
	}
}

func TestStripResult_NoOptions(t *testing.T) {
	r := makeResult()
	out := stripper.StripResult(r, stripper.Options{})
	if len(out.MissingInRight) != len(r.MissingInRight) {
		t.Error("expected MissingInRight unchanged")
	}
	if len(out.MissingInLeft) != len(r.MissingInLeft) {
		t.Error("expected MissingInLeft unchanged")
	}
	if len(out.Mismatched) != len(r.Mismatched) {
		t.Error("expected Mismatched unchanged")
	}
}
