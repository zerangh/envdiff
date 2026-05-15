package encoder

import (
	"fmt"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/differ"
)

// Format represents a supported encoding format.
type Format string

const (
	FormatDotenv  Format = "dotenv"
	FormatExports Format = "exports"
	FormatDocker  Format = "docker"
)

// Options controls encoding behaviour.
type Options struct {
	Format      Format
	OmitEmpty   bool
	QuoteAll    bool
	SortKeys    bool
}

// DefaultOptions returns sensible encoding defaults.
func DefaultOptions() Options {
	return Options{
		Format:   FormatDotenv,
		SortKeys: true,
	}
}

// Encode serialises an env map to a string using the requested format.
func Encode(env map[string]string, opts Options) (string, error) {
	switch opts.Format {
	case FormatDotenv:
		return encodeDotenv(env, opts), nil
	case FormatExports:
		return encodeExports(env, opts), nil
	case FormatDocker:
		return encodeDocker(env, opts), nil
	default:
		return "", fmt.Errorf("encoder: unknown format %q", opts.Format)
	}
}

// EncodeResult encodes only the keys present in a differ.Result's left-hand
// environment, applying the same formatting rules.
func EncodeResult(result differ.Result, env map[string]string, opts Options) (string, error) {
	subset := make(map[string]string, len(env))
	for k, v := range env {
		subset[k] = v
	}
	return Encode(subset, opts)
}

func sortedKeys(env map[string]string) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func encodeDotenv(env map[string]string, opts Options) string {
	var sb strings.Builder
	keys := sortedKeys(env)
	for _, k := range keys {
		v := env[k]
		if opts.OmitEmpty && v == "" {
			continue
		}
		if opts.QuoteAll || strings.ContainsAny(v, " \t#") {
			v = fmt.Sprintf("%q", v)
		}
		fmt.Fprintf(&sb, "%s=%s\n", k, v)
	}
	return sb.String()
}

func encodeExports(env map[string]string, opts Options) string {
	var sb strings.Builder
	keys := sortedKeys(env)
	for _, k := range keys {
		v := env[k]
		if opts.OmitEmpty && v == "" {
			continue
		}
		if opts.QuoteAll || strings.ContainsAny(v, " \t#") {
			v = fmt.Sprintf("%q", v)
		}
		fmt.Fprintf(&sb, "export %s=%s\n", k, v)
	}
	return sb.String()
}

func encodeDocker(env map[string]string, opts Options) string {
	var sb strings.Builder
	keys := sortedKeys(env)
	for _, k := range keys {
		v := env[k]
		if opts.OmitEmpty && v == "" {
			continue
		}
		fmt.Fprintf(&sb, "%s=%s\n", k, v)
	}
	return sb.String()
}
