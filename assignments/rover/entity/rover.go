// Package entity defines rover domain types including heading, rotation, and position.
package entity

// UnitID is a unique identifier for a rover unit.
type UnitID string

// Heading represents the direction a rover is facing.
type Heading int

const (
	// North represents the north direction.
	North Heading = iota
	// East represents the east direction.
	East
	// South represents the south direction.
	South
	// West represents the west direction.
	West
)

// String returns the single-character representation of the heading.
func (h Heading) String() string {
	switch h {
	case North:
		return "N"
	case East:
		return "E"
	case South:
		return "S"
	case West:
		return "W"
	default:
		return ""
	}
}

// Rotation represents a turn direction.
type Rotation int

const (
	// Left represents a counterclockwise turn.
	Left Rotation = iota
	// Right represents a clockwise turn.
	Right
)

// Rover represents a rover unit with position and heading.
type Rover struct {
	id      UnitID
	x       int
	y       int
	heading Heading
}

// NewRover creates a new rover with the given ID, position, and heading.
func NewRover(id UnitID, x, y int, h Heading) *Rover {
	return &Rover{
		id:      id,
		x:       x,
		y:       y,
		heading: h,
	}
}

// ID returns the rover's unit ID.
func (r *Rover) ID() UnitID {
	return r.id
}

// Move advances the rover one cell in the direction it is facing.
func (r *Rover) Move() {
	switch r.heading {
	case North:
		r.y++
	case East:
		r.x++
	case South:
		r.y--
	case West:
		r.x--
	}
}

// Turn rotates the rover left or right.
func (r *Rover) Turn(d Rotation) {
	if d == Right {
		r.heading = (r.heading + 1) % 4
	} else {
		r.heading = (r.heading + 3) % 4
	}
}

// Position returns the rover's current x, y coordinates and heading.
func (r *Rover) Position() (x, y int, h Heading) {
	return r.x, r.y, r.heading
}
