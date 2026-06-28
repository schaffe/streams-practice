// Package cli provides an interactive command-line interface for autocomplete search.
package cli

import (
	"fmt"
	"io"
	"os"
	"unicode/utf8"

	"golang.org/x/term"

	"streams-practice/assignments/autocomplete/corpus"
)

// Searcher defines the interface for prefix-based search.
type Searcher interface {
	Search(prefix string) []corpus.Entry
}

// CLI handles user interaction for autocomplete search.
type CLI struct {
	searcher Searcher
	in       io.Reader
	out      io.Writer
}

// NewCLI creates a new CLI with the given searcher, input, and output.
func NewCLI(s Searcher, in io.Reader, out io.Writer) *CLI {
	return &CLI{searcher: s, in: in, out: out}
}

// RunBatch runs a batch search for the given prefix and prints results.
func (c *CLI) RunBatch(prefix string) error {
	results := c.searcher.Search(prefix)
	for _, e := range results {
		if _, err := fmt.Fprintf(c.out, "%s (%d)\n", e.Term, e.Frequency); err != nil {
			return err
		}
	}
	return nil
}

// RunInteractive runs the interactive terminal UI.
func (c *CLI) RunInteractive() error {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return err
	}
	defer func() { _ = term.Restore(fd, oldState) }()
	return c.drive()
}

func (c *CLI) drive() error {
	var prefix string
	var results []corpus.Entry
	c.render("", nil)

	var buf [4]byte
	bufLen := 0

	for {
		if bufLen >= len(buf) {
			bufLen = 0
			continue
		}
		_, err := c.in.Read(buf[bufLen : bufLen+1])
		if err != nil {
			return err
		}
		bufLen++

		r, size := utf8.DecodeRune(buf[:bufLen])
		if r == utf8.RuneError && size == 1 {
			continue
		}

		bufLen = 0

		switch {
		case r == 0x03 || r == 0x51 || r == 0x71:
			return nil
		case r >= 0x20 && r <= 0x7e:
			prefix += string(r)
			results = c.searcher.Search(prefix)
			c.render(prefix, results)
		case r == 0x7f || r == 0x08:
			if len(prefix) > 0 {
				runes := []rune(prefix)
				prefix = string(runes[:len(runes)-1])
				results = c.searcher.Search(prefix)
				c.render(prefix, results)
			}
		case r == 0x09:
			if len(results) > 0 {
				prefix = results[0].Term
				results = c.searcher.Search(prefix)
				c.render(prefix, results)
			}
		case r == 0x0d || r == 0x0a:
		default:
		}
	}
}

func (c *CLI) render(prefix string, results []corpus.Entry) {
	_, _ = fmt.Fprintf(c.out, "\r\033[J")
	if prefix != "" {
		_, _ = fmt.Fprintf(c.out, "> %s\n", prefix)
	}
	for i, e := range results {
		_, _ = fmt.Fprintf(c.out, "  %d. %s (%d)\n", i+1, e.Term, e.Frequency)
	}
}
