// Package tagger assigns semantic tags to environment variable keys based on
// naming heuristics. It categorises keys into groups such as "secret",
// "url", "database", "feature_flag", "deprecated", or "unknown".
//
// Tags are derived purely from key names (case-insensitive substring matching)
// and can be used downstream for filtering, reporting, or policy enforcement.
//
// Usage:
//
//	keys := []string{"DB_PASSWORD", "ENABLE_BETA", "APP_ENV"}
//	tags := tagger.TagKeys(keys)
//	// tags["DB_PASSWORD"] => ["secret"]
//	// tags["ENABLE_BETA"]  => ["feature_flag"]
//	// tags["APP_ENV"]      => ["unknown"]
//
//	res := tagger.TagResult(diffResult)
//	// res.Summary[tagger.TagSecret] => ["DB_PASSWORD", ...]
package tagger
