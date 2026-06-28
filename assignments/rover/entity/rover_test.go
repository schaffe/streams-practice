package entity

import "testing"

// TestNewRoverAndPosition verifies that NewRover stores and returns the provided values
func TestNewRoverAndPosition(t *testing.T) {
	tests := []struct {
		name            string
		id              UnitID
		x               int
		y               int
		heading         Heading
		expectedID      UnitID
		expectedX       int
		expectedY       int
		expectedHeading Heading
	}{
		{
			name:            "rover at origin facing north",
			id:              "rover1",
			x:               0,
			y:               0,
			heading:         North,
			expectedID:      "rover1",
			expectedX:       0,
			expectedY:       0,
			expectedHeading: North,
		},
		{
			name:            "rover at positive coordinates facing east",
			id:              "rover2",
			x:               5,
			y:               3,
			heading:         East,
			expectedID:      "rover2",
			expectedX:       5,
			expectedY:       3,
			expectedHeading: East,
		},
		{
			name:            "rover at negative coordinates facing south",
			id:              "rover3",
			x:               -2,
			y:               -4,
			heading:         South,
			expectedID:      "rover3",
			expectedX:       -2,
			expectedY:       -4,
			expectedHeading: South,
		},
		{
			name:            "rover facing west",
			id:              "rover4",
			x:               10,
			y:               7,
			heading:         West,
			expectedID:      "rover4",
			expectedX:       10,
			expectedY:       7,
			expectedHeading: West,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRover(tt.id, tt.x, tt.y, tt.heading)

			id := r.ID()
			if id != tt.expectedID {
				t.Errorf("ID() = %v, want %v", id, tt.expectedID)
			}

			x, y, h := r.Position()
			if x != tt.expectedX {
				t.Errorf("Position() x = %d, want %d", x, tt.expectedX)
			}
			if y != tt.expectedY {
				t.Errorf("Position() y = %d, want %d", y, tt.expectedY)
			}
			if h != tt.expectedHeading {
				t.Errorf("Position() heading = %v, want %v", h, tt.expectedHeading)
			}
		})
	}
}

// TestMoveNorth verifies that Move increments y by 1 when facing north
func TestMoveNorth(t *testing.T) {
	tests := []struct {
		name      string
		startX    int
		startY    int
		expectedX int
		expectedY int
	}{
		{"from origin", 0, 0, 0, 1},
		{"from positive", 5, 3, 5, 4},
		{"from negative y", 2, -1, 2, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRover("test", tt.startX, tt.startY, North)
			r.Move()
			x, y, _ := r.Position()
			if x != tt.expectedX || y != tt.expectedY {
				t.Errorf("After Move(), got (%d, %d), want (%d, %d)", x, y, tt.expectedX, tt.expectedY)
			}
		})
	}
}

// TestMoveEast verifies that Move increments x by 1 when facing east
func TestMoveEast(t *testing.T) {
	tests := []struct {
		name      string
		startX    int
		startY    int
		expectedX int
		expectedY int
	}{
		{"from origin", 0, 0, 1, 0},
		{"from positive", 3, 5, 4, 5},
		{"from negative x", -1, 2, 0, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRover("test", tt.startX, tt.startY, East)
			r.Move()
			x, y, _ := r.Position()
			if x != tt.expectedX || y != tt.expectedY {
				t.Errorf("After Move(), got (%d, %d), want (%d, %d)", x, y, tt.expectedX, tt.expectedY)
			}
		})
	}
}

// TestMoveSouth verifies that Move decrements y by 1 when facing south
func TestMoveSouth(t *testing.T) {
	tests := []struct {
		name      string
		startX    int
		startY    int
		expectedX int
		expectedY int
	}{
		{"from origin", 0, 0, 0, -1},
		{"from positive", 3, 5, 3, 4},
		{"from positive y=1", 2, 1, 2, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRover("test", tt.startX, tt.startY, South)
			r.Move()
			x, y, _ := r.Position()
			if x != tt.expectedX || y != tt.expectedY {
				t.Errorf("After Move(), got (%d, %d), want (%d, %d)", x, y, tt.expectedX, tt.expectedY)
			}
		})
	}
}

// TestMoveWest verifies that Move decrements x by 1 when facing west
func TestMoveWest(t *testing.T) {
	tests := []struct {
		name      string
		startX    int
		startY    int
		expectedX int
		expectedY int
	}{
		{"from origin", 0, 0, -1, 0},
		{"from positive", 3, 5, 2, 5},
		{"from positive x=1", 1, 2, 0, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRover("test", tt.startX, tt.startY, West)
			r.Move()
			x, y, _ := r.Position()
			if x != tt.expectedX || y != tt.expectedY {
				t.Errorf("After Move(), got (%d, %d), want (%d, %d)", x, y, tt.expectedX, tt.expectedY)
			}
		})
	}
}

// TestTurnLeftCycle verifies that Turn(Left) cycles: North→West→South→East→North
func TestTurnLeftCycle(t *testing.T) {
	r := NewRover("test", 0, 0, North)

	// First left turn: North → West
	r.Turn(Left)
	_, _, h := r.Position()
	if h != West {
		t.Errorf("After first Turn(Left), got %v, want West", h)
	}

	// Second left turn: West → South
	r.Turn(Left)
	_, _, h = r.Position()
	if h != South {
		t.Errorf("After second Turn(Left), got %v, want South", h)
	}

	// Third left turn: South → East
	r.Turn(Left)
	_, _, h = r.Position()
	if h != East {
		t.Errorf("After third Turn(Left), got %v, want East", h)
	}

	// Fourth left turn: East → North
	r.Turn(Left)
	_, _, h = r.Position()
	if h != North {
		t.Errorf("After fourth Turn(Left), got %v, want North", h)
	}
}

// TestTurnRightCycle verifies that Turn(Right) cycles: North→East→South→West→North
func TestTurnRightCycle(t *testing.T) {
	r := NewRover("test", 0, 0, North)

	// First right turn: North → East
	r.Turn(Right)
	_, _, h := r.Position()
	if h != East {
		t.Errorf("After first Turn(Right), got %v, want East", h)
	}

	// Second right turn: East → South
	r.Turn(Right)
	_, _, h = r.Position()
	if h != South {
		t.Errorf("After second Turn(Right), got %v, want South", h)
	}

	// Third right turn: South → West
	r.Turn(Right)
	_, _, h = r.Position()
	if h != West {
		t.Errorf("After third Turn(Right), got %v, want West", h)
	}

	// Fourth right turn: West → North
	r.Turn(Right)
	_, _, h = r.Position()
	if h != North {
		t.Errorf("After fourth Turn(Right), got %v, want North", h)
	}
}

// TestFourLeftTurnsReturnToOriginal verifies that four Left turns return to the original heading
func TestFourLeftTurnsReturnToOriginal(t *testing.T) {
	startHeadings := []Heading{North, East, South, West}

	for _, startHeading := range startHeadings {
		t.Run("from "+startHeading.String(), func(t *testing.T) {
			r := NewRover("test", 0, 0, startHeading)

			r.Turn(Left)
			r.Turn(Left)
			r.Turn(Left)
			r.Turn(Left)

			_, _, h := r.Position()
			if h != startHeading {
				t.Errorf("After four Left turns, got %v, want %v", h, startHeading)
			}
		})
	}
}

// TestFourRightTurnsReturnToOriginal verifies that four Right turns return to the original heading
func TestFourRightTurnsReturnToOriginal(t *testing.T) {
	startHeadings := []Heading{North, East, South, West}

	for _, startHeading := range startHeadings {
		t.Run("from "+startHeading.String(), func(t *testing.T) {
			r := NewRover("test", 0, 0, startHeading)

			r.Turn(Right)
			r.Turn(Right)
			r.Turn(Right)
			r.Turn(Right)

			_, _, h := r.Position()
			if h != startHeading {
				t.Errorf("After four Right turns, got %v, want %v", h, startHeading)
			}
		})
	}
}

// TestHeadingString verifies that Heading.String() returns the correct string representation
func TestHeadingString(t *testing.T) {
	tests := []struct {
		heading  Heading
		expected string
	}{
		{North, "N"},
		{East, "E"},
		{South, "S"},
		{West, "W"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.heading.String()
			if result != tt.expected {
				t.Errorf("Heading.String() = %q, want %q", result, tt.expected)
			}
		})
	}
}
