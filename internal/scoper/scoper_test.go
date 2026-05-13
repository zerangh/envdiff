package scoper_test

import (
	"testing"

	"github.com/your-org/envdiff/internal/differ"
	"github.com/your-org/envdiff/internal/scoper"
)

func makeResult() differ.Result {
	return differ.Result{
		MissingInRight: []string{"DB_HOST", "APP_NAME"},
		MissingInLeft:  []string{"APP_VERSION"},
		Mismatched: []differ.Mismatch{
			{Key: "DB_PORT", LeftVal: "5432", RightVal: "3306"},
			{Key: "REDIS_URL", LeftVal: "localhost", RightVal: "redis:6379"},
		},
	}
}

func TestExtract_ScopesKeys(t *testing.T) {
	result := makeResult()
	scopes := scoper.Extract(result, scoper.DefaultOptions())

	names := map[string]bool{}
	for _, s := range scopes {
		names[s.Name] = true
	}

	if !names["db"] {
		t.Error("expected scope 'db'")
	}
	if !names["app"] {
		t.Error("expected scope 'app'")
	}
	if !names["redis"] {
		t.Error("expected scope 'redis'")
	}
}

func TestExtract_ScopeContainsCorrectKeys(t *testing.T) {
	result := makeResult()
	scopes := scoper.Extract(result, scoper.DefaultOptions())

	for _, s := range scopes {
		if s.Name == "db" {
			if len(s.Keys) != 2 {
				t.Errorf("expected 2 db keys, got %d", len(s.Keys))
			}
		}
	}
}

func TestExtract_FilteredResultMatchesScope(t *testing.T) {
	result := makeResult()
	scopes := scoper.Extract(result, scoper.DefaultOptions())

	for _, s := range scopes {
		if s.Name == "redis" {
			if len(s.Result.Mismatched) != 1 {
				t.Errorf("expected 1 mismatch in redis scope, got %d", len(s.Result.Mismatched))
			}
			if s.Result.Mismatched[0].Key != "REDIS_URL" {
				t.Errorf("unexpected key: %s", s.Result.Mismatched[0].Key)
			}
		}
	}
}

func TestExtract_MinKeysFilters(t *testing.T) {
	result := differ.Result{
		MissingInRight: []string{"DB_HOST"},
		Mismatched: []differ.Mismatch{
			{Key: "DB_PORT", LeftVal: "5432", RightVal: "3306"},
			{Key: "APP_KEY", LeftVal: "a", RightVal: "b"},
		},
	}
	opts := scoper.Options{MinKeys: 2}
	scopes := scoper.Extract(result, opts)

	for _, s := range scopes {
		if len(s.Keys) < 2 {
			t.Errorf("scope %q has fewer than MinKeys keys", s.Name)
		}
	}
}

func TestExtract_EmptyResult(t *testing.T) {
	scopes := scoper.Extract(differ.Result{}, scoper.DefaultOptions())
	if len(scopes) != 0 {
		t.Errorf("expected no scopes for empty result, got %d", len(scopes))
	}
}

func TestExtract_SortedOutput(t *testing.T) {
	result := makeResult()
	scopes := scoper.Extract(result, scoper.DefaultOptions())

	for i := 1; i < len(scopes); i++ {
		if scopes[i].Name < scopes[i-1].Name {
			t.Errorf("scopes not sorted: %q before %q", scopes[i-1].Name, scopes[i].Name)
		}
	}
}
