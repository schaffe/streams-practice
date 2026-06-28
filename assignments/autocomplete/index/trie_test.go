package index

import (
	"testing"

	"streams-practice/assignments/autocomplete/corpus"
)

func TestEmptyPrefixReturnsTop10Overall(t *testing.T) {
	entries := []corpus.Entry{
		{Term: "the", Frequency: 1000},
		{Term: "be", Frequency: 900},
		{Term: "and", Frequency: 800},
		{Term: "to", Frequency: 700},
		{Term: "of", Frequency: 600},
		{Term: "a", Frequency: 500},
		{Term: "in", Frequency: 400},
		{Term: "that", Frequency: 300},
		{Term: "have", Frequency: 200},
		{Term: "it", Frequency: 100},
		{Term: "extra", Frequency: 50},
	}

	trie := New(entries)
	results := trie.Search("")

	if len(results) != 10 {
		t.Fatalf("Search(\"\") returned %d results, want 10", len(results))
	}

	for i := 1; i < len(results); i++ {
		if results[i-1].Frequency < results[i].Frequency {
			t.Errorf("results not sorted by frequency desc at index %d: %d < %d",
				i, results[i-1].Frequency, results[i].Frequency)
		}
	}

	expected := []string{"the", "be", "and", "to", "of", "a", "in", "that", "have", "it"}
	for i, e := range expected {
		if results[i].Term != e {
			t.Errorf("result[%d].Term = %q, want %q", i, results[i].Term, e)
		}
	}
}

func TestPrefixMatchReturnsAtMost10(t *testing.T) {
	entries := []corpus.Entry{
		{Term: "apple", Frequency: 100},
		{Term: "application", Frequency: 90},
		{Term: "appetite", Frequency: 80},
		{Term: "appliance", Frequency: 70},
		{Term: "applause", Frequency: 60},
		{Term: "append", Frequency: 50},
		{Term: "appendix", Frequency: 40},
		{Term: "appoint", Frequency: 30},
		{Term: "appreciate", Frequency: 20},
		{Term: "approach", Frequency: 10},
		{Term: "approve", Frequency: 5},
		{Term: "apricot", Frequency: 1},
		{Term: "banana", Frequency: 200},
		{Term: "cherry", Frequency: 150},
	}

	trie := New(entries)
	results := trie.Search("app")

	if len(results) == 0 {
		t.Fatal("Search(\"app\") returned 0 results, want > 0")
	}
	if len(results) > 10 {
		t.Fatalf("Search(\"app\") returned %d results, want <= 10", len(results))
	}

	for i := 1; i < len(results); i++ {
		if results[i-1].Frequency < results[i].Frequency {
			t.Errorf("results not sorted by frequency desc at index %d: %d < %d",
				i, results[i-1].Frequency, results[i].Frequency)
		}
	}

	for _, r := range results {
		if r.Term == "banana" || r.Term == "cherry" {
			t.Errorf("Search(\"app\") returned %q, which does not start with \"app\"", r.Term)
		}
	}
}

func TestNoMatchReturnsNil(t *testing.T) {
	entries := []corpus.Entry{
		{Term: "apple", Frequency: 100},
		{Term: "banana", Frequency: 80},
	}

	trie := New(entries)

	if results := trie.Search("zzz"); results != nil {
		t.Errorf("Search(\"zzz\") = %v, want nil", results)
	}

	if results := trie.Search("xyz"); results != nil {
		t.Errorf("Search(\"xyz\") = %v, want nil", results)
	}
}

func TestCaseInsensitive(t *testing.T) {
	entries := []corpus.Entry{
		{Term: "Apple", Frequency: 100},
		{Term: "apricot", Frequency: 80},
		{Term: "banana", Frequency: 50},
	}

	lower := New(entries).Search("ap")
	upper := New(entries).Search("AP")
	mixed := New(entries).Search("Ap")

	if len(lower) != len(upper) || len(lower) != len(mixed) {
		t.Fatalf("length mismatch: lower=%d, upper=%d, mixed=%d",
			len(lower), len(upper), len(mixed))
	}

	for i := range lower {
		if lower[i].Term != upper[i].Term {
			t.Errorf("lower[%d].Term=%q != upper[%d].Term=%q", i, lower[i].Term, i, upper[i].Term)
		}
		if lower[i].Term != mixed[i].Term {
			t.Errorf("lower[%d].Term=%q != mixed[%d].Term=%q", i, lower[i].Term, i, mixed[i].Term)
		}
	}
}

func TestUnicodeRunes(t *testing.T) {
	entries := []corpus.Entry{
		{Term: "café", Frequency: 90},
		{Term: "cafeteria", Frequency: 80},
		{Term: "cafe", Frequency: 70},
		{Term: "naïve", Frequency: 60},
		{Term: "naivety", Frequency: 50},
		{Term: "über", Frequency: 40},
		{Term: "uber", Frequency: 30},
		{Term: "собор", Frequency: 100},
		{Term: "собака", Frequency: 80},
	}

	t.Run("latin accent", func(t *testing.T) {
		trie := New(entries)
		results := trie.Search("caf")
		if len(results) == 0 {
			t.Fatal("Search(\"caf\") returned 0 results for unicode-containing terms")
		}
	})

	t.Run("cyrillic", func(t *testing.T) {
		trie := New(entries)
		results := trie.Search("со")
		if len(results) == 0 {
			t.Fatal("Search(\"со\") returned 0 results for cyrillic prefix")
		}
		for _, r := range results {
			first := []rune(r.Term)
			if len(first) == 0 || first[0] != 'с' {
				t.Errorf("cyrillic Search result %q does not start with cyrillic 'с'", r.Term)
			}
		}
	})

	t.Run("latin combining char", func(t *testing.T) {
		trie := New(entries)
		results := trie.Search("na")
		if len(results) == 0 {
			t.Fatal("Search(\"na\") returned 0 results for unicode terms with combining chars")
		}
	})
}

func TestSingleCharacterPrefix(t *testing.T) {
	entries := []corpus.Entry{
		{Term: "apple", Frequency: 100},
		{Term: "banana", Frequency: 80},
		{Term: "avocado", Frequency: 70},
		{Term: "cherry", Frequency: 60},
		{Term: "blueberry", Frequency: 50},
	}

	t.Run("prefix a", func(t *testing.T) {
		trie := New(entries)
		results := trie.Search("a")
		if len(results) == 0 {
			t.Fatal("Search(\"a\") returned 0 results")
		}
		if len(results) > 10 {
			t.Fatalf("Search(\"a\") returned %d results, want <= 10", len(results))
		}
		for _, r := range results {
			if r.Term[0] != 'a' && r.Term[0] != 'A' {
				t.Errorf("Search(\"a\") returned %q, which does not start with 'a'", r.Term)
			}
		}
	})

	t.Run("prefix b", func(t *testing.T) {
		trie := New(entries)
		results := trie.Search("b")
		if len(results) == 0 {
			t.Fatal("Search(\"b\") returned 0 results")
		}
		for _, r := range results {
			if r.Term[0] != 'b' && r.Term[0] != 'B' {
				t.Errorf("Search(\"b\") returned %q, which does not start with 'b'", r.Term)
			}
		}
	})
}

func TestWholeWordPrefix(t *testing.T) {
	entries := []corpus.Entry{
		{Term: "apple", Frequency: 100},
		{Term: "application", Frequency: 90},
		{Term: "appetizer", Frequency: 80},
	}

	trie := New(entries)

	results := trie.Search("apple")
	if len(results) == 0 {
		t.Fatal("Search(\"apple\") returned 0 results, term itself should appear")
	}

	found := false
	for _, r := range results {
		if r.Term == "apple" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Search(\"apple\") results = %v, want \"apple\" to be included", results)
	}
}

func TestTrieWithZeroEntries(t *testing.T) {
	trie := New([]corpus.Entry{})

	if results := trie.Search(""); results != nil {
		t.Errorf("Search(\"\") with zero entries = %v, want nil", results)
	}

	if results := trie.Search("a"); results != nil {
		t.Errorf("Search(\"a\") with zero entries = %v, want nil", results)
	}
}

func TestTrieWithThreeEntries(t *testing.T) {
	entries := []corpus.Entry{
		{Term: "apple", Frequency: 100},
		{Term: "banana", Frequency: 80},
		{Term: "cherry", Frequency: 60},
	}

	trie := New(entries)
	results := trie.Search("")

	if len(results) != 3 {
		t.Fatalf("Search(\"\") returned %d results, want 3", len(results))
	}

	for i := 1; i < len(results); i++ {
		if results[i-1].Frequency < results[i].Frequency {
			t.Errorf("results not sorted by frequency desc at index %d: %d < %d",
				i, results[i-1].Frequency, results[i].Frequency)
		}
	}
}

func TestSameFrequencyStableOrder(t *testing.T) {
	entries := []corpus.Entry{
		{Term: "delta", Frequency: 50},
		{Term: "gamma", Frequency: 50},
		{Term: "beta", Frequency: 50},
		{Term: "alpha", Frequency: 50},
	}

	trie := New(entries)
	results := trie.Search("")

	if len(results) == 0 {
		t.Fatal("Search(\"\") returned 0 results")
	}

	for i := 1; i < len(results); i++ {
		if results[i-1].Frequency < results[i].Frequency {
			t.Errorf("results not sorted by frequency desc at index %d: %d < %d",
				i, results[i-1].Frequency, results[i].Frequency)
		}
	}
}

func TestNodeWithFewerThan10EntriesReturnsAll(t *testing.T) {
	entries := []corpus.Entry{
		{Term: "xyzzy", Frequency: 100},
		{Term: "xylophone", Frequency: 80},
		{Term: "xenon", Frequency: 60},
	}

	t.Run("prefix xy", func(t *testing.T) {
		trie := New(entries)
		results := trie.Search("xy")
		if len(results) != 2 {
			t.Fatalf("Search(\"xy\") returned %d results, want 2", len(results))
		}
	})

	t.Run("prefix x", func(t *testing.T) {
		trie := New(entries)
		results := trie.Search("x")
		if len(results) != 3 {
			t.Fatalf("Search(\"x\") returned %d results, want 3", len(results))
		}
	})
}
