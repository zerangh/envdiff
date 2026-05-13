// Package redactor provides utilities for identifying and masking sensitive
// environment variable values before they are displayed, logged, or exported.
//
// Keys are considered sensitive when their names match well-known patterns such
// as PASSWORD, SECRET, TOKEN, API_KEY, PRIVATE, AUTH, or CREDENTIALS.
// Callers may supply additional regexp patterns via Options.ExtraPatterns.
//
// Two primary functions are provided:
//
//   - RedactEnv masks sensitive values in a raw key→value map.
//   - RedactResult masks sensitive mismatch values inside a differ.Result so
//     that diff reports never leak credentials.
//
// Example:
//
//	env := map[string]string{"DB_PASSWORD": "hunter2", "PORT": "5432"}
//	safe := redactor.RedactEnv(env, redactor.Options{})
//	// safe["DB_PASSWORD"] == "***"
//	// safe["PORT"]        == "5432"
package redactor
