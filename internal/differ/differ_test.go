package differ_test

import (
	"testing"

	"github.com/user/envdiff/internal/differ"
)

func TestDiff_NoDifferences(t *testing.T) {
	left := map[string]string{"FOO": "bar", "BAZ": "qux"}
	right := map[string]string{"FOO": "bar", "BAZ": "qux"}
	res := differ.Diff(left, right, "left.env", "right.env")
	if len(res.MissingInLeft) != 0 || len(res.MissingInRight) != 0 || len(res.Mismatched) != 0 {
		t.Errorf("expected no differences, got %+v", res)
	}
}

func TestDiff_MissingInRight(t *testing.T) {
	left := map[string]string{"FOO": "bar", "ONLY_LEFT": "val"}
	right := map[string]string{"FOO": "bar"}
	res := differ.Diff(left, right, "left.env", "right.env")
	if len(res.MissingInRight) != 1 || res.MissingInRight[0] != "ONLY_LEFT" {
		t.Errorf("expected ONLY_LEFT missing in right, got %v", res.MissingInRight)
	}
}

func TestDiff_MissingInLeft(t *testing.T) {
	left := map[string]string{"FOO": "bar"}
	right := map[string]string{"FOO": "bar", "ONLY_RIGHT": "val"}
	res := differ.Diff(left, right, "left.env", "right.env")
	if len(res.MissingInLeft) != 1 || res.MissingInLeft[0] != "ONLY_RIGHT" {
		t.Errorf("expected ONLY_RIGHT missing in left, got %v", res.MissingInLeft)
	}
}

func TestDiff_MismatchedValues(t *testing.T) {
	left := map[string]string{"FOO": "bar"}
	right := map[string]string{"FOO": "baz"}
	res := differ.Diff(left, right, "left.env", "right.env")
	if len(res.Mismatched) != 1 {
		t.Fatalf("expected 1 mismatch, got %d", len(res.Mismatched))
	}
	m := res.Mismatched[0]
	if m.Key != "FOO" || m.LeftValue != "bar" || m.RightValue != "baz" {
		t.Errorf("unexpected mismatch: %+v", m)
	}
}

func TestDiff_SortedOutput(t *testing.T) {
	left := map[string]string{"Z_KEY": "1", "A_KEY": "2", "M_KEY": "3"}
	right := map[string]string{}
	res := differ.Diff(left, right, "left.env", "right.env")
	for i := 1; i < len(res.MissingInRight); i++ {
		if res.MissingInRight[i-1] > res.MissingInRight[i] {
			t.Errorf("MissingInRight not sorted: %v", res.MissingInRight)
		}
	}
}

func TestDiff_Labels(t *testing.T) {
	res := differ.Diff(map[string]string{}, map[string]string{}, "a.env", "b.env")
	if res.LeftFile != "a.env" || res.RightFile != "b.env" {
		t.Errorf("unexpected labels: left=%q right=%q", res.LeftFile, res.RightFile)
	}
}
