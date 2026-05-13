// Package differ computes the difference between two parsed .env maps.
package differ

import "github.com/user/envdiff/internal/parser"

// Result holds the full diff output between two .env files.
type Result struct {
	// MissingInLeft contains keys present in right but absent in left.
	MissingInLeft []string
	// MissingInRight contains keys present in left but absent in right.
	MissingInRight []string
	// Mismatched contains keys present in both files but with differing values.
	Mismatched []MismatchedKey
	// LeftFile is the label/path of the left environment file.
	LeftFile string
	// RightFile is the label/path of the right environment file.
	RightFile string
}

// MismatchedKey describes a single key whose value differs between environments.
type MismatchedKey struct {
	Key        string
	LeftValue  string
	RightValue string
}

// Diff computes the difference between two parsed env maps.
// leftLabel and rightLabel are used for display purposes in the result.
func Diff(left, right map[string]string, leftLabel, rightLabel string) Result {
	result := Result{
		LeftFile:  leftLabel,
		RightFile: rightLabel,
	}

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

	sortStrings(result.MissingInLeft)
	sortStrings(result.MissingInRight)
	sortMismatched(result.Mismatched)

	return result
}

// DiffFiles parses both files and returns the diff result.
func DiffFiles(leftPath, rightPath string) (Result, error) {
	left, err := parser.ParseFile(leftPath)
	if err != nil {
		return Result{}, err
	}
	right, err := parser.ParseFile(rightPath)
	if err != nil {
		return Result{}, err
	}
	return Diff(left, right, leftPath, rightPath), nil
}
