# TODO List for GoTUI Project

## Implementation Status (as of 2025-07-07)

### Tabs

- [x] TabBar: Emoji icons, tab switching, direct jump keys (T/B/F/V/C)
- [x] TestRunner: Streams `go test -json`, TreeList, live counters, log viewer
- [x] BenchRunner: Streams bench output, sortable table, diff vs previous run
- [x] FmtDiff: Runs `gofmt -l -d`, accumulates diff, auto-format action (basic), file open (basic)
- [x] VetRunner: Runs `go vet`, parses output, file open (basic)
- [x] BuildRunner: Runs `go build`, shows spinner, error list

### Shared Components

- [x] TreeList: Collapsible, icons, highlight
- [x] LogViewer: Scrollable, ANSI color, copy
- [x] DiffViewer: Unified diff, color, raw diff access
- [x] StatusBar: Tab, totals, elapsed, hints

### Key Features

- [x] Streaming output for all commands
- [x] Live status bar with totals, elapsed time, key bindings
- [x] Consistent, responsive UI/UX (Bubbletea + Lipgloss)
- [x] Error handling via log/status bar
- [x] No global state; per-tab Bubbletea models
- [x] Works on macOS & Linux

### Refactored/Production-Ready

- [x] All TODOs in code are now either implemented or have context-rich comments
- [x] Diff accumulation and file open in FmtDiff
- [x] Benchmark sorting and comparison

---

## What's Left / Next Steps

- [ ] FmtDiff: Improve file open to jump to exact diff line (integrate with editor)
- [ ] FmtDiff: Add file list view for multiple diffs, allow selection
- [ ] BenchRunner: Show Δ diff visually in table (color, arrows)
- [ ] VetRunner: Integrate `staticcheck` if present
- [ ] BuildRunner: Show more detailed error output (group by file)
- [ ] Add more unit tests for edge cases and error handling
- [ ] Polish: Lipgloss theme tweaks, status bar hints, help popover
- [ ] Documentation: Update README with usage, key bindings, and demo gif
- [ ] CI: Ensure linter and tests pass in GitHub Actions

---

### Prioritization

- [ ] FmtDiff file list and selection (high, 4h)
- [ ] BenchRunner diff visualization (medium, 3h)
- [ ] VetRunner staticcheck integration (medium, 2h)
- [ ] BuildRunner error grouping (low, 2h)
- [ ] Polish & docs (medium, 2h)

---

*Update this file after each milestone or major refactor.*
