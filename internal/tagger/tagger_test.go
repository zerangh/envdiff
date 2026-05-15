package tagger

import (
	"testing"

	"github.com/your-org/envdiff/internal/differ"
)

func makeResult(missingRight, missingLeft []string, mismatched []differ.Mismatch) differ.Result {
	return differ.Result{
		MissingInRight: missingRight,
		MissingInLeft:  missingLeft,
		Mismatched:     mismatched,
	}
}

func TestTagKeys_Secret(t *testing.T) {
	tags := TagKeys([]string{"DB_PASSWORD", "API_TOKEN", "STRIPE_SECRET"})
	for _, k := range []string{"DB_PASSWORD", "API_TOKEN", "STRIPE_SECRET"} {
		if !containsTag(tags[k], TagSecret) {
			t.Errorf("expected %q to be tagged as secret", k)
		}
	}
}

func TestTagKeys_URL(t *testing.T) {
	tags := TagKeys([]string{"DATABASE_URL", "API_ENDPOINT", "SERVICE_HOST"})
	for _, k := range []string{"DATABASE_URL", "API_ENDPOINT", "SERVICE_HOST"} {
		if !containsTag(tags[k], TagURL) {
			t.Errorf("expected %q to be tagged as url", k)
		}
	}
}

func TestTagKeys_Database(t *testing.T) {
	tags := TagKeys([]string{"POSTGRES_DB", "REDIS_HOST", "MYSQL_USER"})
	for _, k := range []string{"POSTGRES_DB", "REDIS_HOST", "MYSQL_USER"} {
		if !containsTag(tags[k], TagDatabase) {
			t.Errorf("expected %q to be tagged as database", k)
		}
	}
}

func TestTagKeys_FeatureFlag(t *testing.T) {
	tags := TagKeys([]string{"ENABLE_DARK_MODE", "FEATURE_X", "BETA_ENABLED"})
	for _, k := range []string{"ENABLE_DARK_MODE", "FEATURE_X", "BETA_ENABLED"} {
		if !containsTag(tags[k], TagFeatureFlag) {
			t.Errorf("expected %q to be tagged as feature_flag", k)
		}
	}
}

func TestTagKeys_Unknown(t *testing.T) {
	tags := TagKeys([]string{"APP_ENV", "LOG_LEVEL"})
	for _, k := range []string{"APP_ENV", "LOG_LEVEL"} {
		if !containsTag(tags[k], TagUnknown) {
			t.Errorf("expected %q to be tagged as unknown", k)
		}
	}
}

func TestTagResult_CollectsAllKeys(t *testing.T) {
	r := makeResult(
		[]string{"DB_PASSWORD"},
		[]string{"ENABLE_BETA"},
		[]differ.Mismatch{{Key: "DATABASE_URL", Left: "a", Right: "b"}},
	)
	res := TagResult(r)
	if _, ok := res.Tags["DB_PASSWORD"]; !ok {
		t.Error("expected DB_PASSWORD in tags")
	}
	if _, ok := res.Tags["ENABLE_BETA"]; !ok {
		t.Error("expected ENABLE_BETA in tags")
	}
	if _, ok := res.Tags["DATABASE_URL"]; !ok {
		t.Error("expected DATABASE_URL in tags")
	}
}

func TestTagResult_SummaryGroupsKeys(t *testing.T) {
	r := makeResult(
		[]string{"DB_PASSWORD", "API_TOKEN"},
		nil, nil,
	)
	res := TagResult(r)
	secrets := res.Summary[TagSecret]
	if len(secrets) != 2 {
		t.Errorf("expected 2 secret keys, got %d", len(secrets))
	}
}

func TestTagKeys_Deprecated(t *testing.T) {
	tags := TagKeys([]string{"OLD_API_URL", "LEGACY_TOKEN"})
	for _, k := range []string{"OLD_API_URL", "LEGACY_TOKEN"} {
		if !containsTag(tags[k], TagDeprecated) {
			t.Errorf("expected %q to be tagged as deprecated", k)
		}
	}
}

func containsTag(tags []Tag, target Tag) bool {
	for _, t := range tags {
		if t == target {
			return true
		}
	}
	return false
}
