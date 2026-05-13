// Package annotator produces human-readable, severity-tagged annotations
// for the differences found between two .env environments.
//
// Given a differ.Result, Annotate returns a set of Annotation values, each
// describing a specific key-level issue with a recommended action and a
// severity level: info, warning, or critical.
//
// Severity levels:
//
//	- info     : key exists only in the right env (possible new addition)
//	- warning  : key is missing in right env, or values differ
//	- critical : one side has an empty value where the other does not
//
// Example usage:
//
//	result := differ.Diff(leftEnv, rightEnv)
//	annotations := annotator.Annotate(result)
//	fmt.Println(annotations.Format())
package annotator
