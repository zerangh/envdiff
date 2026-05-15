package classifier

import (
	"sort"
	"strings"

	"github.com/user/envdiff/internal/differ"
)

// Category represents a semantic classification for an env key.
type Category string

const (
	CategoryDatabase    Category = "database"
	CategoryAuth        Category = "auth"
	CategoryNetwork     Category = "network"
	CategoryFeatureFlag Category = "feature_flag"
	CategoryLogging     Category = "logging"
	CategoryStorage     Category = "storage"
	CategoryUnknown     Category = "unknown"
)

// KeyClass holds a key and its assigned category.
type KeyClass struct {
	Key      string
	Category Category
}

// Result holds all classified keys from a diff result.
type Result struct {
	Classes    []KeyClass
	Categories map[Category][]string
}

var categoryPatterns = map[Category][]string{
	CategoryDatabase:    {"DB_", "DATABASE_", "POSTGRES", "MYSQL", "MONGO", "REDIS", "DSN"},
	CategoryAuth:        {"AUTH_", "JWT_", "SECRET", "TOKEN", "API_KEY", "OAUTH", "PASSWORD", "PASSWD"},
	CategoryNetwork:     {"HOST", "PORT", "URL", "ADDR", "ENDPOINT", "PROXY", "TLS_", "SSL_"},
	CategoryFeatureFlag: {"FEATURE_", "FLAG_", "ENABLE_", "DISABLE_", "FF_"},
	CategoryLogging:     {"LOG_", "LOGGING_", "LOG_LEVEL", "DEBUG", "TRACE"},
	CategoryStorage:     {"S3_", "BUCKET", "STORAGE_", "GCS_", "BLOB_", "MINIO_"},
}

// Classify assigns a Category to each unique key found in the diff result.
func Classify(result differ.Result) Result {
	keys := collectKeys(result)
	classes := make([]KeyClass, 0, len(keys))
	byCategory := make(map[Category][]string)

	for _, key := range keys {
		cat := classifyKey(key)
		classes = append(classes, KeyClass{Key: key, Category: cat})
		byCategory[cat] = append(byCategory[cat], key)
	}

	sort.Slice(classes, func(i, j int) bool {
		if classes[i].Category != classes[j].Category {
			return classes[i].Category < classes[j].Category
		}
		return classes[i].Key < classes[j].Key
	})

	return Result{Classes: classes, Categories: byCategory}
}

func classifyKey(key string) Category {
	upper := strings.ToUpper(key)
	for cat, patterns := range categoryPatterns {
		for _, p := range patterns {
			if strings.Contains(upper, p) {
				return cat
			}
		}
	}
	return CategoryUnknown
}

func collectKeys(result differ.Result) []string {
	seen := make(map[string]struct{})
	for _, k := range result.MissingInRight {
		seen[k] = struct{}{}
	}
	for _, k := range result.MissingInLeft {
		seen[k] = struct{}{}
	}
	for _, m := range result.Mismatched {
		seen[m.Key] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
