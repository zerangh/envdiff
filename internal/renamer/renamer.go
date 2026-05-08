// Package renamer provides utilities for detecting and suggesting key renames
// between two .env files. It uses similarity scoring to identify keys that
// likely represent the same configuration value under a different name.
package renamer

import (
	"strings"

	"github.com/user/envdiff/internal/differ"
)

// Suggestion represents a possible key rename between two environments.
type Suggestion struct {
	LeftKey  string
	RightKey string
	Score    float64
}

// Options controls the behaviour of the rename detector.
type Options struct {
	// MinScore is the minimum similarity score (0.0–1.0) required to emit a
	// suggestion. Defaults to 0.6 if zero.
	MinScore float64
}

// Detect analyses the missing keys in a differ.Result and returns a list of
// rename suggestions ordered by descending similarity score.
func Detect(result differ.Result, opts Options) []Suggestion {
	if opts.MinScore == 0 {
		opts.MinScore = 0.6
	}

	var suggestions []Suggestion

	for _, left := range result.MissingInRight {
		for _, right := range result.MissingInLeft {
			score := similarity(left, right)
			if score >= opts.MinScore {
				suggestions = append(suggestions, Suggestion{
					LeftKey:  left,
					RightKey: right,
					Score:    score,
				})
			}
		}
	}

	sortSuggestions(suggestions)
	return suggestions
}

// similarity returns a value in [0, 1] representing how alike two key strings
// are. It normalises both strings to lowercase and computes a bigram-based
// Dice coefficient.
func similarity(a, b string) float64 {
	a = strings.ToLower(a)
	b = strings.ToLower(b)

	if a == b {
		return 1.0
	}

	bigramsA := bigrams(a)
	bigramsB := bigrams(b)

	if len(bigramsA) == 0 && len(bigramsB) == 0 {
		return 1.0
	}
	if len(bigramsA) == 0 || len(bigramsB) == 0 {
		return 0.0
	}

	intersection := 0
	for bg, countA := range bigramsA {
		if countB, ok := bigramsB[bg]; ok {
			if countA < countB {
				intersection += countA
			} else {
				intersection += countB
			}
		}
	}

	return float64(2*intersection) / float64(len(bigramsA)+len(bigramsB))
}

// bigrams returns a frequency map of character bigrams for s.
func bigrams(s string) map[string]int {
	m := make(map[string]int)
	for i := 0; i < len(s)-1; i++ {
		m[s[i:i+2]]++
	}
	return m
}

// sortSuggestions sorts suggestions by descending Score, then alphabetically
// by LeftKey for determinism.
func sortSuggestions(ss []Suggestion) {
	for i := 1; i < len(ss); i++ {
		for j := i; j > 0; j-- {
			a, b := ss[j-1], ss[j]
			if a.Score < b.Score || (a.Score == b.Score && a.LeftKey > b.LeftKey) {
				ss[j-1], ss[j] = ss[j], ss[j-1]
			}
		}
	}
}
