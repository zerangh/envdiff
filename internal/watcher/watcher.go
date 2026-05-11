// Package watcher provides functionality for watching .env files for changes
// and reporting a diff summary when modifications are detected.
package watcher

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/parser"
)

// Event represents a change detected between two .env file states.
type Event struct {
	File      string
	ChangedAt time.Time
	Result    differ.Result
}

// Options configures the watcher behaviour.
type Options struct {
	// PollInterval is how often the file is checked for changes.
	PollInterval time.Duration
}

var defaultOptions = Options{
	PollInterval: 2 * time.Second,
}

// Watch polls the given file for changes relative to a reference file.
// When a change is detected it sends an Event to the returned channel.
// Call the returned stop function to halt watching.
func Watch(referenceFile, watchedFile string, opts *Options) (<-chan Event, func(), error) {
	if opts == nil {
		opts = &defaultOptions
	}

	if _, err := os.Stat(referenceFile); err != nil {
		return nil, nil, fmt.Errorf("reference file: %w", err)
	}
	if _, err := os.Stat(watchedFile); err != nil {
		return nil, nil, fmt.Errorf("watched file: %w", err)
	}

	lastHash, err := fileHash(watchedFile)
	if err != nil {
		return nil, nil, fmt.Errorf("initial hash: %w", err)
	}

	events := make(chan Event, 1)
	stopCh := make(chan struct{})

	go func() {
		ticker := time.NewTicker(opts.PollInterval)
		defer ticker.Stop()
		defer close(events)

		for {
			select {
			case <-stopCh:
				return
			case <-ticker.C:
				currentHash, err := fileHash(watchedFile)
				if err != nil || currentHash == lastHash {
					continue
				}
				lastHash = currentHash

				ref, err := parser.ParseFile(referenceFile)
				if err != nil {
					continue
				}
				watched, err := parser.ParseFile(watchedFile)
				if err != nil {
					continue
				}

				result := differ.Diff(ref, watched)
				events <- Event{
					File:      watchedFile,
					ChangedAt: time.Now(),
					Result:    result,
				}
			}
		}
	}()

	stop := func() { close(stopCh) }
	return events, stop, nil
}

// fileHash returns a SHA-256 hex digest of the file contents.
func fileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
