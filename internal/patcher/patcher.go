// Package patcher provides functionality to apply a patch (set of key-value
// pairs) to an existing .env file, adding missing keys and optionally
// updating mismatched ones.
package patcher

import (
	"fmt"
	"os"\n	"strings"

	"github.com/user/envdiff/internal/differ"
)

// Options controls the behaviour of the Patch operation.
type Options struct {
	// UpdateMismatched, when true, overwrites keys whose values differ.
	UpdateMismatched bool
	// DryRun, when true, returns the patched content without writing to disk.
	DryRun bool
}

// Result describes what the patcher did (or would do).
type Result struct {
	Added   []string
	Updated []string
	Skipped []string
}

// Patch reads destPath, applies keys from patch map according to opts, and
// writes the result back to destPath (unless DryRun is set).
// It returns the final file content and a Result summary.
func Patch(destPath string, patch map[string]string, diff differ.Result, opts Options) (string, Result, error) {
	raw, err := os.ReadFile(destPath)
	if err != nil {
		return "", Result{}, fmt.Errorf("patcher: read %s: %w", destPath, err)
	}

	lines := strings.Split(string(raw), "\n")
	existing := map[string]int{} // key -> line index
	for i, line := range lines {
		if idx := strings.IndexByte(line, '='); idx > 0 {
			key := strings.TrimSpace(line[:idx])
			if key != "" && !strings.HasPrefix(key, "#") {
				existing[key] = i
			}
		}
	}

	var res Result

	// Update mismatched keys in-place.
	if opts.UpdateMismatched {
		for _, mm := range diff.Mismatched {
			if val, ok := patch[mm.Key]; ok {
				if idx, found := existing[mm.Key]; found {
					lines[idx] = mm.Key + "=" + val
					res.Updated = append(res.Updated, mm.Key)
				}
			} else {
				res.Skipped = append(res.Skipped, mm.Key)
			}
		}
	} else {
		for _, mm := range diff.Mismatched {
			res.Skipped = append(res.Skipped, mm.Key)
		}
	}

	// Append keys missing in dest.
	for _, key := range diff.MissingInRight {
		if val, ok := patch[key]; ok {
			lines = append(lines, key+"="+val)
			res.Added = append(res.Added, key)
		}
	}

	output := strings.Join(lines, "\n")

	if !opts.DryRun {
		if err := os.WriteFile(destPath, []byte(output), 0o644); err != nil {
			return "", res, fmt.Errorf("patcher: write %s: %w", destPath, err)
		}
	}

	return output, res, nil
}
