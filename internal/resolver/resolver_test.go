package resolver_test

import (
	"testing"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/resolver"
)

func makeResult(missingRight, missingLeft []string) differ.Result {
	return differ.Result{
		MissingInRight: missingRight,
		MissingInLeft:  missingLeft,
	}
}

func TestResolve_NothingMissing(t *testing.T) {
	diff := makeResult(nil, nil)
	envs := map[string]map[string]string{
		"staging": {"DB_HOST": "localhost"},
	}
	res := resolver.Resolve(diff, envs)
	if len(res.Suggestions) != 0 {
		t.Errorf("expected no suggestions, got %d", len(res.Suggestions))
	}
	if len(res.Unresolved) != 0 {
		t.Errorf("expected no unresolved, got %d", len(res.Unresolved))
	}
}

func TestResolve_FindsSuggestion(t *testing.T) {
	diff := makeResult([]string{"DB_HOST"}, nil)
	envs := map[string]map[string]string{
		"staging": {"DB_HOST": "staging.db.local"},
	}
	res := resolver.Resolve(diff, envs)
	if len(res.Suggestions) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(res.Suggestions))
	}
	s := res.Suggestions[0]
	if s.Key != "DB_HOST" {
		t.Errorf("expected key DB_HOST, got %s", s.Key)
	}
	if s.Value != "staging.db.local" {
		t.Errorf("unexpected value: %s", s.Value)
	}
	if s.Source != "staging" {
		t.Errorf("unexpected source: %s", s.Source)
	}
}

func TestResolve_Unresolved(t *testing.T) {
	diff := makeResult([]string{"SECRET_KEY"}, nil)
	envs := map[string]map[string]string{
		"staging": {"DB_HOST": "localhost"},
	}
	res := resolver.Resolve(diff, envs)
	if len(res.Suggestions) != 0 {
		t.Errorf("expected 0 suggestions, got %d", len(res.Suggestions))
	}
	if len(res.Unresolved) != 1 || res.Unresolved[0] != "SECRET_KEY" {
		t.Errorf("expected SECRET_KEY unresolved, got %v", res.Unresolved)
	}
}

func TestResolve_SkipsEmptyValues(t *testing.T) {
	diff := makeResult([]string{"API_URL"}, nil)
	envs := map[string]map[string]string{
		"staging": {"API_URL": ""},
		"prod":    {"API_URL": "https://api.example.com"},
	}
	res := resolver.Resolve(diff, envs)
	if len(res.Suggestions) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(res.Suggestions))
	}
	if res.Suggestions[0].Value == "" {
		t.Error("expected non-empty suggestion value")
	}
}

func TestResolve_SortedOutput(t *testing.T) {
	diff := makeResult([]string{"Z_KEY", "A_KEY", "M_KEY"}, nil)
	envs := map[string]map[string]string{
		"staging": {"Z_KEY": "z", "A_KEY": "a", "M_KEY": "m"},
	}
	res := resolver.Resolve(diff, envs)
	if len(res.Suggestions) != 3 {
		t.Fatalf("expected 3 suggestions, got %d", len(res.Suggestions))
	}
	if res.Suggestions[0].Key != "A_KEY" || res.Suggestions[1].Key != "M_KEY" || res.Suggestions[2].Key != "Z_KEY" {
		t.Errorf("suggestions not sorted: %v", res.Suggestions)
	}
}

func TestResolve_MissingFromBothSides(t *testing.T) {
	diff := makeResult([]string{"LEFT_ONLY"}, []string{"RIGHT_ONLY"})
	envs := map[string]map[string]string{
		"extra": {"LEFT_ONLY": "val1", "RIGHT_ONLY": "val2"},
	}
	res := resolver.Resolve(diff, envs)
	if len(res.Suggestions) != 2 {
		t.Errorf("expected 2 suggestions, got %d", len(res.Suggestions))
	}
}
