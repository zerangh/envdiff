package differ

import (
	"os"
	"path/filepath"
	"testing"
)

func writeUnionTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestUnion_EmptyInput(t *testing.T) {
	result := Union(map[string]map[string]string{})
	if len(result.AllKeys) != 0 {
		t.Errorf("expected no keys, got %v", result.AllKeys)
	}
}

func TestUnion_SingleEnv(t *testing.T) {
	envs := map[string]map[string]string{
		"prod": {"FOO": "bar", "BAZ": "qux"},
	}
	result := Union(envs)
	if len(result.AllKeys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(result.AllKeys))
	}
	if result.AllKeys[0] != "BAZ" || result.AllKeys[1] != "FOO" {
		t.Errorf("keys not sorted: %v", result.AllKeys)
	}
}

func TestUnion_MergesAllKeys(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"A": "1", "B": "2"},
		"prod": {"B": "2", "C": "3"},
	}
	result := Union(envs)
	if len(result.AllKeys) != 3 {
		t.Fatalf("expected 3 keys, got %d: %v", len(result.AllKeys), result.AllKeys)
	}
}

func TestUnion_ValuesPerEnv(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"A": "dev-a"},
		"prod": {"A": "prod-a", "B": "prod-b"},
	}
	result := Union(envs)

	if result.Values["A"]["dev"] != "dev-a" {
		t.Errorf("expected dev-a for A in dev")
	}
	if result.Values["A"]["prod"] != "prod-a" {
		t.Errorf("expected prod-a for A in prod")
	}
	if _, ok := result.Values["B"]["dev"]; ok {
		t.Errorf("B should be absent in dev")
	}
}

func TestUnion_EnvsSorted(t *testing.T) {
	envs := map[string]map[string]string{
		"staging": {"X": "1"},
		"dev":     {"X": "2"},
		"prod":    {"X": "3"},
	}
	result := Union(envs)
	expected := []string{"dev", "prod", "staging"}
	for i, e := range expected {
		if result.Envs[i] != e {
			t.Errorf("expected env[%d]=%s, got %s", i, e, result.Envs[i])
		}
	}
}

func TestUnionFiles_ParsesFiles(t *testing.T) {
	dir := t.TempDir()
	p1 := filepath.Join(dir, "a.env")
	p2 := filepath.Join(dir, "b.env")
	os.WriteFile(p1, []byte("FOO=1\nBAR=2\n"), 0644)
	os.WriteFile(p2, []byte("BAR=2\nBAZ=3\n"), 0644)

	result, err := UnionFiles([]string{p1, p2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.AllKeys) != 3 {
		t.Errorf("expected 3 union keys, got %d: %v", len(result.AllKeys), result.AllKeys)
	}
}

func TestUnionFiles_InvalidFile(t *testing.T) {
	_, err := UnionFiles([]string{"/nonexistent/path.env"})
	if err == nil {
		t.Error("expected error for missing file")
	}
}
