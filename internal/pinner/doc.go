// Package pinner checks that a set of required (pinned) environment variable
// keys are present and non-empty across one or more parsed env maps.
//
// Typical usage:
//
//	pins := []string{"DATABASE_URL", "SECRET_KEY", "API_TOKEN"}
//
//	// Check a single environment.
//	report := pinner.CheckEnv("production", envMap, pins)
//	if !report.OK() {
//		fmt.Print(report.Format())
//	}
//
//	// Check multiple environments at once.
//	report = pinner.CheckAll(allEnvMaps, pins)
//
The Report type exposes Missing and Empty slices so callers can
programmatically act on individual failures in addition to using
the human-readable Format() helper.
package pinner
