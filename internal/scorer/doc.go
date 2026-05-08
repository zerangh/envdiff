// Package scorer computes a numeric health score for an env diff result.
//
// A score of 100 indicates that two .env files are identical. Penalties
// are applied for missing keys and mismatched values, with configurable
// weights for each category.
//
// Example usage:
//
//	result := differ.Diff(left, right)
//	score := scorer.Calculate(result, scorer.Options{})
//	fmt.Printf("Score: %d (%s)\n", score.Total, scorer.Grade(score.Total))
package scorer
