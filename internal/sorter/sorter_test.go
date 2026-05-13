package sorter_test

import (
	"testing"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/sorter"
)

func makeResult() differ.Result {
	return differ.Result{
		MissingInRight: []string{"DB_PORT", "AWS_SECRET", "APP_NAME"},
		MissingInLeft:  []string{"DB_HOST", "AWS_KEY"},
		Mismatched: []differ.Mismatch{
			{Key: "DB_PASS", LeftVal: "x", RightVal: "y"},
			{Key: "APP_ENV", LeftVal: "dev", RightVal: "prod"},
		},
	}
}

func TestSortResult_Alpha(t *testing.T) {
	r := makeResult()
	out := sorter.SortResult(r, sorter.Alpha)

	if out.MissingInRight[0] != "APP_NAME" {
		t.Errorf("expected APP_NAME first, got %s", out.MissingInRight[0])
	}
	if out.MissingInLeft[0] != "AWS_KEY" {
		t.Errorf("expected AWS_KEY first, got %s", out.MissingInLeft[0])
	}
	if out.Mismatched[0].Key != "APP_ENV" {
		t.Errorf("expected APP_ENV first, got %s", out.Mismatched[0].Key)
	}
}

func TestSortResult_ChangeType_DoesNotMix(t *testing.T) {
	r := makeResult()
	out := sorter.SortResult(r, sorter.ChangeType)

	// Each bucket should still be internally sorted
	for i := 1; i < len(out.MissingInRight); i++ {
		if out.MissingInRight[i] < out.MissingInRight[i-1] {
			t.Errorf("MissingInRight not sorted at index %d", i)
		}
	}
	for i := 1; i < len(out.Mismatched); i++ {
		if out.Mismatched[i].Key < out.Mismatched[i-1].Key {
			t.Errorf("Mismatched not sorted at index %d", i)
		}
	}
}

func TestSortResult_Prefix(t *testing.T) {
	r := makeResult()
	out := sorter.SortResult(r, sorter.Prefix)

	// APP_ prefix should come before AWS_ and DB_
	if out.MissingInRight[0] != "APP_NAME" {
		t.Errorf("expected APP_NAME first under prefix sort, got %s", out.MissingInRight[0])
	}
}

func TestSortResult_DoesNotMutateOriginal(t *testing.T) {
	r := makeResult()
	origFirst := r.MissingInRight[0]
	sorter.SortResult(r, sorter.Alpha)
	if r.MissingInRight[0] != origFirst {
		t.Error("SortResult mutated the original result")
	}
}

func TestSortedEnv_Alpha(t *testing.T) {
	env := map[string]string{"Z_KEY": "1", "A_KEY": "2", "M_KEY": "3"}
	keys := sorter.SortedEnv(env, sorter.Alpha)
	if keys[0] != "A_KEY" || keys[1] != "M_KEY" || keys[2] != "Z_KEY" {
		t.Errorf("unexpected order: %v", keys)
	}
}

func TestSortedEnv_Prefix(t *testing.T) {
	env := map[string]string{"DB_HOST": "h", "APP_ENV": "e", "DB_PORT": "p", "APP_NAME": "n"}
	keys := sorter.SortedEnv(env, sorter.Prefix)
	// APP_ keys should precede DB_ keys
	appDone := false
	for _, k := range keys {
		if k[:3] == "DB_" {
			appDone = true
		}
		if appDone && k[:3] == "APP" {
			t.Errorf("APP_ key appeared after DB_ key: %s", k)
		}
	}
}
