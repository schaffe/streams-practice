// Package ui provides a terminal CLI for controlling a rover world.
package ui

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"

	"golang.org/x/term"

	"streams-practice/assignments/rover/world"
)

// CLI drives a World using a Renderer, reading commands from in and writing output to out.
type CLI struct {
	world    World
	renderer Renderer
	in       io.Reader
	out      io.Writer
}

// NewCLI constructs a CLI.
func NewCLI(w World, r Renderer, in io.Reader, out io.Writer) *CLI {
	return &CLI{world: w, renderer: r, in: in, out: out}
}

// parse maps a rune to a Command, a quit signal, or unknown.
// Returns (cmd, ok, quit): quit=true means stop; ok=false means unrecognized.
func parse(r rune) (world.Command, bool, bool) {
	switch unicode.ToUpper(r) {
	case 'L':
		return world.Command{Unit: ActiveUnit, Action: world.TurnLeft}, true, false
	case 'R':
		return world.Command{Unit: ActiveUnit, Action: world.TurnRight}, true, false
	case 'M':
		return world.Command{Unit: ActiveUnit, Action: world.Forward}, true, false
	case 'Q':
		return world.Command{}, false, true
	default:
		return world.Command{}, false, false
	}
}

// RunBatch executes a string of commands and renders the final state.
// Unrecognized characters return an error. 'q'/'Q' are ignored in batch mode.
func (c *CLI) RunBatch(commands string) error {
	for _, r := range commands {
		cmd, ok, quit := parse(r)
		if quit {
			// quit chars ignored in batch
			continue
		}
		if !ok {
			return fmt.Errorf("unrecognized command %q", r)
		}
		if err := c.world.Apply(cmd); err != nil {
			return err
		}
	}
	c.renderer.Render(c.out, c.world.Snapshot())
	return nil
}

// drive is the terminal-free read-parse-apply-render loop.
// It reads bytes from c.in, applies recognized commands, re-renders the
// status line to c.out (with a leading \r), and returns nil on q/Q or Ctrl-C (0x03).
// Unrecognized characters are silently ignored.
func (c *CLI) drive() error {
	br := bufio.NewReader(c.in)
	for {
		b, err := br.ReadByte()
		if err != nil {
			return err
		}
		// Ctrl-C in raw mode arrives as 0x03.
		if b == 0x03 {
			return nil
		}
		r := rune(b)
		cmd, ok, quit := parse(r)
		if quit {
			return nil
		}
		if !ok {
			// Ignore unrecognized characters in interactive mode.
			continue
		}
		if applyErr := c.world.Apply(cmd); applyErr != nil {
			// Show error on the line; keep looping.
			_, _ = fmt.Fprintf(c.out, "\r%v", applyErr)
			continue
		}
		// Re-render the status line in place using carriage return.
		_, _ = fmt.Fprintf(c.out, "\r")
		c.renderer.Render(c.out, c.world.Snapshot())
	}
}

// RunInteractive runs an interactive raw-mode terminal loop.
// Each keystroke applies a command and re-renders the status line in place.
func (c *CLI) RunInteractive() error {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return err
	}
	defer func() { _ = term.Restore(fd, oldState) }()
	return c.drive()
}
