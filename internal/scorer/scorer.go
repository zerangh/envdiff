package scorer

import (
	"github.com/user/envdiff/internal/differ"
)

// Score represents a numeric health score for an env diff result.
type Score struct {
	// Total is the overall score from 0 to 100.
	Total int
	// MissingPenalty is the total points deducted for missing keys.
	MissingPenalty int
	// MismatchPenalty is the total points deducted for mismatched values.
	MismatchPenalty int
	// KeyCount is the total number of unique keys across both files.
	KeyCount int
}

// Options configures how penalties are calculated.
type Options struct {
	// MissingWeight is the penalty per missing key (default: 10).
	MissingWeight int
	// MismatchWeight is the penalty per mismatched value (default: 5).
	MismatchWeight int
}

func defaultOptions(o Options) Options {
	if o.MissingWeight == 0 {
		o.MissingWeight = 10
	}
	if o.MismatchWeight == 0 {
		o.MismatchWeight = 5
	}
	return o
}

// Calculate computes a health score given a diff result.
// A score of 100 means the environments are identical.
// Penalties are applied for missing and mismatched keys.
func Calculate(result differ.Result, opts Options) Score {
	opts = defaultOptions(opts)

	totalKeys := result.KeyCount()
	if totalKeys == 0 {
		return Score{Total: 100, KeyCount: 0}
	}

	missingCount := len(result.MissingInRight) + len(result.MissingInLeft)
	mismatchCount := len(result.Mismatched)

	missingPenalty := missingCount * opts.MissingWeight
	mismatchPenalty := mismatchCount * opts.MismatchWeight
	totalPenalty := missingPenalty + mismatchPenalty

	score := 100 - totalPenalty
	if score < 0 {
		score = 0
	}

	return Score{
		Total:           score,
		MissingPenalty:  missingPenalty,
		MismatchPenalty: mismatchPenalty,
		KeyCount:        totalKeys,
	}
}

// Grade returns a letter grade for the given score total.
func Grade(score int) string {
	switch {
	case score >= 90:
		return "A"
	case score >= 75:
		return "B"
	case score >= 60:
		return "C"
	case score >= 40:
		return "D"
	default:
		return "F"
	}
}
