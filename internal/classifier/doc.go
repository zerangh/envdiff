// Package classifier assigns semantic categories to environment variable keys
// found in a diff result.
//
// Keys are matched against well-known naming patterns to determine whether they
// belong to categories such as database, auth, network, feature flags, logging,
// or storage. Keys that do not match any pattern are classified as "unknown".
//
// Usage:
//
//	res := classifier.Classify(diffResult)
//	for _, c := range res.Classes {
//		fmt.Printf("%s => %s\n", c.Key, c.Category)
//	}
package classifier
