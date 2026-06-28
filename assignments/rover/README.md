# Rover — Handoff PRD

A robot ("rover") that moves on an **infinite** 2D integer grid, driven by three
commands: `L` (turn left), `R` (turn right), `M` (move forward one cell). Starts
at `(0, 0)` facing **North**. Two entry points: a live interactive terminal loop
(default) and a scriptable batch mode (`-cmd`).

Status: **implemented, tested, review-clean.** `go test ./assignments/rover/...`
is green across all packages; `go vet` is clean. Nothing committed yet.

---

## 1. Problem & Goals

Provide a small, extensible simulation of a grid robot with a clean separation
between the moving unit, the world that hosts it, and the UI that drives it.

- Correct rover kinematics on an infinite grid (no bounds).
- Three packages — `entity`, `world`, `ui` — integrated **only through
  interfaces**. No package references another's concrete types; `main` does all
  wiring. Dependencies flow one way: `main → ui → world → entity`. No cycles.
- Both interactive and batch entry points share the same world + entity.
- Extensible seams: addressable commands (carry a unit id), a render-agnostic
  snapshot, and a `Renderer` interface so a future 2D terminal renderer drops in
  without touching `world` or `entity`.

### Non-goals (deliberately out of scope)

- No 2D terminal graphics — only the `Renderer` seam exists for it.
- No multi-unit *driving* in the UI — the world supports many units, but the CLI
  drives a single hardcoded active unit (`ui.ActiveUnit = "rover-1"`).
- No persistence, networking, obstacles, or grid bounds.
- No Windows-specific raw-mode work beyond what `golang.org/x/term` provides.

---

## 2. How to run

```bash
# Interactive (raw-mode terminal loop): press L / R / M to drive; q or Ctrl-C to quit.
go run ./assignments/rover

# Batch (scriptable): apply a command string, print final "(x, y) H", exit.
go run ./assignments/rover -cmd "MMRMM"     # → (2, 2) E   (exit 0)
go run ./assignments/rover -cmd "X"         # → error on stderr (exit 1)

# Tests
go test ./assignments/rover/...
```

Status line format: `(x, y) H`, where `H` is `N` / `E` / `S` / `W`.

---

## 3. Architecture

```
assignments/rover/
├── main.go              // flag parsing (-cmd), wiring
├── entity/
│   └── rover.go         // UnitID, Heading, Rotation, Rover (implements world.Unit)
├── world/
│   └── world.go         // Unit iface, Action, Command, Snapshot, UnitState, New (unexported grid)
└── ui/
    ├── ui.go            // World iface, Renderer iface, ActiveUnit, lineRenderer
    └── cli.go           // CLI: RunBatch + RunInteractive (raw mode) + drive() loop
```

Module: `streams-practice` (go 1.21). Import paths:
`streams-practice/assignments/rover/{entity,world,ui}`.

**Interfaces are consumer-defined**, so dependencies point one way and concrete
types never cross a package boundary:

- `entity.Rover` implements `world.Unit` (defined in `world`).
- `world`'s concrete type (`grid`) is **unexported**; callers get it from
  `world.New(units...)` and only ever use it through `ui.World`.
- `ui.CLI` drives behavior through the `ui.World` interface it defines — never
  the concrete `grid`. `ui` imports `world` only for the `Command`/`Snapshot`
  value types.
- `main` is the only place concrete types meet.

### Data flow (one `L` keypress)

```
ui:    read 'l' → world.Command{Unit: ActiveUnit, Action: world.TurnLeft}
ui:    World.Apply(cmd)
world: look up unit by cmd.Unit → unit.Turn(entity.Left)
entity:Rover heading rotates North → West
ui:    World.Snapshot() → Renderer.Render → "(0, 0) W"
```

---

## 4. Package contracts (current API)

### entity

```go
type UnitID string
type Heading int            // North, East, South, West (clockwise); String() → "N"/"E"/"S"/"W"
type Rotation int           // Left, Right

type Rover struct { /* id, x, y, heading */ }
func NewRover(id UnitID, x, y int, h Heading) *Rover
func (r *Rover) ID() UnitID
func (r *Rover) Move()                            // advance one cell along heading
func (r *Rover) Turn(d Rotation)                  // rotate left/right (mod-4)
func (r *Rover) Position() (x, y int, h Heading)
```

Kinematics (North = +Y, East = +X): North (0,+1), East (+1,0), South (0,-1),
West (-1,0). Right = +1 mod 4, Left = +3 mod 4.

### world

```go
type Unit interface { ID() entity.UnitID; Move(); Turn(entity.Rotation); Position() (int, int, entity.Heading) }
type Action int             // TurnLeft, TurnRight, Forward
type Command struct { Unit entity.UnitID; Action Action }   // addressed instruction
type UnitState struct { ID entity.UnitID; X, Y int; Heading entity.Heading }
type Snapshot struct { Units []UnitState }                  // render-agnostic

func New(units ...Unit) *grid          // grid is unexported; use via ui.World
func (g *grid) Apply(cmd Command) error   // routes by id; error on unknown unit
func (g *grid) Snapshot() Snapshot        // all units, insertion order
```

### ui

```go
type World interface { Apply(world.Command) error; Snapshot() world.Snapshot }   // consumer-defined
type Renderer interface { Render(w io.Writer, s world.Snapshot) }
const ActiveUnit entity.UnitID = "rover-1"
func NewLineRenderer() Renderer                          // writes "(x, y) H" for ActiveUnit

type CLI struct { /* world, renderer, in, out */ }
func NewCLI(w World, r Renderer, in io.Reader, out io.Writer) *CLI
func (c *CLI) RunBatch(commands string) error            // parse string, apply, render final line
func (c *CLI) RunInteractive() error                     // raw mode + drive()
```

Input mapping (case-insensitive): `L`→TurnLeft, `R`→TurnRight, `M`→Forward,
`q`/`Q` or Ctrl-C → quit (interactive). `RunInteractive` puts stdin in raw mode
via `golang.org/x/term` (`MakeRaw`/`Restore`, always restored on exit) then
delegates to the terminal-free `drive()` loop, which is what the tests exercise.

---

## 5. Behavior reference (verified by tests)

| Input (`-cmd`) | Output    |
|----------------|-----------|
| `""`           | `(0, 0) N` |
| `M`            | `(0, 1) N` |
| `RM`           | `(1, 0) E` |
| `MMLM`         | `(-1, 2) W` |
| `LLLL`         | `(0, 0) N` |
| `MMRMM`        | `(2, 2) E` |
| `X` (invalid)  | error → stderr, exit 1 |

### Error handling

- Unknown unit id in `Apply` → error. Batch: stderr + exit 1. Interactive:
  surfaced on the status line, loop continues.
- Unrecognized input char — Batch: error (strict, it's the test path).
  Interactive: ignored.
- Raw-mode setup failure (e.g. stdin not a TTY) → `RunInteractive` returns the
  error; `main` prints it and exits 1. Terminal state is always restored.

---

## 6. Tests

- `entity/rover_test.go` — move deltas per heading, turn cycles, `Heading.String()`.
- `world/world_test.go` — `Apply` routing by id, unknown-id error, snapshot
  state, multi-unit independence and insertion order.
- `ui/cli_test.go` — `RunBatch` over the canonical strings (full ui+world+entity
  stack), asserting exact `(x, y) H`; invalid-char error case.
- `ui/interactive_test.go` — drives the terminal-free `drive()` loop with a
  scripted `io.Reader`: command sequences, `q`/`Q` quit, Ctrl-C (0x03).
- `smoke_test.go` (`package main`) — runs the real binary via `go run`:
  `-cmd "MMRMM"` → `(2, 2) E`, exit 0; `-cmd "X"` → stderr contains `X`, non-zero exit.

---

## 7. Extensibility (designed, not built)

- **2D terminal renderer:** implement `ui.Renderer` over the same `Snapshot`;
  swap it in `main`. `world`/`entity` untouched.
- **More commands:** add an `Action` value + a case in `grid.Apply`, plus a key
  mapping in the UI. `entity` gains a method only if it's new mechanics.
- **More units / UI targeting:** `Command` already carries a unit id and
  `Snapshot` holds all units, so a future UI can switch the active unit or render
  every unit with no interface changes.

---

## 8. Known limitations / follow-ups

- Interactive raw-mode TTY behavior is validated manually, not in CI (no TTY in
  `go test`); the logic underneath is covered via `drive()`.
- `drive()` surfaces `io.EOF` if interactive stdin closes without a `q`/Ctrl-C
  terminator — outside the spec's exit paths, currently harmless. Tighten if a
  piped-stdin interactive mode is ever wanted.
- `q` inside a batch string is silently skipped (neither acted on nor rejected);
  decide and document if batch should treat it as invalid.

Full design rationale: `.rapid/specs/2026-06-27-rover.md`.
Plan / DAG: `.rapid/plans/rover/`.
