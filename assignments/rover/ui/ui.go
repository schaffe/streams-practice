package ui

import (
	"fmt"
	"io"

	"streams-practice/assignments/rover/entity"
	"streams-practice/assignments/rover/world"
)

// World is the interface the CLI drives. Implemented by *world.grid (via world.New).
type World interface {
	Apply(cmd world.Command) error
	Snapshot() world.Snapshot
}

// Renderer knows how to render a world snapshot to a writer.
type Renderer interface {
	Render(w io.Writer, s world.Snapshot)
}

// ActiveUnit is the unit ID the CLI controls.
const ActiveUnit entity.UnitID = "rover-1"

// lineRenderer renders the active unit state as "(x, y) H".
type lineRenderer struct{}

// NewLineRenderer returns a Renderer that writes "(x, y) H" for ActiveUnit.
func NewLineRenderer() Renderer {
	return lineRenderer{}
}

func (lineRenderer) Render(w io.Writer, s world.Snapshot) {
	for _, u := range s.Units {
		if u.ID == ActiveUnit {
			fmt.Fprintf(w, "(%d, %d) %s", u.X, u.Y, u.Heading.String())
			return
		}
	}
	// ActiveUnit absent: render nothing (should not happen in normal wiring).
}
