// Package world provides a grid-based world for managing rover units.
package world

import (
	"fmt"

	"streams-practice/assignments/rover/entity"
)

// Action represents a command action type
type Action int

const (
	// TurnLeft represents a left turn action.
	TurnLeft Action = iota
	// TurnRight represents a right turn action.
	TurnRight
	// Forward represents a forward movement action.
	Forward
)

// Command represents a command to be applied to a unit
type Command struct {
	Unit   entity.UnitID
	Action Action
}

// UnitState represents the state of a unit at a point in time
type UnitState struct {
	ID      entity.UnitID
	X       int
	Y       int
	Heading entity.Heading
}

// Snapshot represents the state of all units at a point in time
type Snapshot struct {
	Units []UnitState
}

// Unit is the interface for a unit in the world
type Unit interface {
	ID() entity.UnitID
	Move()
	Turn(d entity.Rotation)
	Position() (x, y int, h entity.Heading)
}

// grid represents the world and manages units in it
type grid struct {
	units     map[entity.UnitID]Unit
	unitOrder []entity.UnitID // preserve insertion order
}

// New creates a new world with the given units
func New(units ...Unit) *grid {
	g := &grid{
		units:     make(map[entity.UnitID]Unit),
		unitOrder: make([]entity.UnitID, 0, len(units)),
	}
	for _, u := range units {
		id := u.ID()
		g.units[id] = u
		g.unitOrder = append(g.unitOrder, id)
	}
	return g
}

// Apply applies a command to the world
func (g *grid) Apply(cmd Command) error {
	u, ok := g.units[cmd.Unit]
	if !ok {
		return fmt.Errorf("unknown unit %q", cmd.Unit)
	}

	switch cmd.Action {
	case TurnLeft:
		u.Turn(entity.Left)
	case TurnRight:
		u.Turn(entity.Right)
	case Forward:
		u.Move()
	default:
		return fmt.Errorf("unknown action %v", cmd.Action)
	}

	return nil
}

// Snapshot returns a snapshot of the current state of all units
func (g *grid) Snapshot() Snapshot {
	states := make([]UnitState, 0, len(g.unitOrder))
	for _, id := range g.unitOrder {
		u := g.units[id]
		x, y, h := u.Position()
		states = append(states, UnitState{
			ID:      id,
			X:       x,
			Y:       y,
			Heading: h,
		})
	}
	return Snapshot{Units: states}
}
