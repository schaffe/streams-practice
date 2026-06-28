package ui

import (
	"bytes"
	"strings"
	"testing"

	"streams-practice/assignments/rover/entity"
	"streams-practice/assignments/rover/world"
)

func TestRunBatch(t *testing.T) {
	tests := []struct {
		name     string
		commands string
		want     string
		wantErr  bool
	}{
		{
			name:     "empty commands",
			commands: "",
			want:     "(0, 0) N",
			wantErr:  false,
		},
		{
			name:     "move forward once",
			commands: "M",
			want:     "(0, 1) N",
			wantErr:  false,
		},
		{
			name:     "rotate right then move",
			commands: "RM",
			want:     "(1, 0) E",
			wantErr:  false,
		},
		{
			name:     "complex sequence",
			commands: "MMLM",
			want:     "(-1, 2) W",
			wantErr:  false,
		},
		{
			name:     "full rotation",
			commands: "LLLL",
			want:     "(0, 0) N",
			wantErr:  false,
		},
		{
			name:     "case insensitive lowercase",
			commands: "m",
			want:     "(0, 1) N",
			wantErr:  false,
		},
		{
			name:     "invalid command",
			commands: "X",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create fresh world and rover for each test
			rover := entity.NewRover(ActiveUnit, 0, 0, entity.North)
			w := world.New(rover)

			// Create CLI with line renderer and buffer
			buf := &bytes.Buffer{}
			r := NewLineRenderer()
			cli := NewCLI(w, r, nil, buf)

			// Run batch commands
			err := cli.RunBatch(tt.commands)

			// Check error expectation
			if (err != nil) != tt.wantErr {
				t.Fatalf("RunBatch() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return // Skip output check on error case
			}

			// Get output and trim single trailing newline
			output := buf.String()
			output = strings.TrimSuffix(output, "\n")

			if output != tt.want {
				t.Errorf("RunBatch() output = %q, want %q", output, tt.want)
			}
		})
	}
}
