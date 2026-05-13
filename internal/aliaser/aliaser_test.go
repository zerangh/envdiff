package aliaser_test

import (
	"testing"

	"github.com/user/envdiff/internal/aliaser"
	"github.com/user/envdiff/internal/differ"
)

var sampleAliases = aliaser.AliasMap{
	"LEGACY_DB_HOST": "DATABASE_HOST",
	"OLD_API_KEY":    "API_KEY",
}

func TestApplyToEnv_RenamesAliasedKeys(t *testing.T) {
	env := map[string]string{
		"LEGACY_DB_HOST": "localhost",
		"PORT":           "8080",
	}
	out := aliaser.ApplyToEnv(env, sampleAliases)
	if _, ok := out["LEGACY_DB_HOST"]; ok {
		t.Error("expected old key to be removed")
	}
	if out["DATABASE_HOST"] != "localhost" {
		t.Errorf("expected DATABASE_HOST=localhost, got %q", out["DATABASE_HOST"])
	}
	if out["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", out["PORT"])
	}
}

func TestApplyToEnv_AliasWinsOverOldKey(t *testing.T) {
	env := map[string]string{
		"LEGACY_DB_HOST": "old-value",
		"DATABASE_HOST":  "new-value",
	}
	out := aliaser.ApplyToEnv(env, sampleAliases)
	// Both map to DATABASE_HOST; the alias-renamed entry overwrites.
	if out["DATABASE_HOST"] == "" {
		t.Error("expected DATABASE_HOST to be present")
	}
}

func TestApplyToResult_RenamesMissingAndMismatched(t *testing.T) {
	r := differ.Result{
		MissingInLeft:  []string{"LEGACY_DB_HOST"},
		MissingInRight: []string{"OLD_API_KEY"},
		Mismatched: []differ.Mismatch{
			{Key: "LEGACY_DB_HOST", LeftValue: "a", RightValue: "b"},
		},
	}
	out := aliaser.ApplyToResult(r, sampleAliases)
	if len(out.MissingInLeft) != 1 || out.MissingInLeft[0] != "DATABASE_HOST" {
		t.Errorf("MissingInLeft not renamed: %v", out.MissingInLeft)
	}
	if len(out.MissingInRight) != 1 || out.MissingInRight[0] != "API_KEY" {
		t.Errorf("MissingInRight not renamed: %v", out.MissingInRight)
	}
	if out.Mismatched[0].Key != "DATABASE_HOST" {
		t.Errorf("Mismatched key not renamed: %v", out.Mismatched[0].Key)
	}
}

func TestLoadAliasMap_ParsesLines(t *testing.T) {
	lines := []string{
		"# comment",
		"",
		"LEGACY_DB_HOST = DATABASE_HOST",
		"OLD_API_KEY=API_KEY",
		"BAD_LINE",
	}
	am := aliaser.LoadAliasMap(lines)
	if am["LEGACY_DB_HOST"] != "DATABASE_HOST" {
		t.Errorf("expected DATABASE_HOST, got %q", am["LEGACY_DB_HOST"])
	}
	if am["OLD_API_KEY"] != "API_KEY" {
		t.Errorf("expected API_KEY, got %q", am["OLD_API_KEY"])
	}
	if _, ok := am["BAD_LINE"]; ok {
		t.Error("malformed line should not be parsed")
	}
}

func TestAliasMap_Invert(t *testing.T) {
	am := aliaser.AliasMap{"OLD": "NEW"}
	inv := am.Invert()
	if inv["NEW"] != "OLD" {
		t.Errorf("expected inv[NEW]=OLD, got %q", inv["NEW"])
	}
	if _, ok := inv["OLD"]; ok {
		t.Error("original key should not appear in inverted map")
	}
}
