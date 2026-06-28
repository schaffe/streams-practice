package cli

import (
	"bytes"
	"strings"
	"testing"

	"streams-practice/assignments/autocomplete/corpus"
)

type mockSearcher struct {
	entries []corpus.Entry
}

func (m *mockSearcher) Search(prefix string) []corpus.Entry {
	return m.entries
}

func TestRunBatch_EmptyPrefixReturnsNoResults(t *testing.T) {
	var out bytes.Buffer
	s := &mockSearcher{}
	cli := NewCLI(s, nil, &out)

	err := cli.RunBatch("")
	if err != nil {
		t.Fatalf("RunBatch() error = %v", err)
	}
	if out.String() != "" {
		t.Errorf("RunBatch() output = %q, want empty", out.String())
	}
}

func TestRunBatch_WithResults(t *testing.T) {
	var out bytes.Buffer
	s := &mockSearcher{entries: []corpus.Entry{
		{Term: "apple", Frequency: 10},
		{Term: "application", Frequency: 5},
	}}
	cli := NewCLI(s, nil, &out)

	err := cli.RunBatch("app")
	if err != nil {
		t.Fatalf("RunBatch() error = %v", err)
	}

	want := "apple (10)\napplication (5)\n"
	if out.String() != want {
		t.Errorf("RunBatch() output = %q, want %q", out.String(), want)
	}
}

func TestRunBatch_NoResultsIsNotAnError(t *testing.T) {
	var out bytes.Buffer
	s := &mockSearcher{}
	cli := NewCLI(s, nil, &out)

	err := cli.RunBatch("zzz")
	if err != nil {
		t.Fatalf("RunBatch() error = %v, want nil", err)
	}
	if out.String() != "" {
		t.Errorf("RunBatch() output = %q, want empty", out.String())
	}
}

func TestDrive_QuitOnCtrlC(t *testing.T) {
	var out bytes.Buffer
	s := &mockSearcher{}
	cli := NewCLI(s, strings.NewReader("\x03"), &out)

	err := cli.drive()
	if err != nil {
		t.Fatalf("drive() error = %v", err)
	}
}

func TestDrive_QuitOnQ(t *testing.T) {
	var out bytes.Buffer
	s := &mockSearcher{}
	cli := NewCLI(s, strings.NewReader("q"), &out)

	err := cli.drive()
	if err != nil {
		t.Fatalf("drive() error = %v", err)
	}
}

func TestDrive_PrintableAddsToPrefixAndRenders(t *testing.T) {
	var out bytes.Buffer
	s := &mockSearcher{entries: []corpus.Entry{
		{Term: "apple", Frequency: 10},
	}}
	cli := NewCLI(s, strings.NewReader("aq"), &out)

	err := cli.drive()
	if err != nil {
		t.Fatalf("drive() error = %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "> a") {
		t.Errorf("output should contain prompt '> a', got: %q", output)
	}
	if !strings.Contains(output, "apple") {
		t.Errorf("output should contain result 'apple', got: %q", output)
	}
}

func TestDrive_BackspaceRemovesLastChar(t *testing.T) {
	var out bytes.Buffer
	s := &mockSearcher{
		entries: []corpus.Entry{
			{Term: "apple", Frequency: 10},
		},
	}
	// Type 'a', then backspace, then 'a' again, then quit
	cli := NewCLI(s, strings.NewReader("a\baq"), &out)

	err := cli.drive()
	if err != nil {
		t.Fatalf("drive() error = %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "> a") {
		t.Errorf("output should contain prompt after final 'a', got: %q", output)
	}
}

func TestDrive_BackspaceOnEmptyPrefixIsNoop(t *testing.T) {
	var out bytes.Buffer
	s := &mockSearcher{}
	cli := NewCLI(s, strings.NewReader("\bq"), &out)

	err := cli.drive()
	if err != nil {
		t.Fatalf("drive() error = %v", err)
	}
}

func TestDrive_TabCompletesToFirstResult(t *testing.T) {
	var out bytes.Buffer
	s := &mockSearcher{
		entries: []corpus.Entry{
			{Term: "apple", Frequency: 10},
			{Term: "application", Frequency: 5},
		},
	}
	// Type 'a', then tab → prefix becomes "apple"
	cli := NewCLI(s, strings.NewReader("a\tq"), &out)

	err := cli.drive()
	if err != nil {
		t.Fatalf("drive() error = %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "> apple") {
		t.Errorf("output should contain completed prefix '> apple', got: %q", output)
	}
}

func TestDrive_TabWithNoResultsIsNoop(t *testing.T) {
	var out bytes.Buffer
	s := &mockSearcher{}
	cli := NewCLI(s, strings.NewReader("\tq"), &out)

	err := cli.drive()
	if err != nil {
		t.Fatalf("drive() error = %v", err)
	}
}

func TestDrive_EnterIsNoop(t *testing.T) {
	var out bytes.Buffer
	s := &mockSearcher{}
	cli := NewCLI(s, strings.NewReader("\nq"), &out)

	err := cli.drive()
	if err != nil {
		t.Fatalf("drive() error = %v", err)
	}
}

func TestDrive_CarriageReturnIsNoop(t *testing.T) {
	var out bytes.Buffer
	s := &mockSearcher{}
	cli := NewCLI(s, strings.NewReader("\rq"), &out)

	err := cli.drive()
	if err != nil {
		t.Fatalf("drive() error = %v", err)
	}
}

func TestDrive_NonPrintableIgnored(t *testing.T) {
	var out bytes.Buffer
	s := &mockSearcher{}
	// 0x01 is Ctrl-A, should be ignored
	cli := NewCLI(s, strings.NewReader("\x01q"), &out)

	err := cli.drive()
	if err != nil {
		t.Fatalf("drive() error = %v", err)
	}
}

func TestDrive_RendersNumberedResults(t *testing.T) {
	var out bytes.Buffer
	s := &mockSearcher{
		entries: []corpus.Entry{
			{Term: "apple", Frequency: 10},
			{Term: "application", Frequency: 5},
		},
	}
	cli := NewCLI(s, strings.NewReader("aq"), &out)

	err := cli.drive()
	if err != nil {
		t.Fatalf("drive() error = %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "  1. apple (10)") {
		t.Errorf("output should contain numbered result, got: %q", output)
	}
	if !strings.Contains(output, "  2. application (5)") {
		t.Errorf("output should contain numbered result, got: %q", output)
	}
}

func TestDrive_EmptyResultsAfterPrefixShowsPromptOnly(t *testing.T) {
	var out bytes.Buffer
	s := &mockSearcher{}
	cli := NewCLI(s, strings.NewReader("xq"), &out)

	err := cli.drive()
	if err != nil {
		t.Fatalf("drive() error = %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "> x") {
		t.Errorf("output should contain prompt, got: %q", output)
	}
}
