package splitter

import (
	"sort"

	"github.com/your-org/envdiff/internal/differ"
)

// Strategy controls how keys are assigned to buckets when splitting.
type Strategy string

const (
	// StrategyAlpha groups keys alphabetically by first character.
	StrategyAlpha Strategy = "alpha"
	// StrategyPrefix groups keys by underscore-delimited prefix.
	StrategyPrefix Strategy = "prefix"
)

// Options configures the Split operation.
type Options struct {
	Strategy Strategy
	// MaxBuckets limits the number of output buckets (0 = unlimited).
	MaxBuckets int
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Strategy:   StrategyPrefix,
		MaxBuckets: 0,
	}
}

// Bucket holds a named slice of keys and their associated diff result.
type Bucket struct {
	Name   string
	Keys   []string
	Result differ.Result
}

// Split partitions a differ.Result into named buckets according to opts.
func Split(result differ.Result, opts Options) []Bucket {
	keySet := collectKeys(result)

	groups := make(map[string][]string)
	for _, k := range keySet {
		var label string
		switch opts.Strategy {
		case StrategyAlpha:
			if len(k) > 0 {
				label = string([]rune(k)[:1])
			} else {
				label = "_"
			}
		default: // StrategyPrefix
			label = extractPrefix(k)
		}
		groups[label] = append(groups[label], k)
	}

	labels := make([]string, 0, len(groups))
	for l := range groups {
		labels = append(labels, l)
	}
	sort.Strings(labels)

	if opts.MaxBuckets > 0 && len(labels) > opts.MaxBuckets {
		labels = labels[:opts.MaxBuckets]
	}

	buckets := make([]Bucket, 0, len(labels))
	for _, label := range labels {
		keys := groups[label]
		sort.Strings(keys)
		buckets = append(buckets, Bucket{
			Name:   label,
			Keys:   keys,
			Result: filterResult(result, keys),
		})
	}
	return buckets
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
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	return out
}

func extractPrefix(key string) string {
	for i, ch := range key {
		if ch == '_' && i > 0 {
			return key[:i]
		}
	}
	return key
}

func filterResult(result differ.Result, keys []string) differ.Result {
	keySet := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		keySet[k] = struct{}{}
	}

	var r differ.Result
	for _, k := range result.MissingInRight {
		if _, ok := keySet[k]; ok {
			r.MissingInRight = append(r.MissingInRight, k)
		}
	}
	for _, k := range result.MissingInLeft {
		if _, ok := keySet[k]; ok {
			r.MissingInLeft = append(r.MissingInLeft, k)
		}
	}
	for _, m := range result.Mismatched {
		if _, ok := keySet[m.Key]; ok {
			r.Mismatched = append(r.Mismatched, m)
		}
	}
	return r
}
