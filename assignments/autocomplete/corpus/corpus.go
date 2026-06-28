package corpus

import (
	"encoding/json"
	"os"
)

type Entry struct {
	Term      string `json:"term"`
	Frequency int    `json:"frequency"`
}

type Loader interface {
	Load() ([]Entry, error)
}

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
