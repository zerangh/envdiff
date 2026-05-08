package templater_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/templater"
)

func TestGenerate_Empty(t *testing.T) {
	out := templater.Generate(map[string]string{}, templater.Options{})
	if out != "" {
		t.Errorf("expected empty string, got %q", out)
	}
}

func TestGenerate_BlankValues(t *testing.T) {
	env := map[string]string{"DB_HOST": "localhost", "APP_PORT": "8080"}
	out := templater.Generate(env, templater.Options{})

	if !strings.Contains(out, "APP_PORT=\n") {
		t.Errorf("expected blank value for APP_PORT, got:\n%s", out)
	}
	if !strings.Contains(out, "DB_HOST=\n") {
		t.Errorf("expected blank value for DB_HOST, got:\n%s", out)
	}
}

func TestGenerate_WithPlaceholder(t *testing.T) {
	env := map[string]string{"SECRET_KEY": "abc123"}
	out := templater.Generate(env, templater.Options{Placeholder: "CHANGE_ME"})

	if !strings.Contains(out, "SECRET_KEY=CHANGE_ME\n") {
		t.Errorf("expected placeholder value, got:\n%s", out)
	}
}

func TestGenerate_IncludeValues(t *testing.T) {
	env := map[string]string{"LOG_LEVEL": "debug"}
	out := templater.Generate(env, templater.Options{IncludeValues: true})

	if !strings.Contains(out, "# LOG_LEVEL=debug\n") {
		t.Errorf("expected comment with original value, got:\n%s", out)
	}
	if !strings.Contains(out, "LOG_LEVEL=\n") {
		t.Errorf("expected blank key line after comment, got:\n%s", out)
	}
}

func TestGenerate_SortedKeys(t *testing.T) {
	env := map[string]string{"ZEBRA": "1", "ALPHA": "2", "MIDDLE": "3"}
	out := templater.Generate(env, templater.Options{})

	alphaIdx := strings.Index(out, "ALPHA")
	middleIdx := strings.Index(out, "MIDDLE")
	zebraIdx := strings.Index(out, "ZEBRA")

	if !(alphaIdx < middleIdx && middleIdx < zebraIdx) {
		t.Errorf("keys not sorted: ALPHA=%d MIDDLE=%d ZEBRA=%d", alphaIdx, middleIdx, zebraIdx)
	}
}

func TestMerge_UniqueKeys(t *testing.T) {
	a := map[string]string{"FOO": "1", "BAR": "2"}
	b := map[string]string{"BAZ": "3"}
	result := templater.Merge([]map[string]string{a, b})

	if len(result) != 3 {
		t.Errorf("expected 3 keys, got %d", len(result))
	}
}

func TestMerge_FirstWins(t *testing.T) {
	a := map[string]string{"KEY": "first"}
	b := map[string]string{"KEY": "second"}
	result := templater.Merge([]map[string]string{a, b})

	if result["KEY"] != "first" {
		t.Errorf("expected 'first', got %q", result["KEY"])
	}
}
