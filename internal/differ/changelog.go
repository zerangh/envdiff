package differ

import "time"

// ChangeKind describes the type of change recorded in a Changelog entry.
type ChangeKind string

const (
	ChangeAdded    ChangeKind = "added"
	ChangeRemoved  ChangeKind = "removed"
	ChangeMismatch ChangeKind = "mismatch"
)

// ChangeEntry represents a single change between two environments.
type ChangeEntry struct {
	Key        string     `json:"key"`
	Kind       ChangeKind `json:"kind"`
	LeftValue  string     `json:"left_value,omitempty"`
	RightValue string     `json:"right_value,omitempty"`
	RecordedAt time.Time  `json:"recorded_at"`
}

// Changelog is an ordered list of changes derived from a diff Result.
type Changelog struct {
	Entries    []ChangeEntry `json:"entries"`
	LeftFile   string        `json:"left_file"`
	RightFile  string        `json:"right_file"`
	GeneratedAt time.Time   `json:"generated_at"`
}

// BuildChangelog converts a diff Result into a Changelog for auditing or export.
func BuildChangelog(r Result) Changelog {
	now := time.Now().UTC()
	cl := Changelog{
		LeftFile:    r.LeftFile,
		RightFile:   r.RightFile,
		GeneratedAt: now,
	}

	for _, k := range r.MissingInRight {
		cl.Entries = append(cl.Entries, ChangeEntry{
			Key:        k,
			Kind:       ChangeRemoved,
			LeftValue:  r.LeftEnv[k],
			RecordedAt: now,
		})
	}

	for _, k := range r.MissingInLeft {
		cl.Entries = append(cl.Entries, ChangeEntry{
			Key:        k,
			Kind:       ChangeAdded,
			RightValue: r.RightEnv[k],
			RecordedAt: now,
		})
	}

	for _, m := range r.Mismatched {
		cl.Entries = append(cl.Entries, ChangeEntry{
			Key:        m.Key,
			Kind:       ChangeMismatch,
			LeftValue:  m.LeftValue,
			RightValue: m.RightValue,
			RecordedAt: now,
		})
	}

	return cl
}
