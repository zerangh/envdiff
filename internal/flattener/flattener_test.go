package flattener_test

import (
	"testing"

	"github.com/user/envdiff/internal/flattener"
)

func named(name string, kv map[string]string) flattener.NamedEnv {
	return flattener.NamedEnv{Name: name, Env: kv}
}

func TestFlatten_EmptyInput(t *testing.T) {
	res := flattener.Flatten(nil, flattener.DefaultOptions())
	if len(res.Env) != 0 {
		t.Errorf("expected empty env, got %v", res.Env)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %v", res.Conflicts)
	}
}

func TestFlatten_SingleEnv(t *testing.T) {
	envs := []flattener.NamedEnv{
		named("dev", map[string]string{"HOST": "localhost", "PORT": "5432"}),
	}
	res := flattener.Flatten(envs, flattener.DefaultOptions())
	if res.Env["HOST"] != "localhost" {
		t.Errorf("expected localhost, got %s", res.Env["HOST"])
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts")
	}
}

func TestFlatten_StrategyFirst_KeepsFirst(t *testing.T) {
	envs := []flattener.NamedEnv{
		named("dev", map[string]string{"DB": "dev-db"}),
		named("prod", map[string]string{"DB": "prod-db"}),
	}
	opts := flattener.Options{Strategy: flattener.StrategyFirst}
	res := flattener.Flatten(envs, opts)
	if res.Env["DB"] != "dev-db" {
		t.Errorf("expected dev-db, got %s", res.Env["DB"])
	}
	if len(res.Conflicts) != 1 || res.Conflicts[0] != "DB" {
		t.Errorf("expected DB in conflicts, got %v", res.Conflicts)
	}
}

func TestFlatten_StrategyLast_KeepsLast(t *testing.T) {
	envs := []flattener.NamedEnv{
		named("dev", map[string]string{"DB": "dev-db"}),
		named("prod", map[string]string{"DB": "prod-db"}),
	}
	opts := flattener.Options{Strategy: flattener.StrategyLast}
	res := flattener.Flatten(envs, opts)
	if res.Env["DB"] != "prod-db" {
		t.Errorf("expected prod-db, got %s", res.Env["DB"])
	}
}

func TestFlatten_ConflictsSorted(t *testing.T) {
	envs := []flattener.NamedEnv{
		named("a", map[string]string{"Z": "1", "A": "1", "M": "1"}),
		named("b", map[string]string{"Z": "2", "A": "2", "M": "2"}),
	}
	res := flattener.Flatten(envs, flattener.DefaultOptions())
	if len(res.Conflicts) != 3 {
		t.Fatalf("expected 3 conflicts, got %d", len(res.Conflicts))
	}
	if res.Conflicts[0] != "A" || res.Conflicts[1] != "M" || res.Conflicts[2] != "Z" {
		t.Errorf("conflicts not sorted: %v", res.Conflicts)
	}
}

func TestFlatten_OriginsTracked(t *testing.T) {
	envs := []flattener.NamedEnv{
		named("dev", map[string]string{"HOST": "localhost"}),
		named("prod", map[string]string{"PORT": "443"}),
	}
	res := flattener.Flatten(envs, flattener.DefaultOptions())
	if len(res.Origins) != 2 {
		t.Fatalf("expected 2 origins, got %d", len(res.Origins))
	}
	// origins are sorted by key: HOST < PORT
	if res.Origins[0].Key != "HOST" || res.Origins[0].EnvName != "dev" {
		t.Errorf("unexpected origin[0]: %+v", res.Origins[0])
	}
	if res.Origins[1].Key != "PORT" || res.Origins[1].EnvName != "prod" {
		t.Errorf("unexpected origin[1]: %+v", res.Origins[1])
	}
}
