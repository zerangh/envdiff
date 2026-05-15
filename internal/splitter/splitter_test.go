package splitter_test

import (
	"testing"

	"github.com/your-org/envdiff/internal/differ"
	"github.com/your-org/envdiff/internal/splitter"
)

func makeResult() differ.Result {
	return differ.Result{
		MissingInRight: []string{"APP_HOST", "DB_PORT", "CACHE_TTL"},
		MissingInLeft:  []string{"APP_SECRET"},
		Mismatched: []differ.Mismatch{
			{Key: "DB_NAME", Left: "dev", Right: "prod"},
			{Key: "APP_ENV", Left: "development", Right: "production"},
		},
	}
}

func TestSplit_PrefixStrategy(t *testing.T) {
	result := makeResult()
	opts := splitter.DefaultOptions()
	buckets := splitter.Split(result, opts)

	if len(buckets) == 0 {
		t.Fatal("expected at least one bucket")
	}

	names := make(map[string]bool)
	for _, b := range buckets {
		names[b.Name] = true
	}
	if !names["APP"] {
		t.Error("expected APP prefix bucket")
	}
	if !names["DB"] {
		t.Error("expected DB prefix bucket")
	}
	if !names["CACHE"] {
		t.Error("expected CACHE prefix bucket")
	}
}

func TestSplit_AlphaStrategy(t *testing.T) {
	result := makeResult()
	opts := splitter.Options{Strategy: splitter.StrategyAlpha}
	buckets := splitter.Split(result, opts)

	names := make(map[string]bool)
	for _, b := range buckets {
		names[b.Name] = true
	}
	if !names["A"] {
		t.Error("expected 'A' alpha bucket")
	}
	if !names["D"] {
		t.Error("expected 'D' alpha bucket")
	}
	if !names["C"] {
		t.Error("expected 'C' alpha bucket")
	}
}

func TestSplit_MaxBuckets(t *testing.T) {
	result := makeResult()
	opts := splitter.Options{Strategy: splitter.StrategyPrefix, MaxBuckets: 2}
	buckets := splitter.Split(result, opts)

	if len(buckets) > 2 {
		t.Errorf("expected at most 2 buckets, got %d", len(buckets))
	}
}

func TestSplit_EmptyResult(t *testing.T) {
	result := differ.Result{}
	buckets := splitter.Split(result, splitter.DefaultOptions())
	if len(buckets) != 0 {
		t.Errorf("expected 0 buckets for empty result, got %d", len(buckets))
	}
}

func TestSplit_BucketResultCorrect(t *testing.T) {
	result := makeResult()
	buckets := splitter.Split(result, splitter.DefaultOptions())

	for _, b := range buckets {
		if b.Name == "DB" {
			if len(b.Result.MissingInRight) != 1 || b.Result.MissingInRight[0] != "DB_PORT" {
				t.Errorf("DB bucket MissingInRight wrong: %v", b.Result.MissingInRight)
			}
			if len(b.Result.Mismatched) != 1 || b.Result.Mismatched[0].Key != "DB_NAME" {
				t.Errorf("DB bucket Mismatched wrong: %v", b.Result.Mismatched)
			}
			return
		}
	}
	t.Error("DB bucket not found")
}

func TestSplit_SortedBucketNames(t *testing.T) {
	result := makeResult()
	buckets := splitter.Split(result, splitter.DefaultOptions())

	for i := 1; i < len(buckets); i++ {
		if buckets[i].Name < buckets[i-1].Name {
			t.Errorf("buckets not sorted: %s before %s", buckets[i-1].Name, buckets[i].Name)
		}
	}
}
