package differ_test

import (
	"testing"

	"github.com/user/envdiff/internal/differ"
)

func TestDiff_NoDifferences(t *testing.T) {
	left := map[string]string{"APP_ENV": "dev", "PORT": "8080"}
	right := map[string]string{"APP_ENV": "dev", "PORT": "8080"}

	result := differ.Diff(left, right)

	if result.HasDiff() {
		t.Errorf("expected no diff, got %+v", result)
	}
}

func TestDiff_MissingInRight(t *testing.T) {
	left := map[string]string{"APP_ENV": "dev", "SECRET": "abc"}
	right := map[string]string{"APP_ENV": "dev"}

	result := differ.Diff(left, right)

	if len(result.MissingInRight) != 1 || result.MissingInRight[0] != "SECRET" {
		t.Errorf("expected SECRET missing in right, got %v", result.MissingInRight)
	}
	if len(result.MissingInLeft) != 0 {
		t.Errorf("expected no keys missing in left, got %v", result.MissingInLeft)
	}
}

func TestDiff_MissingInLeft(t *testing.T) {
	left := map[string]string{"APP_ENV": "dev"}
	right := map[string]string{"APP_ENV": "dev", "DB_URL": "postgres://localhost/db"}

	result := differ.Diff(left, right)

	if len(result.MissingInLeft) != 1 || result.MissingInLeft[0] != "DB_URL" {
		t.Errorf("expected DB_URL missing in left, got %v", result.MissingInLeft)
	}
}

func TestDiff_MismatchedValues(t *testing.T) {
	left := map[string]string{"PORT": "8080", "APP_ENV": "dev"}
	right := map[string]string{"PORT": "9090", "APP_ENV": "dev"}

	result := differ.Diff(left, right)

	if len(result.Mismatched) != 1 {
		t.Fatalf("expected 1 mismatch, got %d", len(result.Mismatched))
	}
	m := result.Mismatched[0]
	if m.Key != "PORT" || m.LeftValue != "8080" || m.RightValue != "9090" {
		t.Errorf("unexpected mismatch entry: %+v", m)
	}
}

func TestDiff_SortedOutput(t *testing.T) {
	left := map[string]string{"Z_KEY": "1", "A_KEY": "2", "M_KEY": "3"}
	right := map[string]string{}

	result := differ.Diff(left, right)

	expected := []string{"A_KEY", "M_KEY", "Z_KEY"}
	for i, k := range result.MissingInRight {
		if k != expected[i] {
			t.Errorf("expected sorted key %q at index %d, got %q", expected[i], i, k)
		}
	}
}

func TestDiff_EmptyMaps(t *testing.T) {
	result := differ.Diff(map[string]string{}, map[string]string{})
	if result.HasDiff() {
		t.Errorf("expected no diff for two empty maps")
	}
}
