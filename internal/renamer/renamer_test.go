package renamer_test

import (
	"testing"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/renamer"
)

func makeResult(missingInRight, missingInLeft []string) differ.Result {
	return differ.Result{
		MissingInRight: missingInRight,
		MissingInLeft:  missingInLeft,
	}
}

func TestDetect_NoMissingKeys(t *testing.T) {
	result := makeResult(nil, nil)
	suggestions := renamer.Detect(result, renamer.Options{})
	if len(suggestions) != 0 {
		t.Errorf("expected no suggestions, got %d", len(suggestions))
	}
}

func TestDetect_ObviousRename(t *testing.T) {
	// DB_HOST → DATABASE_HOST should score above default threshold.
	result := makeResult(
		[]string{"DB_HOST"},
		[]string{"DATABASE_HOST"},
	)
	suggestions := renamer.Detect(result, renamer.Options{})
	if len(suggestions) == 0 {
		t.Fatal("expected at least one suggestion")
	}
	s := suggestions[0]
	if s.LeftKey != "DB_HOST" || s.RightKey != "DATABASE_HOST" {
		t.Errorf("unexpected suggestion: %+v", s)
	}
	if s.Score < 0.6 {
		t.Errorf("score too low: %.2f", s.Score)
	}
}

func TestDetect_UnrelatedKeys(t *testing.T) {
	result := makeResult(
		[]string{"FOO"},
		[]string{"XYZ_TOTALLY_DIFFERENT"},
	)
	suggestions := renamer.Detect(result, renamer.Options{})
	if len(suggestions) != 0 {
		t.Errorf("expected no suggestions for unrelated keys, got %d", len(suggestions))
	}
}

func TestDetect_CustomMinScore(t *testing.T) {
	result := makeResult(
		[]string{"API_KEY"},
		[]string{"API_TOKEN"},
	)
	// Very high threshold — should suppress borderline matches.
	suggestions := renamer.Detect(result, renamer.Options{MinScore: 0.99})
	if len(suggestions) != 0 {
		t.Errorf("expected no suggestions at 0.99 threshold, got %d", len(suggestions))
	}

	// Permissive threshold — should surface the match.
	suggestions = renamer.Detect(result, renamer.Options{MinScore: 0.3})
	if len(suggestions) == 0 {
		t.Error("expected a suggestion at 0.3 threshold")
	}
}

func TestDetect_SortedByScore(t *testing.T) {
	result := makeResult(
		[]string{"DB_HOST", "DB_PORT"},
		[]string{"DATABASE_HOST", "DATABASE_PORT"},
	)
	suggestions := renamer.Detect(result, renamer.Options{MinScore: 0.4})
	for i := 1; i < len(suggestions); i++ {
		if suggestions[i-1].Score < suggestions[i].Score {
			t.Errorf("suggestions not sorted by descending score at index %d", i)
		}
	}
}

func TestDetect_ExactMatchScoresOne(t *testing.T) {
	// If the same key appears on both sides (unusual but edge-case worthy).
	result := makeResult(
		[]string{"SAME_KEY"},
		[]string{"SAME_KEY"},
	)
	suggestions := renamer.Detect(result, renamer.Options{MinScore: 0.0})
	if len(suggestions) == 0 {
		t.Fatal("expected a suggestion for identical keys")
	}
	if suggestions[0].Score != 1.0 {
		t.Errorf("expected score 1.0 for identical keys, got %.2f", suggestions[0].Score)
	}
}
