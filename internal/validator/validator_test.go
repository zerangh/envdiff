package validator_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/validator"
)

func makeDiff(missingLeft, missingRight []string, mismatched []differ.Mismatch) differ.Result {
	return differ.Result{
		MissingInLeft:  missingLeft,
		MissingInRight: missingRight,
		Mismatched:     mismatched,
	}
}

func TestValidate_NoWarnings(t *testing.T) {
	diff := makeDiff(nil, nil, []differ.Mismatch{
		{Key: "HOST", LeftValue: "localhost", RightValue: "prod.example.com"},
	})
	left := map[string]string{"HOST": "localhost"}
	right := map[string]string{"HOST": "prod.example.com"}

	res := validator.Validate(diff, left, right)
	if res.HasWarnings() {
		t.Errorf("expected no warnings, got %d", len(res.Warnings))
	}
}

func TestValidate_EmptyLeftValue(t *testing.T) {
	diff := makeDiff(nil, nil, []differ.Mismatch{
		{Key: "SECRET", LeftValue: "", RightValue: "abc123"},
	})
	left := map[string]string{"SECRET": ""}
	right := map[string]string{"SECRET": "abc123"}

	res := validator.Validate(diff, left, right)
	if !res.HasWarnings() {
		t.Fatal("expected warnings for empty left value")
	}
	if res.Warnings[0].Key != "SECRET" {
		t.Errorf("unexpected key: %s", res.Warnings[0].Key)
	}
}

func TestValidate_WhitespaceValue(t *testing.T) {
	diff := makeDiff(nil, nil, []differ.Mismatch{
		{Key: "TOKEN", LeftValue: " abc", RightValue: "abc"},
	})
	left := map[string]string{"TOKEN": " abc"}
	right := map[string]string{"TOKEN": "abc"}

	res := validator.Validate(diff, left, right)
	if !res.HasWarnings() {
		t.Fatal("expected whitespace warning")
	}
	found := false
	for _, w := range res.Warnings {
		if strings.Contains(w.Message, "whitespace") {
			found = true
		}
	}
	if !found {
		t.Error("expected a whitespace warning message")
	}
}

func TestValidate_MissingInRightEmptyLeft(t *testing.T) {
	diff := makeDiff(nil, []string{"ORPHAN"}, nil)
	left := map[string]string{"ORPHAN": ""}
	right := map[string]string{}

	res := validator.Validate(diff, left, right)
	if !res.HasWarnings() {
		t.Fatal("expected warning for missing-in-right with empty left value")
	}
	if res.Warnings[0].Key != "ORPHAN" {
		t.Errorf("unexpected key: %s", res.Warnings[0].Key)
	}
}

func TestValidate_Format_NoWarnings(t *testing.T) {
	res := validator.Result{}
	out := res.Format()
	if !strings.Contains(out, "no validation warnings") {
		t.Errorf("unexpected format output: %s", out)
	}
}

func TestValidate_Format_WithWarnings(t *testing.T) {
	res := validator.Result{
		Warnings: []validator.Warning{
			{Key: "FOO", Message: "left value is empty"},
		},
	}
	out := res.Format()
	if !strings.Contains(out, "FOO") {
		t.Errorf("expected key FOO in output: %s", out)
	}
	if !strings.Contains(out, "1 validation warning") {
		t.Errorf("expected warning count in output: %s", out)
	}
}

func TestValidate_Format_MultipleWarnings(t *testing.T) {
	res := validator.Result{
		Warnings: []validator.Warning{
			{Key: "FOO", Message: "left value is empty"},
			{Key: "BAR", Message: "leading/trailing whitespace"},
		},
	}
	out := res.Format()
	if !strings.Contains(out, "FOO") {
		t.Errorf("expected key FOO in output: %s", out)
	}
	if !strings.Contains(out, "BAR") {
		t.Errorf("expected key BAR in output: %s", out)
	}
	if !strings.Contains(out, "2 validation warning") {
		t.Errorf("expected warning count in output: %s", out)
	}
}
