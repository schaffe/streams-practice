package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

// TestSmokeHappyPath runs the rover CLI with a valid command and asserts success
func TestSmokeHappyPath(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "-cmd", "MMRMM")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("expected exit code 0, got error: %v\nstderr: %s", err, stderr.String())
	}

	output := strings.TrimSpace(stdout.String())
	if !strings.Contains(output, "(2, 2) E") {
		t.Errorf("expected stdout to contain '(2, 2) E', got: %s", output)
	}
}

// TestSmokeErrorPath runs the rover CLI with an invalid command and asserts error
func TestSmokeErrorPath(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "-cmd", "X")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err == nil {
		t.Fatal("expected non-zero exit code, but command succeeded")
	}

	// Verify it's an exit error with non-zero code
	if _, ok := err.(*exec.ExitError); !ok {
		t.Fatalf("expected exec.ExitError, got: %T", err)
	}

	// Verify stderr is non-empty and mentions the unrecognized command
	stderrStr := stderr.String()
	if stderrStr == "" {
		t.Error("expected non-empty stderr, got empty string")
	}
	if !strings.Contains(stderrStr, "X") {
		t.Errorf("expected stderr to contain 'X', got: %q", stderrStr)
	}
}
