package tagger

import (
	"sort"
	"strings"

	"github.com/your-org/envdiff/internal/differ"
)

// Tag represents a label assigned to an env key based on heuristics.
type Tag string

const (
	TagSecret    Tag = "secret"
	TagURL       Tag = "url"
	TagFeatureFlag Tag = "feature_flag"
	TagDatabase  Tag = "database"
	TagDeprecated Tag = "deprecated"
	TagUnknown   Tag = "unknown"
)

// KeyTags maps a key name to its assigned tags.
type KeyTags map[string][]Tag

// Result holds the tagging output for a diff result.
type Result struct {
	Tags    KeyTags  `json:"tags"`
	Summary map[Tag][]string `json:"summary"`
}

var secretPatterns = []string{"secret", "password", "passwd", "token", "apikey", "api_key", "private", "credential"}
var urlPatterns = []string{"url", "uri", "endpoint", "host", "addr", "address", "dsn"}
var dbPatterns = []string{"db", "database", "postgres", "mysql", "mongo", "redis", "sqlite"}
var featurePatterns = []string{"enable_", "disable_", "feature_", "flag_", "_enabled", "_disabled"}
var deprecatedPatterns = []string{"old_", "legacy_", "deprecated_", "_old", "_legacy"}

// TagKeys assigns tags to a set of keys.
func TagKeys(keys []string) KeyTags {
	result := make(KeyTags, len(keys))
	for _, k := range keys {
		result[k] = assignTags(k)
	}
	return result
}

// TagResult derives all unique keys from a differ.Result and tags them.
func TagResult(r differ.Result) Result {
	keys := collectKeys(r)
	tags := TagKeys(keys)
	summary := buildSummary(tags)
	return Result{Tags: tags, Summary: summary}
}

func assignTags(key string) []Tag {
	lower := strings.ToLower(key)
	var tags []Tag
	for _, p := range secretPatterns {
		if strings.Contains(lower, p) {
			tags = append(tags, TagSecret)
			break
		}
	}
	for _, p := range urlPatterns {
		if strings.Contains(lower, p) {
			tags = append(tags, TagURL)
			break
		}
	}
	for _, p := range dbPatterns {
		if strings.Contains(lower, p) {
			tags = append(tags, TagDatabase)
			break
		}
	}
	for _, p := range featurePatterns {
		if strings.Contains(lower, p) {
			tags = append(tags, TagFeatureFlag)
			break
		}
	}
	for _, p := range deprecatedPatterns {
		if strings.Contains(lower, p) {
			tags = append(tags, TagDeprecated)
			break
		}
	}
	if len(tags) == 0 {
		tags = append(tags, TagUnknown)
	}
	return tags
}

func collectKeys(r differ.Result) []string {
	seen := make(map[string]struct{})
	for _, k := range r.MissingInRight {
		seen[k] = struct{}{}
	}
	for _, k := range r.MissingInLeft {
		seen[k] = struct{}{}
	}
	for _, m := range r.Mismatched {
		seen[m.Key] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func buildSummary(tags KeyTags) map[Tag][]string {
	summary := make(map[Tag][]string)
	for key, ts := range tags {
		for _, t := range ts {
			summary[t] = append(summary[t], key)
		}
	}
	for t := range summary {
		sort.Strings(summary[t])
	}
	return summary
}
