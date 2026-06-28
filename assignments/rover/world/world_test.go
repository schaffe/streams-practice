package world

import (
	"testing"

	"streams-practice/assignments/rover/entity"
)

func TestNewCreatesGridWithUnits(t *testing.T) {
	rover := entity.NewRover("rover", 0, 0, entity.North)
	g := New(rover)

	snapshot := g.Snapshot()
	if len(snapshot.Units) != 1 {
		t.Fatalf("expected 1 unit, got %d", len(snapshot.Units))
	}
	if snapshot.Units[0].ID != "rover" {
		t.Errorf("expected ID 'rover', got %q", snapshot.Units[0].ID)
	}
	if snapshot.Units[0].X != 0 || snapshot.Units[0].Y != 0 {
		t.Errorf("expected position (0,0), got (%d,%d)", snapshot.Units[0].X, snapshot.Units[0].Y)
	}
	if snapshot.Units[0].Heading != entity.North {
		t.Errorf("expected heading North, got %v", snapshot.Units[0].Heading)
	}
}

// Test Forward action moves unit north
func TestApplyForwardMovesUnitNorth(t *testing.T) {
	rover := entity.NewRover("rover", 0, 0, entity.North)
	g := New(rover)

	err := g.Apply(Command{Unit: "rover", Action: Forward})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	snapshot := g.Snapshot()
	if snapshot.Units[0].Y != 1 {
		t.Errorf("expected Y=1 after moving forward north, got Y=%d", snapshot.Units[0].Y)
	}
	if snapshot.Units[0].X != 0 {
		t.Errorf("expected X=0, got X=%d", snapshot.Units[0].X)
	}
	if snapshot.Units[0].Heading != entity.North {
		t.Errorf("expected heading North, got %v", snapshot.Units[0].Heading)
	}
}

// Test Forward action moves unit east
func TestApplyForwardMovesUnitEast(t *testing.T) {
	rover := entity.NewRover("rover", 0, 0, entity.East)
	g := New(rover)

	err := g.Apply(Command{Unit: "rover", Action: Forward})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	snapshot := g.Snapshot()
	if snapshot.Units[0].X != 1 {
		t.Errorf("expected X=1 after moving forward east, got X=%d", snapshot.Units[0].X)
	}
	if snapshot.Units[0].Y != 0 {
		t.Errorf("expected Y=0, got Y=%d", snapshot.Units[0].Y)
	}
	if snapshot.Units[0].Heading != entity.East {
		t.Errorf("expected heading East, got %v", snapshot.Units[0].Heading)
	}
}

// Test Forward action moves unit south
func TestApplyForwardMovesUnitSouth(t *testing.T) {
	rover := entity.NewRover("rover", 0, 0, entity.South)
	g := New(rover)

	err := g.Apply(Command{Unit: "rover", Action: Forward})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	snapshot := g.Snapshot()
	if snapshot.Units[0].Y != -1 {
		t.Errorf("expected Y=-1 after moving forward south, got Y=%d", snapshot.Units[0].Y)
	}
	if snapshot.Units[0].X != 0 {
		t.Errorf("expected X=0, got X=%d", snapshot.Units[0].X)
	}
	if snapshot.Units[0].Heading != entity.South {
		t.Errorf("expected heading South, got %v", snapshot.Units[0].Heading)
	}
}

// Test Forward action moves unit west
func TestApplyForwardMovesUnitWest(t *testing.T) {
	rover := entity.NewRover("rover", 0, 0, entity.West)
	g := New(rover)

	err := g.Apply(Command{Unit: "rover", Action: Forward})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	snapshot := g.Snapshot()
	if snapshot.Units[0].X != -1 {
		t.Errorf("expected X=-1 after moving forward west, got X=%d", snapshot.Units[0].X)
	}
	if snapshot.Units[0].Y != 0 {
		t.Errorf("expected Y=0, got Y=%d", snapshot.Units[0].Y)
	}
	if snapshot.Units[0].Heading != entity.West {
		t.Errorf("expected heading West, got %v", snapshot.Units[0].Heading)
	}
}

// Test TurnLeft action rotates unit
func TestApplyTurnLeftRotatesUnitCounterclockwise(t *testing.T) {
	rover := entity.NewRover("rover", 0, 0, entity.North)
	g := New(rover)

	err := g.Apply(Command{Unit: "rover", Action: TurnLeft})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	snapshot := g.Snapshot()
	if snapshot.Units[0].Heading != entity.West {
		t.Errorf("expected heading West after TurnLeft from North, got %v", snapshot.Units[0].Heading)
	}
	// Position should not change
	if snapshot.Units[0].X != 0 || snapshot.Units[0].Y != 0 {
		t.Errorf("position should not change on rotation, got (%d,%d)", snapshot.Units[0].X, snapshot.Units[0].Y)
	}
}

// Test TurnRight action rotates unit
func TestApplyTurnRightRotatesUnitClockwise(t *testing.T) {
	rover := entity.NewRover("rover", 0, 0, entity.North)
	g := New(rover)

	err := g.Apply(Command{Unit: "rover", Action: TurnRight})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	snapshot := g.Snapshot()
	if snapshot.Units[0].Heading != entity.East {
		t.Errorf("expected heading East after TurnRight from North, got %v", snapshot.Units[0].Heading)
	}
	// Position should not change
	if snapshot.Units[0].X != 0 || snapshot.Units[0].Y != 0 {
		t.Errorf("position should not change on rotation, got (%d,%d)", snapshot.Units[0].X, snapshot.Units[0].Y)
	}
}

// Test Apply with unknown unit ID returns error
func TestApplyUnknownUnitIDReturnsError(t *testing.T) {
	rover := entity.NewRover("rover", 0, 0, entity.North)
	g := New(rover)

	err := g.Apply(Command{Unit: "unknown", Action: Forward})
	if err == nil {
		t.Fatal("expected error for unknown unit ID, got nil")
	}
}

// Test Apply with unknown unit ID does not mutate grid
func TestApplyUnknownUnitIDDoesNotMutateGrid(t *testing.T) {
	rover := entity.NewRover("rover", 0, 0, entity.North)
	g := New(rover)

	snapshotBefore := g.Snapshot()
	_ = g.Apply(Command{Unit: "unknown", Action: Forward})
	snapshotAfter := g.Snapshot()

	if snapshotBefore.Units[0].X != snapshotAfter.Units[0].X ||
		snapshotBefore.Units[0].Y != snapshotAfter.Units[0].Y ||
		snapshotBefore.Units[0].Heading != snapshotAfter.Units[0].Heading {
		t.Error("grid was mutated after applying command with unknown unit ID")
	}
}

// Test multiple units: applying to one doesn't affect another
func TestApplyToOneUnitDoesNotAffectAnother(t *testing.T) {
	rover1 := entity.NewRover("rover1", 0, 0, entity.North)
	rover2 := entity.NewRover("rover2", 5, 5, entity.East)
	g := New(rover1, rover2)

	err := g.Apply(Command{Unit: "rover1", Action: Forward})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	snapshot := g.Snapshot()
	// rover1 should move
	if snapshot.Units[0].Y != 1 {
		t.Errorf("rover1 should move to Y=1, got Y=%d", snapshot.Units[0].Y)
	}
	// rover2 should remain unchanged
	if snapshot.Units[1].X != 5 || snapshot.Units[1].Y != 5 || snapshot.Units[1].Heading != entity.East {
		t.Errorf("rover2 should not change, got (%d,%d,%v)", snapshot.Units[1].X, snapshot.Units[1].Y, snapshot.Units[1].Heading)
	}
}

// Test multiple units in snapshot are in insertion order
func TestSnapshotUnitsInInsertionOrder(t *testing.T) {
	rover1 := entity.NewRover("rover1", 0, 0, entity.North)
	rover2 := entity.NewRover("rover2", 5, 5, entity.East)
	rover3 := entity.NewRover("rover3", 3, 3, entity.South)
	g := New(rover1, rover2, rover3)

	snapshot := g.Snapshot()
	if len(snapshot.Units) != 3 {
		t.Fatalf("expected 3 units, got %d", len(snapshot.Units))
	}
	if snapshot.Units[0].ID != "rover1" {
		t.Errorf("first unit should be rover1, got %s", snapshot.Units[0].ID)
	}
	if snapshot.Units[1].ID != "rover2" {
		t.Errorf("second unit should be rover2, got %s", snapshot.Units[1].ID)
	}
	if snapshot.Units[2].ID != "rover3" {
		t.Errorf("third unit should be rover3, got %s", snapshot.Units[2].ID)
	}
}

// Test sequence of commands on same unit
func TestMultipleCommandsOnSameUnit(t *testing.T) {
	rover := entity.NewRover("rover", 0, 0, entity.North)
	g := New(rover)

	// Forward twice
	if err := g.Apply(Command{Unit: "rover", Action: Forward}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := g.Apply(Command{Unit: "rover", Action: Forward}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	snapshot := g.Snapshot()
	if snapshot.Units[0].Y != 2 {
		t.Errorf("expected Y=2 after two forwards, got Y=%d", snapshot.Units[0].Y)
	}

	// Turn right
	if err := g.Apply(Command{Unit: "rover", Action: TurnRight}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	snapshot = g.Snapshot()
	if snapshot.Units[0].Heading != entity.East {
		t.Errorf("expected heading East after TurnRight, got %v", snapshot.Units[0].Heading)
	}

	// Forward once more
	if err := g.Apply(Command{Unit: "rover", Action: Forward}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	snapshot = g.Snapshot()
	if snapshot.Units[0].X != 1 || snapshot.Units[0].Y != 2 {
		t.Errorf("expected position (1,2), got (%d,%d)", snapshot.Units[0].X, snapshot.Units[0].Y)
	}
}

// Test full rotation returns to original heading
func TestFullRotationReturnsToOriginalHeading(t *testing.T) {
	rover := entity.NewRover("rover", 0, 0, entity.North)
	g := New(rover)

	// Turn left 4 times (full rotation)
	for i := 0; i < 4; i++ {
		if err := g.Apply(Command{Unit: "rover", Action: TurnLeft}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	snapshot := g.Snapshot()
	if snapshot.Units[0].Heading != entity.North {
		t.Errorf("expected heading North after full rotation, got %v", snapshot.Units[0].Heading)
	}
}
