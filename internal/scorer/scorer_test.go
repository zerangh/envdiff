package scorer_test

import (
	"testing"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/scorer"
)

func makeResult(missingRight, missingLeft []string, mismatched []differ.Mismatch) differ.Result {
	return differ.Result{
		MissingInRight: missingRight,
		MissingInLeft:  missingLeft,
		Matched:        []string{"COMMON_KEY"},
		Mismatched:     mismatched,
	}
}

func TestCalculate_PerfectScore(t *testing.T) {
	result := makeResult(nil, nil, nil)
	s := scorer.Calculate(result, scorer.Options{})
	if s.Total != 100 {
		t.Errorf("expected 100, got %d", s.Total)
	}
	if s.MissingPenalty != 0 || s.MismatchPenalty != 0 {
		t.Errorf("expected zero penalties, got missing=%d mismatch=%d", s.MissingPenalty, s.MismatchPenalty)
	}
}

func TestCalculate_MissingKeys(t *testing.T) {
	result := makeResult([]string{"KEY_A", "KEY_B"}, nil, nil)
	s := scorer.Calculate(result, scorer.Options{})
	// 2 missing * 10 = 20 penalty => score 80
	if s.Total != 80 {
		t.Errorf("expected 80, got %d", s.Total)
	}
	if s.MissingPenalty != 20 {
		t.Errorf("expected missing penalty 20, got %d", s.MissingPenalty)
	}
}

func TestCalculate_MissingKeysFromBothSides(t *testing.T) {
	// Keys missing in right AND left should both contribute to the missing penalty.
	result := makeResult([]string{"KEY_A"}, []string{"KEY_B"}, nil)
	s := scorer.Calculate(result, scorer.Options{})
	// 2 missing total * 10 = 20 penalty => score 80
	if s.Total != 80 {
		t.Errorf("expected 80, got %d", s.Total)
	}
	if s.MissingPenalty != 20 {
		t.Errorf("expected missing penalty 20, got %d", s.MissingPenalty)
	}
}

func TestCalculate_MismatchedValues(t *testing.T) {
	mismatches := []differ.Mismatch{
		{Key: "DB_HOST", LeftVal: "localhost", RightVal: "prod.db"},
	}
	result := makeResult(nil, nil, mismatches)
	s := scorer.Calculate(result, scorer.Options{})
	// 1 mismatch * 5 = 5 penalty => score 95
	if s.Total != 95 {
		t.Errorf("expected 95, got %d", s.Total)
	}
	if s.MismatchPenalty != 5 {
		t.Errorf("expected mismatch penalty 5, got %d", s.MismatchPenalty)
	}
}

func TestCalculate_FloorAtZero(t *testing.T) {
	mismatches := make([]differ.Mismatch, 20)
	for i := range mismatches {
		mismatches[i] = differ.Mismatch{Key: "KEY", LeftVal: "a", RightVal: "b"}
	}
	result := makeResult([]string{"K1", "K2", "K3"}, []string{"K4"}, mismatches)
	s := scorer.Calculate(result, scorer.Options{})
	if s.Total < 0 {
		t.Errorf("score should not be negative, got %d", s.Total)
	}
	if s.Total != 0 {
		t.Errorf("expected 0, got %d", s.Total)
	}
}

func TestCalculate_CustomWeights(t *testing.T) {
	result := makeResult([]string{"KEY_A"}, nil, nil)
	s := scorer.Calculate(result, scorer.Options{MissingWeight: 20, MismatchWeight: 2})
	// 1 missing * 20 = 20 => score 80
	if s.Total != 80 {
		t.Errorf("expected 80, got %d", s.Total)
	}
}

func TestGrade(t *testing.T) {
	cases := []struct {
		score int
		want  string
	}{
		{100, "A"}, {90, "A"}, {89, "B"}, {75, "B"},
		{74, "C"}, {60, "C"}, {59, "D"}, {40, "D"},
		{39, "F"}, {0, "F"},
	}
	for _, c := range cases {
		got := scorer.Grade(c.score)
		if got != c.want {
			t.Errorf("Grade(%d) = %q, want %q", c.score, got, c.want)
		}
	}
}
