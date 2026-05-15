package pinner_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/pinner"
)

func TestCheckEnv_AllPresent(t *testing.T) {
	env := map[string]string{"DB_URL": "postgres://localhost", "API_KEY": "secret"}
	r := pinner.CheckEnv("prod", env, []string{"DB_URL", "API_KEY"})
	if !r.OK() {
		t.Fatalf("expected OK, got missing=%v empty=%v", r.Missing, r.Empty)
	}
}

func TestCheckEnv_MissingKey(t *testing.T) {
	env := map[string]string{"DB_URL": "postgres://localhost"}
	r := pinner.CheckEnv("staging", env, []string{"DB_URL", "API_KEY"})
	if r.OK() {
		t.Fatal("expected not OK")
	}
	if len(r.Missing) != 1 || r.Missing[0].Key != "API_KEY" {
		t.Fatalf("unexpected missing: %v", r.Missing)
	}
}

func TestCheckEnv_EmptyValue(t *testing.T) {
	env := map[string]string{"DB_URL": "postgres://localhost", "API_KEY": ""}
	r := pinner.CheckEnv("dev", env, []string{"DB_URL", "API_KEY"})
	if r.OK() {
		t.Fatal("expected not OK due to empty value")
	}
	if len(r.Empty) != 1 || r.Empty[0].Key != "API_KEY" {
		t.Fatalf("unexpected empty: %v", r.Empty)
	}
}

func TestCheckEnv_EnvLabelSet(t *testing.T) {
	env := map[string]string{}
	r := pinner.CheckEnv("prod", env, []string{"SECRET"})
	if len(r.Missing) == 0 {
		t.Fatal("expected missing entry")
	}
	if r.Missing[0].Env != "prod" {
		t.Fatalf("expected env label 'prod', got %q", r.Missing[0].Env)
	}
}

func TestCheckAll_MultipleEnvs(t *testing.T) {
	envs := map[string]map[string]string{
		"prod":    {"DB_URL": "postgres://prod", "API_KEY": "key"},
		"staging": {"DB_URL": "postgres://staging"},
	}
	r := pinner.CheckAll(envs, []string{"DB_URL", "API_KEY"})
	if r.OK() {
		t.Fatal("expected not OK")
	}
	if len(r.Missing) != 1 || r.Missing[0].Key != "API_KEY" || r.Missing[0].Env != "staging" {
		t.Fatalf("unexpected missing: %v", r.Missing)
	}
}

func TestReport_Format_Clean(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	r := pinner.CheckEnv("test", env, []string{"A", "B"})
	out := r.Format()
	if !strings.Contains(out, "all 2 pinned keys") {
		t.Fatalf("unexpected format output: %q", out)
	}
}

func TestReport_Format_WithIssues(t *testing.T) {
	env := map[string]string{"A": ""}
	r := pinner.CheckEnv("test", env, []string{"A", "B"})
	out := r.Format()
	if !strings.Contains(out, "MISSING") || !strings.Contains(out, "EMPTY") {
		t.Fatalf("expected MISSING and EMPTY in output: %q", out)
	}
}
