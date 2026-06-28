package corpus

import (
	"testing"
)

func TestLoadValidJSON(t *testing.T) {
	entries, err := Load("testdata/valid.json")
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("Load() returned %d entries, want 3", len(entries))
	}
	want := []Entry{
		{Term: "foo", Frequency: 42},
		{Term: "bar", Frequency: 17},
		{Term: "baz", Frequency: 99},
	}
	for i, e := range entries {
		if e != want[i] {
			t.Errorf("entries[%d] = %+v, want %+v", i, e, want[i])
		}
	}
}

func TestLoadEmptyArray(t *testing.T) {
	entries, err := Load("testdata/empty.json")
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("Load() returned %d entries, want 0", len(entries))
	}
}

func TestLoadMalformedJSON(t *testing.T) {
	_, err := Load("testdata/malformed.json")
	if err == nil {
		t.Fatal("Load() expected error for malformed JSON")
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := Load("testdata/nonexistent.json")
	if err == nil {
		t.Fatal("Load() expected error for missing file")
	}
}

func TestLoadNonArrayJSON(t *testing.T) {
	_, err := Load("testdata/nonarray.json")
	if err == nil {
		t.Fatal("Load() expected error for non-array JSON")
	}
}
