package highlighter

import (
	"fmt"
	"strings"

	"github.com/user/envdiff/internal/differ"
)

// Highlight represents a single highlighted segment within a value.
type Highlight struct {
	Key    string
	Left   string
	Right  string
	Diff   []Segment
}

// Segment is a portion of a value string with a change type annotation.
type Segment struct {
	Text    string
	Added   bool
	Removed bool
}

// HighlightResult holds all inline highlights for mismatched keys.
type HighlightResult struct {
	Highlights []Highlight
}

// Compute produces inline diff highlights for all mismatched key-value pairs
// in the given differ.Result.
func Compute(result differ.Result) HighlightResult {
	var highlights []Highlight
	for _, m := range result.Mismatched {
		h := Highlight{
			Key:   m.Key,
			Left:  m.Left,
			Right: m.Right,
			Diff:  computeSegments(m.Left, m.Right),
		}
		highlights = append(highlights, h)
	}
	return HighlightResult{Highlights: highlights}
}

// computeSegments produces a simple word-level diff between two strings.
func computeSegments(left, right string) []Segment {
	lWords := strings.Fields(left)
	rWords := strings.Fields(right)

	var segments []Segment

	lSet := toSet(lWords)
	rSet := toSet(rWords)

	for _, w := range lWords {
		if !rSet[w] {
			segments = append(segments, Segment{Text: w, Removed: true})
		} else {
			segments = append(segments, Segment{Text: w})
		}
	}
	for _, w := range rWords {
		if !lSet[w] {
			segments = append(segments, Segment{Text: w, Added: true})
		}
	}
	return segments
}

func toSet(words []string) map[string]bool {
	m := make(map[string]bool, len(words))
	for _, w := range words {
		m[w] = true
	}
	return m
}

// Format returns a human-readable string showing the inline diff segments.
func Format(h Highlight) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Key: %s\n", h.Key))
	for _, seg := range h.Diff {
		switch {
		case seg.Added:
			sb.WriteString(fmt.Sprintf("  [+] %s\n", seg.Text))
		case seg.Removed:
			sb.WriteString(fmt.Sprintf("  [-] %s\n", seg.Text))
		default:
			sb.WriteString(fmt.Sprintf("  [ ] %s\n", seg.Text))
		}
	}
	return sb.String()
}
