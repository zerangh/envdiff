package differ

import (
	"testing"

	"github.com/your-org/envdiff/internal/parser"
)

func TestIntersect_EmptyInput(t *testing.T) {
	result := Intersect(nil)
	if len(result.CommonKeys) != 0 {
		t.Errorf("expected no common keys, got %v", result.CommonKeys)
	}
}

func TestIntersect_SingleEnv(t *testing.T) {
	envs := map[string]parser.Env{
		"prod": {"FOO": "bar", "BAZ": "qux"},
	}
	result := Intersect(envs)
	if len(result.CommonKeys) != 2 {
		t.Fatalf("expected 2 common keys, got %d", len(result.CommonKeys))
	}
	if len(result.Consistent) != 2 {
		t.Errorf("expected 2 consistent keys with single env, got %d", len(result.Consistent))
	}
	if len(result.Divergent) != 0 {
		t.Errorf("expected 0 divergent keys, got %d", len(result.Divergent))
	}
}

func TestIntersect_AllCommonAllSame(t *testing.T) {
	envs := map[string]parser.Env{
		"dev":  {"HOST": "localhost", "PORT": "8080"},
		"prod": {"HOST": "localhost", "PORT": "8080"},
	}
	result := Intersect(envs)
	if len(result.CommonKeys) != 2 {
		t.Fatalf("expected 2 common keys, got %d", len(result.CommonKeys))
	}
	if len(result.Consistent) != 2 {
		t.Errorf("expected 2 consistent, got %d", len(result.Consistent))
	}
	if len(result.Divergent) != 0 {
		t.Errorf("expected 0 divergent, got %d", len(result.Divergent))
	}
}

func TestIntersect_SomeCommonDivergent(t *testing.T) {
	envs := map[string]parser.Env{
		"dev":  {"HOST": "localhost", "PORT": "3000", "ONLY_DEV": "1"},
		"prod": {"HOST": "example.com", "PORT": "3000"},
	}
	result := Intersect(envs)

	if len(result.CommonKeys) != 2 {
		t.Fatalf("expected 2 common keys, got %v", result.CommonKeys)
	}
	if len(result.Consistent) != 1 || result.Consistent[0] != "PORT" {
		t.Errorf("expected PORT as consistent, got %v", result.Consistent)
	}
	if len(result.Divergent) != 1 || result.Divergent[0] != "HOST" {
		t.Errorf("expected HOST as divergent, got %v", result.Divergent)
	}
}

func TestIntersect_ValueMapPopulated(t *testing.T) {
	envs := map[string]parser.Env{
		"dev":  {"DB": "dev-db"},
		"prod": {"DB": "prod-db"},
	}
	result := Intersect(envs)

	vals, ok := result.ValueMap["DB"]
	if !ok {
		t.Fatal("expected DB in ValueMap")
	}
	if vals["dev"] != "dev-db" {
		t.Errorf("expected dev-db for dev, got %s", vals["dev"])
	}
	if vals["prod"] != "prod-db" {
		t.Errorf("expected prod-db for prod, got %s", vals["prod"])
	}
}

func TestIntersect_SortedOutput(t *testing.T) {
	envs := map[string]parser.Env{
		"a": {"Z": "1", "A": "1", "M": "1"},
		"b": {"Z": "2", "A": "1", "M": "1"},
	}
	result := Intersect(envs)

	for i := 1; i < len(result.CommonKeys); i++ {
		if result.CommonKeys[i] < result.CommonKeys[i-1] {
			t.Errorf("CommonKeys not sorted: %v", result.CommonKeys)
		}
	}
	for i := 1; i < len(result.Consistent); i++ {
		if result.Consistent[i] < result.Consistent[i-1] {
			t.Errorf("Consistent not sorted: %v", result.Consistent)
		}
	}
}
