# Autocomplete — Handoff PRD

A search autocomplete CLI that suggests the top 10 completions for a partial
query, ranked by search frequency. Built as a trie with precomputed top-k per
node for O(k) lookup. Two entry points: a live interactive terminal loop
(default) and a scriptable batch mode (`-prefix`).

Status: **implemented, tested, review-clean.** `go test ./assignments/autocomplete/...`
is green across all packages; `go vet` is clean. Nothing committed yet.

---

## 1. Problem & Goals

Provide a search-autocomplete tool with a clean separation between the data
layer (corpus), the indexing layer (trie), and the UI that drives it.

- Correct prefix matching with top-10 suggestions ranked by frequency.
- Three packages — `corpus`, `index`, `cli` — integrated **only through
  interfaces**. No package references another's concrete types; `main` does all
  wiring. Dependencies flow one way: `main → cli → index → corpus`. No cycles.
- Both interactive and batch entry points share the same corpus + index.
- Extensible seams: a `corpus.Loader` interface for alternate data sources, a
  `cli.Searcher` interface for alternate indexing strategies.

### Non-goals (deliberately out of scope)

- No HTTP server or REST API.
- No fuzzy matching, spell correction, or partial-infix search — only prefix
  autocomplete.
- No persistence of query history or personalization.
- No concurrent query handling.
- No Windows-specific raw-mode work beyond what `golang.org/x/term` provides.

---

## 2. How to run

```bash
# Interactive (raw-mode terminal loop): type characters; suggestions update live.
# Tab accepts first suggestion; q or Ctrl-C to quit.
go run ./assignments/autocomplete

# Batch (scriptable): print top 10 for a prefix, one per line, exit.
go run ./assignments/autocomplete -prefix "ap"   # → apple (100), application (85), ...
go run ./assignments/autocomplete -prefix ""      # → top 10 overall

# Tests
go test ./assignments/autocomplete/...
```

Output format (batch): one `term (frequency)` per line, highest frequency first.

---

## 3. Architecture

```
assignments/autocomplete/
├── main.go                   // flag parsing (-prefix), wiring
├── README.md
├── smoke_test.go             // binary smoke tests (package main)
├── corpus/
│   ├── frequencies.json      // 1000 entries sorted by frequency desc
│   ├── corpus.go             // Entry, Load, Loader interface
│   └── corpus_test.go
├── index/
│   ├── trie.go               // Trie with max-heap per node (top 10)
│   └── trie_test.go
└── cli/
    ├── cli.go                // Searcher interface, CLI (interactive + batch)
    └── cli_test.go
```

Module: `streams-practice` (go 1.25). Import paths:
`streams-practice/assignments/autocomplete/{corpus,index,cli}`.

**Interfaces are consumer-defined**, so dependencies point one way and concrete
types never cross a package boundary:

- `index.Trie` implements `cli.Searcher` (defined in `cli`).
- `corpus.Load` is the concrete file loader; a network loader would implement
  `corpus.Loader`.
- `main` is the only place concrete types meet.

### Data flow (one keystroke in interactive)

```
cli:    read 'p' → prefix = "ap"
cli:    searcher.Search("ap")
index:  walk trie a→p, return node's top-10 (sorted by frequency desc)
cli:    render "> ap" + numbered list
```

---

## 4. Package contracts (current API)

### corpus

```go
type Entry struct {
    Term      string `json:"term"`
    Frequency int    `json:"frequency"`
}

type Loader interface { Load() ([]Entry, error) }

func Load(path string) ([]Entry, error)
```

`Load` reads a JSON file at `path` containing `[]Entry` and returns entries in
insertion order. Errors on missing file or malformed JSON.

### index

```go
type Trie struct { /* root node */ }

func New(entries []corpus.Entry) *Trie   // builds trie, precomputes top-10 per node
func (t *Trie) Search(prefix string) []corpus.Entry  // returns top 10 for prefix
```

Each trie node holds a max-heap of up to 10 entries from its subtree. During
insertion, each node along the path receives `node.Add(entry)`, which pushes
into the heap and evicts the lowest-frequency entry if size > 10. `Search`
walks the prefix (case-insensitive) and returns the heap contents sorted by
frequency descending. Non-destructive — repeated calls return the same results.

### cli

```go
type Searcher interface { Search(prefix string) []corpus.Entry }

type CLI struct { /* searcher, in, out */ }

func NewCLI(s Searcher, in io.Reader, out io.Writer) *CLI
func (c *CLI) RunBatch(prefix string) error         // print top 10, one per line
func (c *CLI) RunInteractive() error                 // raw-mode keystroke loop
```

Input mapping (interactive):

| Key | Behavior |
|-----|----------|
| Printable | Append to prefix, search, render |
| Backspace (0x7f / 0x08) | Pop last rune, search, render |
| Tab (0x09) | Replace prefix with first suggestion's term |
| `q` / `Q` / Ctrl-C | Restore terminal, exit |
| Enter | No-op |
| Other | Ignored |

Output format (batch): `term (frequency)` per line, highest frequency first.
Interactive: `> prefix` followed by numbered results, e.g.:
```
> ap
  1. apple (100)
  2. application (85)
```

---

## 5. Behavior reference (verified by tests)

| Input (`-prefix`) | Output |
|-------------------|--------|
| `""` | 10 lines, highest-frequency terms |
| `"the"` | ≤10 lines starting with "the", includes `the (1000)` |
| `"zzz"` | (empty) |
| `"ap"` (from generated data) | ≤10 lines starting with "ap", freq-descending |

### Error handling

- Missing/malformed `frequencies.json` → prints error to stderr, exit 1.
- No completions for prefix → batch: empty stdout; interactive: prompt line only.
- Raw-mode setup failure → `RunInteractive` returns error, main exits 1.
- Terminal state is always restored on exit (defer in `RunInteractive`).

---

## 6. Tests

- **corpus/corpus_test.go** — `Load` with valid JSON, empty array, malformed
  JSON, missing file, non-array JSON.
- **index/trie_test.go** — `Search` for empty prefix (top 10), prefix match
  (≤10 results, correct ordering), no match (nil), case-insensitive, Unicode,
  single-char prefix, whole-word prefix, zero/three entries, tie-breaking.
- **cli/cli_test.go** — `RunBatch` (results, no results, empty prefix).
  `drive()` loop (printable chars, backspace, tab, Enter, Ctrl-C, 'q'/Q,
  non-printable, numbered rendering, empty results).
- **smoke_test.go** (`package main`) — runs the real binary via `go run`:
  `-prefix "the"` outputs `the (1000)` with ≤10 lines; `-prefix ""` outputs
  10 lines; `-prefix "zzz"` outputs nothing; all exit 0.

---

## 7. Extensibility (designed, not built)

- **Network source:** implement `corpus.Loader` over HTTP; swap in `main`.
- **Fuzzy search:** implement `cli.Searcher` with Levenshtein matching.
- **More results:** change the heap-size constant in `index` (currently 10).
- **Enter selects first result:** trivial to add in the interactive key handler.
