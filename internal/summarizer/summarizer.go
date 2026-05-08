package summarizer

import (
	"fmt"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/scorer"
	"github.com/user/envdiff/internal/validator"
)

// Summary holds a human-readable overview of a diff result.
type Summary struct {
	TotalKeys      int
	MissingInLeft  int
	MissingInRight int
	Mismatched     int
	Warnings       int
	Score          int
	Grade          string
	Healthy        bool
}

// Summarize produces a Summary from a DiffResult.
func Summarize(result differ.Result) Summary {
	warnings := validator.Validate(result)

	scoreOpts := scorer.Options{
		MissingWeight:   10,
		MismatchWeight:  5,
		WarningWeight:   2,
	}
	sc := scorer.Calculate(result, warnings, scoreOpts)
	grade := scorer.Grade(sc)

	total := countTotalKeys(result)

	return Summary{
		TotalKeys:      total,
		MissingInLeft:  len(result.MissingInLeft),
		MissingInRight: len(result.MissingInRight),
		Mismatched:     len(result.Mismatched),
		Warnings:       len(warnings),
		Score:          sc,
		Grade:          grade,
		Healthy:        sc >= 80,
	}
}

// Format returns a multi-line string representation of the Summary.
func (s Summary) Format() string {
	status := "✔ healthy"
	if !s.Healthy {
		status = "✘ needs attention"
	}
	return fmt.Sprintf(
		"Summary\n"+
			"  Total keys:       %d\n"+
			"  Missing in left:  %d\n"+
			"  Missing in right: %d\n"+
			"  Mismatched:       %d\n"+
			"  Warnings:         %d\n"+
			"  Score:            %d / 100 (%s)\n"+
			"  Status:           %s\n",
		s.TotalKeys,
		s.MissingInLeft,
		s.MissingInRight,
		s.Mismatched,
		s.Warnings,
		s.Score, s.Grade,
		status,
	)
}

func countTotalKeys(result differ.Result) int {
	seen := make(map[string]struct{})
	for _, k := range result.MissingInLeft {
		seen[k] = struct{}{}
	}
	for _, k := range result.MissingInRight {
		seen[k] = struct{}{}
	}
	for _, m := range result.Mismatched {
		seen[m.Key] = struct{}{}
	}
	for _, k := range result.Common {
		seen[k] = struct{}{}
	}
	return len(seen)
}
