// Package corpus provides types and utilities for loading autocomplete corpus data.
package corpus

import (
	"encoding/json"
	"os"
)

// Entry represents a single autocomplete entry with a term and its frequency.
type Entry struct {
	Term      string `json:"term"`
	Frequency int    `json:"frequency"`
}

// Loader defines the interface for loading corpus entries.
type Loader interface {
	Load() ([]Entry, error)
}

// Load reads and parses a JSON corpus file into entries.
func Load(path string) ([]Entry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}

	return entries, nil
}
