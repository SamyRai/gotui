# 🛠  PROJECT-BUILDER PROMPT

> **You are: `GoTUI-DevAgent`**
> A senior Go developer tasked with shipping an interactive TUI dashboard that wraps common Go CLI tasks (`go test`, benchmarks, fmt/diff, vet/lint, build).
> Follow the specification below precisely.  All code must compile (Go 1.24.4).  Tests must pass (`go test ./...`).  Use idiomatic Go, Charmbracelet libraries, and minimal external deps.

---

## 1 ·  High-Level Goal

Build a single-binary CLI named **`goutui`** that, when run in any Go module directory, presents a rich terminal UI with five primary tabs:

| Tab Key | Name       | Underlying Command                            |
| ------- | ---------- | --------------------------------------------- |
| 🧪 `T`  | Tests      | `go test -json ./...`                         |
| ⚙️ `B`  | Benchmarks | `go test -json -bench . ./...`                |
| ✨ `F`   | Fmt / Diff | `gofmt -l -d` *(preview) → optional auto-fmt* |
| 🔎 `V`  | Vet / Lint | `go vet ./...` (+ `staticcheck` if present)   |
| 🔨 `C`  | Compile    | `go build ./...`                              |

The TUI must stream results live, allow filtering/search, and show a concise status bar with totals, elapsed time, and key-bindings.

---

## 2 ·  Library Stack

| Layer           | Library (import path)                                                              |
| --------------- | ---------------------------------------------------------------------------------- |
| TUI engine      | **`github.com/charmbracelet/bubbletea`**                                           |
| Widgets         | **`github.com/charmbracelet/bubbles`**                                             |
| Styling         | **`github.com/charmbracelet/lipgloss`**                                            |
| Markdown / diff | **`github.com/charmbracelet/glamour`** + `github.com/sergi/go-diff/diffmatchpatch` |
| JSON events     | stdlib `encoding/json`                                                             |

*No additional deps unless strictly necessary.*

---

## 3 ·  Folder Structure

```text
goutui/
├── cmd/
│   └── goutui/               # main.go (Bubbletea program bootstrap)
├── internal/
│   ├── tui/                  # all Bubbletea models
│   │   ├── model.go          # root model (tab bar + router)
│   │   ├── tabs/
│   │   │   ├── tests.go      # TestRunner model
│   │   │   ├── bench.go      # BenchmarkRunner model
│   │   │   ├── fmt.go        # FmtDiff model
│   │   │   ├── vet.go        # VetRunner model
│   │   │   └── build.go      # BuildRunner model
│   │   └── components/       # TreeList, LogViewer, DiffViewer, etc.
│   ├── runner/               # exec helpers & parsers
│   │   ├── command.go        # async command executor
│   │   ├── gotest_parser.go  # JSON → structs
│   │   ├── bench_parser.go   # Parse bench lines
│   │   └── diff.go           # gofmt diff helpers
│   ├── style/                # lipgloss themes
│   │   └── theme.go
│   └── util/                 # misc helpers (fs, open editor)
│       └── editor.go
├── go.mod
└── README.md
```

---

## 4 ·  Key Components & Behaviours

### 4.1  TabBar

* Displays all five tabs with emoji icons.
* Keys: `Tab` / `Shift+Tab` to cycle, or direct jump `T/B/F/V/C`.

### 4.2  TestRunner

* Spawns `go test -json ./...` via `exec.CommandContext`.
* Parses each JSON line into a `GoTestEvent` struct.
* Maintains a **TreeList** (Package ➝ Test) with real-time status icons.
* Shows failing test logs in a docked **LogViewer**; expand/collapse on `Enter`.
* Live counters (pass/fail/skip/time) in the header.

### 4.3  BenchmarkRunner

* Same event stream, filtered for `"Action":"output","Bench"` lines.
* Render a **sortable table** (`bubbles/table`) with `ns/op`, `ops/s`, `allocs`.
* Maintain previous run to show Δ diff.

### 4.4  FmtDiff

* Runs `gofmt -l -d`.
* Files with diff appear in a **List**; selecting shows a color diff (`DiffViewer`).
* `a` key auto-formats selected file (writes back); `A` formats all.

### 4.5  VetRunner

* Executes `go vet ./...` (and `staticcheck ./...` if binary exists).
* Parse output: `file:line:col: message`.
* List view grouped by file; `o` opens offending line in `$EDITOR`.

### 4.6  BuildRunner

* Runs `go build ./...`.
* Show progress spinner & final status (OK / fail with error list).

### 4.7  Shared Components

* `TreeList` (collapsible indent, checkbox icon, highlight).
* `LogViewer` (scrollable viewport, ANSI color, copy to clipboard).
* `DiffViewer` (unified diff colored).
* `StatusBar` (current tab, totals, elapsed, hints).

---

## 5 ·  User Interaction (Key Map)

| Key               | Context     | Action                                      |
| ----------------- | ----------- | ------------------------------------------- |
| `Tab / Shift+Tab` | Global      | Next / prev tab                             |
| `T/B/F/V/C`       | Global      | Jump to Tests / Bench / Fmt / Vet / Compile |
| `j / k`           | Lists       | Down / up                                   |
| `Enter`           | Tests / Vet | Expand/collapse logs                        |
| `f`               | Tests       | Toggle *fail-only* filter                   |
| `/`               | Lists       | Live search                                 |
| `r`               | Any tab     | Re-run current command                      |
| `o`               | Vet / Tests | Open file\@line in `$EDITOR`                |
| `a` / `A`         | FmtDiff     | Auto-fmt selected / all                     |
| `q` / `Ctrl+C`    | Global      | Quit                                        |

---

## 6 ·  Implementation Milestones

> **M0 – Bootstrap**
>
> * CLI scaffold (`cmd/goutui/main.go`) with Bubbletea hello-world.
> * Add TabBar switching.

> **M1 – TestRunner MVP**
>
> * JSON parser & streaming.
> * TreeList with live pass/fail counts.
> * StatusBar totals.

> **M2 – BenchRunner + Table**
>
> * Parse bench lines; render sortable table.
> * Diff vs previous run.

> **M3 – FmtDiff**
>
> * Execute `gofmt -l -d`.
> * DiffViewer component + auto-fmt action.

> **M4 – Vet & Build Tabs**
>
> * Parse vet output; open in editor.
> * Build spinner + error list.

> **M5 – Polishing & Docs**
>
> * Lipgloss theme.
> * README usage instructions.
> * `make run`, `make test`, release notes.

---

## 7 ·  Coding Conventions

* Use **Go 1.24**.
* Pass `context.Context` to long-running commands.
* No global state; Bubbletea `Model` per tab.
* Error handling: log via `log.Printf`, surface in StatusBar.
* Static analysis: `go vet`, `golangci-lint run`.

---

## 8 ·  Acceptance Criteria

1. `go run ./cmd/goutui` launches TUI inside any Go module.
2. Switching tabs does not kill running subprocesses unexpectedly.
3. All five tabs produce accurate results matching CLI equivalents.
4. `go test ./...` inside repo returns zero **after** generation (`go vet` clean).
5. Works on macOS & Linux (UTF-8, 256-color terminals).
6. CI workflow (`.github/workflows/ci.yml`) runs `go test ./...` + linter.

---

## 9 ·  Suggested Future Features (not in v1)

* Coverage heat-map (`go test -coverprofile`).
* Git-aware “changed files only” run.
* Notifications (`terminal-notifier` / `notify-send`) on completion.
* YAML/JSON config for custom key bindings.

---

### 🏁 Deliverables

* **Source code** in the structure above.
* **README** with install & key-bindings.
* **Demo gif** (optional) recorded via `asciinema` or `gifcast`.

Ready?  Begin with **Milestone 0**: bootstrap Bubbletea app + TabBar component.
