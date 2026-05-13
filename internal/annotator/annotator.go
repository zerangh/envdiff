package annotator

import (
	"fmt"
	"strings"

	"github.com/your/envdiff/internal/differ"
)

// Severity represents the importance level of an annotation.
type Severity string

const (
	SeverityInfo    Severity = "info"
	SeverityWarning Severity = "warning"
	SeverityCritical Severity = "critical"
)

// Annotation holds a human-readable note attached to a specific key.
type Annotation struct {
	Key      string   `json:"key"`
	Message  string   `json:"message"`
	Severity Severity `json:"severity"`
}

// Result is the full set of annotations produced for a diff result.
type Result struct {
	Annotations []Annotation `json:"annotations"`
}

// Format returns a plain-text representation of all annotations.
func (r Result) Format() string {
	if len(r.Annotations) == 0 {
		return "No annotations."
	}
	var sb strings.Builder
	for _, a := range r.Annotations {
		fmt.Fprintf(&sb, "[%s] %s: %s\n", strings.ToUpper(string(a.Severity)), a.Key, a.Message)
	}
	return strings.TrimRight(sb.String(), "\n")
}

// Annotate inspects a differ.Result and produces contextual annotations
// describing why each difference exists and what action is recommended.
func Annotate(result differ.Result) Result {
	var annotations []Annotation

	for _, key := range result.MissingInRight {
		annotations = append(annotations, Annotation{
			Key:      key,
			Message:  "Key is defined in the left environment but absent in the right. Consider adding it.",
			Severity: SeverityWarning,
		})
	}

	for _, key := range result.MissingInLeft {
		annotations = append(annotations, Annotation{
			Key:      key,
			Message:  "Key is defined in the right environment but absent in the left. It may be a new addition.",
			Severity: SeverityInfo,
		})
	}

	for _, mm := range result.Mismatched {
		sev := SeverityWarning
		msg := fmt.Sprintf("Values differ: left=%q, right=%q. Verify which value is correct.", mm.Left, mm.Right)
		if mm.Left == "" || mm.Right == "" {
			sev = SeverityCritical
			msg = fmt.Sprintf("One side has an empty value: left=%q, right=%q. This may indicate a misconfiguration.", mm.Left, mm.Right)
		}
		annotations = append(annotations, Annotation{
			Key:      mm.Key,
			Message:  msg,
			Severity: sev,
		})
	}

	return Result{Annotations: annotations}
}
