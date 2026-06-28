package ui

import (
	"bytes"
	"strings"
	"testing"

	"streams-practice/assignments/rover/entity"
	"streams-practice/assignments/rover/world"
)

func newTestCLI(input string) (*CLI, *bytes.Buffer) {
	rover := entity.NewRover(ActiveUnit, 0, 0, entity.North)
	w := world.New(rover)
	out := &bytes.Buffer{}
	r := NewLineRenderer()
	cli := NewCLI(w, r, strings.NewReader(input), out)
	return cli, out
}

// TestDriveCommandSequence verifies that a sequence of commands moves the
// rover and produces the expected rendered status line.
func TestDriveCommandSequence(t *testing.T) {
	// "M" moves north to (0,1), "R" turns to face East → final: (0,1) E
	cli, out := newTestCLI("MRq")

	err := cli.drive()
	if err != nil {
		t.Fatalf("drive() returned unexpected error: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "(0, 1) E") {
		t.Errorf("expected output to contain '(0, 1) E', got: %q", output)
	}
}

// TestDriveQuitStopsLoop verifies that 'q' causes drive to return nil and
// bytes after 'q' are never applied.
func TestDriveQuitStopsLoop(t *testing.T) {
	// After q, the "M" should never be applied.
	cli, out := newTestCLI("qM")

	err := cli.drive()
	if err != nil {
		t.Fatalf("drive() returned unexpected error: %v", err)
	}

	// No render should have happened before 'q' because no valid command preceded it.
	// The output should NOT contain a moved position; starting position is (0,0) N.
	output := out.String()
	if strings.Contains(output, "(0, 1)") {
		t.Errorf("expected 'M' after 'q' to be ignored, but output shows movement: %q", output)
	}
}

// TestDriveCtrlCStopsLoop verifies that 0x03 (Ctrl-C) causes drive to return nil.
func TestDriveCtrlCStopsLoop(t *testing.T) {
	// After Ctrl-C, the "M" should never be applied.
	cli, out := newTestCLI("\x03M")

	err := cli.drive()
	if err != nil {
		t.Fatalf("drive() returned unexpected error: %v", err)
	}

	output := out.String()
	if strings.Contains(output, "(0, 1)") {
		t.Errorf("expected 'M' after Ctrl-C to be ignored, but output shows movement: %q", output)
	}
}

// TestDriveUppercaseQuitStopsLoop verifies that 'Q' also stops the loop.
func TestDriveUppercaseQuitStopsLoop(t *testing.T) {
	cli, _ := newTestCLI("Q")

	err := cli.drive()
	if err != nil {
		t.Fatalf("drive() returned unexpected error on Q: %v", err)
	}
}
