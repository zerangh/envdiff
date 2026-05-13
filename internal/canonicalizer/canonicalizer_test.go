package canonicalizer_test

import (
	"testing"

	"github.com/user/envdiff/internal/canonicalizer"
	"github.com/user/envdiff/internal/differ"
)

func TestCanonicalize_DefaultOptions(t *testing.T) {
	opts := canonicalizer.DefaultOptions()
	cases := []struct {
		input string
		want  string
	}{
		{"db-host", "DB_HOST"},
		{"db.port", "DB_PORT"},
		{"API_KEY", "API_KEY"},
		{"my.service-url", "MY_SERVICE_URL"},
		{"lowercase", "LOWERCASE"},
	}
	for _, tc := range cases {
		got := canonicalizer.Canonicalize(tc.input, opts)
		if got != tc.want {
			t.Errorf("Canonicalize(%q) = %q; want %q", tc.input, got, tc.want)
		}
	}
}

func TestCanonicalize_NoCaseNorm(t *testing.T) {
	opts := canonicalizer.Options{NormalizeCase: false, NormalizeSeparators: true}
	got := canonicalizer.Canonicalize("db-host", opts)
	if got != "db_host" {
		t.Errorf("expected db_host, got %q", got)
	}
}

func TestCanonicalize_NoSepNorm(t *testing.T) {
	opts := canonicalizer.Options{NormalizeCase: true, NormalizeSeparators: false}
	got := canonicalizer.Canonicalize("db-host", opts)
	if got != "DB-HOST" {
		t.Errorf("expected DB-HOST, got %q", got)
	}
}

func TestNormalizeEnv_CollapsesSameCanonical(t *testing.T) {
	env := map[string]string{
		"db-host": "localhost",
		"db.host": "remotehost",
	}
	opts := canonicalizer.DefaultOptions()
	norm := canonicalizer.NormalizeEnv(env, opts)
	// Both collapse to DB_HOST; one value survives.
	if _, ok := norm["DB_HOST"]; !ok {
		t.Error("expected DB_HOST key in normalized map")
	}
	if len(norm) != 1 {
		t.Errorf("expected 1 key after collapse, got %d", len(norm))
	}
}

func TestNormalizeEnv_PreservesValues(t *testing.T) {
	env := map[string]string{"api-key": "secret", "PORT": "8080"}
	opts := canonicalizer.DefaultOptions()
	norm := canonicalizer.NormalizeEnv(env, opts)
	if norm["API_KEY"] != "secret" {
		t.Errorf("expected API_KEY=secret, got %q", norm["API_KEY"])
	}
	if norm["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", norm["PORT"])
	}
}

func TestNormalizeResult_NormalizesAllFields(t *testing.T) {
	input := differ.Result{
		MissingInLeft:  []string{"db-host"},
		MissingInRight: []string{"api.key"},
		Mismatched: []differ.Mismatch{
			{Key: "my-service-url", LeftValue: "http://a", RightValue: "http://b"},
		},
	}
	opts := canonicalizer.DefaultOptions()
	got := canonicalizer.NormalizeResult(input, opts)

	if got.MissingInLeft[0] != "DB_HOST" {
		t.Errorf("MissingInLeft[0]: want DB_HOST, got %q", got.MissingInLeft[0])
	}
	if got.MissingInRight[0] != "API_KEY" {
		t.Errorf("MissingInRight[0]: want API_KEY, got %q", got.MissingInRight[0])
	}
	if got.Mismatched[0].Key != "MY_SERVICE_URL" {
		t.Errorf("Mismatched[0].Key: want MY_SERVICE_URL, got %q", got.Mismatched[0].Key)
	}
	if got.Mismatched[0].LeftValue != "http://a" {
		t.Error("LeftValue should be unchanged")
	}
}
