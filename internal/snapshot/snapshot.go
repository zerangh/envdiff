// Package snapshot provides functionality for saving and loading
// .env diff results to/from disk for later comparison or auditing.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/yourusername/envdiff/internal/differ"
)

// Record wraps a diff result with metadata for persistence.
type Record struct {
	CreatedAt time.Time        `json:"created_at"`
	LeftFile  string           `json:"left_file"`
	RightFile string           `json:"right_file"`
	Result    differ.Result    `json:"result"`
}

// Save writes the diff result and metadata to a JSON file at the given path.
func Save(path, leftFile, rightFile string, result differ.Result) error {
	rec := Record{
		CreatedAt: time.Now().UTC(),
		LeftFile:  leftFile,
		RightFile: rightFile,
		Result:    result,
	}

	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal failed: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("snapshot: write failed: %w", err)
	}

	return nil
}

// Load reads a snapshot Record from a JSON file at the given path.
func Load(path string) (Record, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Record{}, fmt.Errorf("snapshot: read failed: %w", err)
	}

	var rec Record
	if err := json.Unmarshal(data, &rec); err != nil {
		return Record{}, fmt.Errorf("snapshot: unmarshal failed: %w", err)
	}

	return rec, nil
}
