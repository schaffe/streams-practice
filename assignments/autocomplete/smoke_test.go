package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestSmokePrefixThe(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "-prefix", "the")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("expected exit code 0, got error: %v\nstderr: %s", err, stderr.String())
	}

	output := stdout.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) == 0 || (len(lines) == 1 && lines[0] == "") {
		t.Fatal("expected non-empty output")
	}
	if len(lines) > 10 {
		t.Errorf("expected ≤10 lines, got %d", len(lines))
	}
	if !strings.Contains(output, "the (1000)") {
		t.Errorf("expected stdout to contain 'the (1000)', got: %s", output)
	}
}

func TestSmokeEmptyPrefix(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "-prefix", "")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("expected exit code 0, got error: %v\nstderr: %s", err, stderr.String())
	}

	lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	if len(lines) != 10 {
		t.Errorf("expected exactly 10 lines, got %d", len(lines))
	}
}

func TestSmokePrefixZzz(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "-prefix", "zzz")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("expected exit code 0, got error: %v\nstderr: %s", err, stderr.String())
	}

	if stdout.String() != "" {
		t.Errorf("expected empty stdout, got: %s", stdout.String())
	}
}
