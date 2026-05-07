package reporter

import (
	"encoding/json"
	"io"

	"github.com/user/envdiff/internal/differ"
)

type jsonMismatch struct {
	Key        string `json:"key"`
	LeftValue  string `json:"left_value"`
	RightValue string `json:"right_value"`
}

type jsonReport struct {
	LeftFile      string         `json:"left_file"`
	RightFile     string         `json:"right_file"`
	MissingInLeft  []string       `json:"missing_in_left"`
	MissingInRight []string       `json:"missing_in_right"`
	Mismatched    []jsonMismatch `json:"mismatched"`
	Clean         bool           `json:"clean"`
}

func reportJSON(w io.Writer, result differ.Result, leftName, rightName string) error {
	mismatches := make([]jsonMismatch, 0, len(result.Mismatched))
	for _, m := range result.Mismatched {
		mismatches = append(mismatches, jsonMismatch{
			Key:        m.Key,
			LeftValue:  m.LeftValue,
			RightValue: m.RightValue,
		})
	}

	missingLeft := result.MissingInLeft
	if missingLeft == nil {
		missingLeft = []string{}
	}
	missingRight := result.MissingInRight
	if missingRight == nil {
		missingRight = []string{}
	}

	report := jsonReport{
		LeftFile:      leftName,
		RightFile:     rightName,
		MissingInLeft:  missingLeft,
		MissingInRight: missingRight,
		Mismatched:    mismatches,
		Clean:         result.IsClean(),
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}
