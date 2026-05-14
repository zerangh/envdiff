// Package differ compares two parsed .env maps and returns structured differences.
package differ

import (
	"fmt"

	"github.com/yourorg/envdiff/internal/parser"
)

// Result holds the outcome of a diff between two .env files.
type Result struct {
	// MissingInRight contains keys present in left but absent in right.
	MissingInRight []string
	// MissingInLeft contains keys present in right but absent in left.
	MissingInLeft []string
	// Mismatched contains keys present in both files but with differing values.
	Mismatched []MismatchedKey
	// LeftFile is the path of the left (reference) file.
	LeftFile string
	// RightFile is the path of the right (target) file.
	RightFile string
	// LeftEnv holds the full parsed contents of the left file.
	LeftEnv map[string]string
	// RightEnv holds the full parsed contents of the right file.
	RightEnv map[string]string
}

// MismatchedKey represents a key whose value differs between two environments.
type MismatchedKey struct {
	Key        string
	LeftValue  string
	RightValue string
}

// Diff compares two env maps and returns a Result describing their differences.
func Diff(left, right map[string]string) Result {
	var result Result
	result.LeftEnv = left
	result.RightEnv = right

	for k, lv := range left {
		rv, ok := right[k]
		if !ok {
			result.MissingInRight = append(result.MissingInRight, k)
		} else if lv != rv {
			result.Mismatched = append(result.Mismatched, MismatchedKey{
				Key:        k,
				LeftValue:  lv,
				RightValue: rv,
			})
		}
	}

	for k := range right {
		if _, ok := left[k]; !ok {
			result.MissingInLeft = append(result.MissingInLeft, k)
		}
	}

	sortStrings(result.MissingInRight)
	sortStrings(result.MissingInLeft)
	sortMismatched(result.Mismatched)

	return result
}

// DiffFiles parses two .env files from disk and returns their diff Result.
func DiffFiles(leftPath, rightPath string) (Result, error) {
	left, err := parser.ParseFile(leftPath)
	if err != nil {
		return Result{}, fmt.Errorf("parsing left file %q: %w", leftPath, err)
	}
	right, err := parser.ParseFile(rightPath)
	if err != nil {
		return Result{}, fmt.Errorf("parsing right file %q: %w", rightPath, err)
	}
	res := Diff(left, right)
	res.LeftFile = leftPath
	res.RightFile = rightPath
	return res, nil
}
