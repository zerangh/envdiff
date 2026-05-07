package reporter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/differ"
)

func TestReportColor_Clean(t *testing.T) {
	var buf bytes.Buffer
	reportColor(&buf, differ.Result{})
	out := buf.String()
	if !strings.Contains(out, "No differences found") {
		t.Errorf("expected clean message, got: %q", out)
	}
	if !strings.Contains(out, colorCode.Green) {
		t.Errorf("expected green color code in clean output")
	}
}

func TestReportColor_MissingInRight(t *testing.T) {
	var buf bytes.Buffer
	reportColor(&buf, differ.Result{
		MissingInRight: []string{"SECRET_KEY", "DB_PASS"},
	})
	out := buf.String()
	if !strings.Contains(out, "MISSING IN RIGHT") {
		t.Errorf("expected MISSING IN RIGHT label, got: %q", out)
	}
	if !strings.Contains(out, "SECRET_KEY") {
		t.Errorf("expected SECRET_KEY in output")
	}
	if !strings.Contains(out, colorCode.Red) {
		t.Errorf("expected red color code for missing keys")
	}
}

func TestReportColor_MissingInLeft(t *testing.T) {
	var buf bytes.Buffer
	reportColor(&buf, differ.Result{
		MissingInLeft: []string{"NEW_FLAG"},
	})
	out := buf.String()
	if !strings.Contains(out, "MISSING IN LEFT") {
		t.Errorf("expected MISSING IN LEFT label, got: %q", out)
	}
	if !strings.Contains(out, "NEW_FLAG") {
		t.Errorf("expected NEW_FLAG in output")
	}
}

func TestReportColor_Mismatched(t *testing.T) {
	var buf bytes.Buffer
	reportColor(&buf, differ.Result{
		Mismatched: []differ.Mismatch{
			{Key: "LOG_LEVEL", LeftValue: "debug", RightValue: "info"},
		},
	})
	out := buf.String()
	if !strings.Contains(out, "MISMATCH") {
		t.Errorf("expected MISMATCH label, got: %q", out)
	}
	if !strings.Contains(out, "LOG_LEVEL") {
		t.Errorf("expected key LOG_LEVEL in output")
	}
	if !strings.Contains(out, colorCode.Yellow) {
		t.Errorf("expected yellow color code for mismatches")
	}
	if !strings.Contains(out, "debug") || !strings.Contains(out, "info") {
		t.Errorf("expected both values in mismatch output")
	}
}

func TestIsTerminal_Buffer(t *testing.T) {
	var buf bytes.Buffer
	if isTerminal(&buf) {
		t.Error("bytes.Buffer should not be detected as terminal")
	}
}
