package profiler_test

import (
	"testing"

	"github.com/user/envdiff/internal/profiler"
)

func envs() map[string]map[string]string {
	return map[string]map[string]string{
		"dev": {
			"APP_NAME": "myapp",
			"DB_HOST":  "localhost",
			"SECRET":   "dev-secret",
		},
		"staging": {
			"APP_NAME": "myapp",
			"DB_HOST":  "staging.db",
			"SECRET":   "stg-secret",
		},
		"prod": {
			"APP_NAME": "myapp",
			"DB_HOST":  "prod.db",
		},
	}
}

func TestProfile_TotalKeys(t *testing.T) {
	r := profiler.Profile(envs())
	if r.TotalKeys != 3 {
		t.Errorf("expected 3 total keys, got %d", r.TotalKeys)
	}
}

func TestProfile_EnvNamesSorted(t *testing.T) {
	r := profiler.Profile(envs())
	if len(r.EnvNames) != 3 {
		t.Fatalf("expected 3 env names, got %d", len(r.EnvNames))
	}
	if r.EnvNames[0] != "dev" || r.EnvNames[1] != "prod" || r.EnvNames[2] != "staging" {
		t.Errorf("unexpected env order: %v", r.EnvNames)
	}
}

func TestProfile_ConsistentKey(t *testing.T) {
	r := profiler.Profile(envs())
	var appProfile *profiler.KeyProfile
	for i := range r.Profiles {
		if r.Profiles[i].Key == "APP_NAME" {
			appProfile = &r.Profiles[i]
		}
	}
	if appProfile == nil {
		t.Fatal("APP_NAME profile not found")
	}
	if appProfile.UniqueValues != 1 {
		t.Errorf("expected 1 unique value for APP_NAME, got %d", appProfile.UniqueValues)
	}
	if len(appProfile.MissingFrom) != 0 {
		t.Errorf("APP_NAME should not be missing from any env")
	}
}

func TestProfile_MissingKey(t *testing.T) {
	r := profiler.Profile(envs())
	var secretProfile *profiler.KeyProfile
	for i := range r.Profiles {
		if r.Profiles[i].Key == "SECRET" {
			secretProfile = &r.Profiles[i]
		}
	}
	if secretProfile == nil {
		t.Fatal("SECRET profile not found")
	}
	if len(secretProfile.MissingFrom) != 1 || secretProfile.MissingFrom[0] != "prod" {
		t.Errorf("expected SECRET missing from prod, got %v", secretProfile.MissingFrom)
	}
}

func TestProfile_InconsistentCount(t *testing.T) {
	r := profiler.Profile(envs())
	// DB_HOST differs across envs; SECRET missing from prod — both inconsistent
	if r.Inconsistent < 2 {
		t.Errorf("expected at least 2 inconsistent keys, got %d", r.Inconsistent)
	}
}

func TestProfile_EmptyInput(t *testing.T) {
	r := profiler.Profile(map[string]map[string]string{})
	if r.TotalKeys != 0 {
		t.Errorf("expected 0 keys for empty input")
	}
}
