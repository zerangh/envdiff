package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// EnvMap represents a parsed .env file as a map of key-value pairs.
type EnvMap map[string]string

// ParseFile reads a .env file and returns an EnvMap.
// It skips blank lines and comments (lines starting with '#').
// It returns an error if the file cannot be opened or contains malformed lines.
func ParseFile(path string) (EnvMap, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("parser: could not open file %q: %w", path, err)
	}
	defer f.Close()

	env := make(EnvMap)
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip blank lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("parser: %s line %d: %w", path, lineNum, err)
		}

		env[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("parser: error reading %q: %w", path, err)
	}

	return env, nil
}

// parseLine splits a single KEY=VALUE line into its components.
func parseLine(line string) (string, string, error) {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("malformed line %q: missing '='"  , line)
	}

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	if key == "" {
		return "", "", fmt.Errorf("malformed line %q: empty key", line)
	}

	// Strip optional surrounding quotes from value
	value = stripQuotes(value)

	return key, value, nil
}

// stripQuotes removes matching surrounding single or double quotes from a string.
func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
