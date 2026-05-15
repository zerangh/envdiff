package classifier_test

import (
	"testing"

	"github.com/user/envdiff/internal/classifier"
	"github.com/user/envdiff/internal/differ"
)

func makeResult(missingRight, missingLeft []string, mismatched []differ.Mismatch) differ.Result {
	return differ.Result{
		MissingInRight: missingRight,
		MissingInLeft:  missingLeft,
		Mismatched:     mismatched,
	}
}

func TestClassify_EmptyResult(t *testing.T) {
	res := classifier.Classify(makeResult(nil, nil, nil))
	if len(res.Classes) != 0 {
		t.Errorf("expected 0 classes, got %d", len(res.Classes))
	}
}

func TestClassify_DatabaseKey(t *testing.T) {
	res := classifier.Classify(makeResult([]string{"DB_HOST", "DATABASE_URL"}, nil, nil))
	for _, c := range res.Classes {
		if c.Category != classifier.CategoryDatabase {
			t.Errorf("key %s: expected database, got %s", c.Key, c.Category)
		}
	}
}

func TestClassify_AuthKey(t *testing.T) {
	res := classifier.Classify(makeResult([]string{"JWT_SECRET", "API_KEY"}, nil, nil))
	for _, c := range res.Classes {
		if c.Category != classifier.CategoryAuth {
			t.Errorf("key %s: expected auth, got %s", c.Key, c.Category)
		}
	}
}

func TestClassify_NetworkKey(t *testing.T) {
	res := classifier.Classify(makeResult(nil, []string{"APP_PORT", "SERVICE_HOST"}, nil))
	for _, c := range res.Classes {
		if c.Category != classifier.CategoryNetwork {
			t.Errorf("key %s: expected network, got %s", c.Key, c.Category)
		}
	}
}

func TestClassify_UnknownKey(t *testing.T) {
	res := classifier.Classify(makeResult([]string{"FOOBAR"}, nil, nil))
	if len(res.Classes) != 1 || res.Classes[0].Category != classifier.CategoryUnknown {
		t.Errorf("expected unknown category for FOOBAR")
	}
}

func TestClassify_MixedKeys(t *testing.T) {
	result := makeResult(
		[]string{"DB_HOST", "LOG_LEVEL"},
		[]string{"S3_BUCKET"},
		[]differ.Mismatch{{Key: "FEATURE_X", Left: "true", Right: "false"}},
	)
	res := classifier.Classify(result)
	if len(res.Classes) != 4 {
		t.Fatalf("expected 4 classes, got %d", len(res.Classes))
	}
	catMap := make(map[string]classifier.Category)
	for _, c := range res.Classes {
		catMap[c.Key] = c.Category
	}
	if catMap["DB_HOST"] != classifier.CategoryDatabase {
		t.Errorf("DB_HOST should be database")
	}
	if catMap["LOG_LEVEL"] != classifier.CategoryLogging {
		t.Errorf("LOG_LEVEL should be logging")
	}
	if catMap["S3_BUCKET"] != classifier.CategoryStorage {
		t.Errorf("S3_BUCKET should be storage")
	}
	if catMap["FEATURE_X"] != classifier.CategoryFeatureFlag {
		t.Errorf("FEATURE_X should be feature_flag")
	}
}

func TestClassify_CategoriesMap(t *testing.T) {
	result := makeResult([]string{"DB_HOST", "DB_PORT"}, nil, nil)
	res := classifier.Classify(result)
	dbKeys := res.Categories[classifier.CategoryDatabase]
	if len(dbKeys) != 2 {
		t.Errorf("expected 2 database keys, got %d", len(dbKeys))
	}
}
