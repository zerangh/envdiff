// Package parser provides utilities for reading and parsing .env files
// into structured key-value maps.
//
// A .env file is expected to follow the common KEY=VALUE format:
//
//	# This is a comment
//	APP_ENV=production
//	DATABASE_URL="postgres://localhost:5432/mydb"
//	SECRET_KEY='s3cr3t'
//
// Rules:
//   - Lines beginning with '#' are treated as comments and ignored.
//   - Blank lines are ignored.
//   - Values may optionally be wrapped in single or double quotes; the quotes
//     are stripped from the parsed value.
//   - Each non-comment, non-blank line must contain exactly one '=' separator.
//   - Keys must not be empty.
//
// Example usage:
//
//	env, err := parser.ParseFile(".env.production")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(env["APP_ENV"]) // "production"
package parser
