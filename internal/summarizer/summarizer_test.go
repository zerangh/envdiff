package summarizer_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/summarizer"
)

func makeResult(missingLeft, missingRight []string, mismatched []differ.Mismatch, common []string) differ.Result {
	return differ.Result{
		MissingInLeft:  missingLeft,
		MissingInRight: missingRight,
		Mismatched:     mismatched,
		Common:         common,
	}
}

func TestSummarize_Clean(t *testing.T) {
	result := makeResult(nil, nil, nil, []string{"KEY_A", "KEY_B"})
	s := summarizer.Summarize(result)

	if s.TotalKeys != 2 {
		t.Errorf("expected TotalKeys=2, got %d", s.TotalKeys)
	}
	if s.MissingInLeft != 0 || s.MissingInRight != 0 || s.Mismatched != 0 {
		t.Error("expected zero diffs for clean result")
	}
	if !s.Healthy {
		t.Error("expected Healthy=true for clean result")
	}
	if s.Grade == "" {
		t.Error("expected non-empty grade")
	}
}

func TestSummarize_WithMissing(t *testing.T) {
	result := makeResult(
		[]string{"MISSING_LEFT"},
		[]string{"MISSING_RIGHT"},
		nil,
		[]string{"SHARED"},
	)
	s := summarizer.Summarize(result)

	if s.MissingInLeft != 1 {
		t.Errorf("expected MissingInLeft=1, got %d", s.MissingInLeft)
	}
	if s.MissingInRight != 1 {
		t.Errorf("expected MissingInRight=1, got %d", s.MissingInRight)
	}
	if s.TotalKeys != 3 {
		t.Errorf("expected TotalKeys=3, got %d", s.TotalKeys)
	}
}

func TestSummarize_WithMismatched(t *testing.T) {
	mismatched := []differ.Mismatch{
		{Key: "DB_HOST", LeftValue: "localhost", RightValue: "prod.db"},
	}
	result := makeResult(nil, nil, mismatched, []string{"OTHER"})
	s := summarizer.Summarize(result)

	if s.Mismatched != 1 {
		t.Errorf("expected Mismatched=1, got %d", s.Mismatched)
	}
}

func TestSummary_Format(t *testing.T) {
	result := makeResult(
		nil,
		[]string{"MISSING_KEY"},
		nil,
		[]string{"PRESENT"},
	)
	s := summarizer.Summarize(result)
	output := s.Format()

	for _, want := range []string{"Summary", "Total keys", "Score", "Status", "Grade"} {
		if !strings.Contains(output, want) {
			t.Errorf("Format() missing expected field %q\nGot:\n%s", want, output)
		}
	}
}

func TestSummary_HealthyFlag(t *testing.T) {
	// Many missing keys should push score below 80
	missing := make([]string, 15)
	for i := range missing {
		missing[i] = fmt.Sprintf("KEY_%d", i)
	}
	result := makeResult(nil, missing, nil, nil)
	s := summarizer.Summarize(result)

	if s.Healthy {
		t.Error("expected Healthy=false when many keys are missing")
	}
}
