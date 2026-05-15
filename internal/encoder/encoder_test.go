package encoder_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/encoder"
)

func TestEncode_UnknownFormat(t *testing.T) {
	_, err := encoder.Encode(map[string]string{"A": "1"}, encoder.Options{Format: "xml"})
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestEncode_DotenvBasic(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	opts := encoder.DefaultOptions()
	out, err := encoder.Encode(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "FOO=bar") {
		t.Errorf("expected FOO=bar in output, got:\n%s", out)
	}
	if !strings.Contains(out, "BAZ=qux") {
		t.Errorf("expected BAZ=qux in output, got:\n%s", out)
	}
}

func TestEncode_DotenvQuoteAll(t *testing.T) {
	env := map[string]string{"KEY": "value"}
	opts := encoder.Options{Format: encoder.FormatDotenv, QuoteAll: true, SortKeys: true}
	out, err := encoder.Encode(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `KEY="value"`) {
		t.Errorf("expected quoted value, got: %s", out)
	}
}

func TestEncode_DotenvOmitEmpty(t *testing.T) {
	env := map[string]string{"PRESENT": "yes", "EMPTY": ""}
	opts := encoder.Options{Format: encoder.FormatDotenv, OmitEmpty: true, SortKeys: true}
	out, err := encoder.Encode(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "EMPTY") {
		t.Errorf("expected EMPTY to be omitted, got: %s", out)
	}
	if !strings.Contains(out, "PRESENT=yes") {
		t.Errorf("expected PRESENT=yes in output, got: %s", out)
	}
}

func TestEncode_ExportsFormat(t *testing.T) {
	env := map[string]string{"PATH_VAR": "/usr/bin"}
	opts := encoder.Options{Format: encoder.FormatExports, SortKeys: true}
	out, err := encoder.Encode(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "export PATH_VAR=/usr/bin") {
		t.Errorf("expected export prefix, got: %s", out)
	}
}

func TestEncode_DockerFormat(t *testing.T) {
	env := map[string]string{"DB_HOST": "localhost"}
	opts := encoder.Options{Format: encoder.FormatDocker, SortKeys: true}
	out, err := encoder.Encode(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "DB_HOST=localhost") {
		t.Errorf("expected DB_HOST=localhost, got: %s", out)
	}
	if strings.Contains(out, "export ") {
		t.Errorf("docker format should not include export prefix")
	}
}

func TestEncode_SortedKeys(t *testing.T) {
	env := map[string]string{"Z_KEY": "z", "A_KEY": "a", "M_KEY": "m"}
	opts := encoder.Options{Format: encoder.FormatDotenv, SortKeys: true}
	out, err := encoder.Encode(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "A_KEY") {
		t.Errorf("expected first line to be A_KEY, got: %s", lines[0])
	}
	if !strings.HasPrefix(lines[2], "Z_KEY") {
		t.Errorf("expected last line to be Z_KEY, got: %s", lines[2])
	}
}
