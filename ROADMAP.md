# GoTUI Development Roadmap

This document outlines the current status and future plans for GoTUI development.

---

## Current Status (v1.0.0-alpha)

### ✅ Completed Features

#### Core Architecture

- [x] Bubbletea-based TUI framework
- [x] Clean separation of concerns (UI, runner, components)
- [x] Event-driven architecture with robust error handling
- [x] Context-based cancellation for all commands
- [x] Graceful shutdown with proper cleanup

#### Tab System

- [x] **TabBar**: Emoji icons, keyboard navigation, direct tab jumping
- [x] **Tests Tab**: Real-time `go test -json` streaming, TreeList navigation, pass/fail counters
- [x] **Benchmarks Tab**: Performance table, sortable results, previous run comparison
- [x] **Format Tab**: `gofmt` diff preview, auto-format capabilities
- [x] **Vet Tab**: `go vet` results parsing, file navigation
- [x] **Build Tab**: Real-time compilation, error reporting

#### UI Components

- [x] **TreeList**: Collapsible hierarchy, status icons, keyboard navigation
- [x] **LogViewer**: Scrollable output, ANSI color support
- [x] **DiffViewer**: Unified diff rendering with syntax highlighting
- [x] **StatusBar**: Live counters, elapsed time, key hints

#### Developer Experience

- [x] Comprehensive key bindings
- [x] File integration (`$EDITOR` support)
- [x] Live streaming output from all commands
- [x] Error notifications with graceful recovery

---

## Upcoming Releases

### v1.0.0-beta (Target: Q3 2025)

#### High Priority Features

- [ ] **Enhanced Format Tab**
  - [ ] File list view for multiple diffs
  - [ ] Per-file selection and formatting
  - [ ] Precise line navigation in editor
  - **Effort**: 4-6 hours

- [ ] **Improved Benchmark Visualization**
  - [ ] Visual diff indicators (↑↓ arrows, colors)
  - [ ] Historical trend tracking
  - [ ] Performance regression alerts
  - **Effort**: 3-4 hours

- [ ] **Extended Vet Integration**
  - [ ] `staticcheck` support (auto-detect binary)
  - [ ] `golangci-lint` integration
  - [ ] Configurable linter selection
  - **Effort**: 2-3 hours

#### Quality Improvements

- [ ] **Enhanced Error Handling**
  - [ ] Better error grouping in Build tab
  - [ ] File-based error navigation
  - [ ] Error context and suggestions
  - **Effort**: 2-3 hours

- [ ] **Testing & Reliability**
  - [ ] Comprehensive unit test coverage (>90%)
  - [ ] Integration tests for command execution
  - [ ] Edge case handling improvements
  - **Effort**: 4-5 hours

- [ ] **UI Polish**
  - [ ] Refined Lipgloss themes
  - [ ] Help overlay/modal
  - [ ] Improved status bar hints
  - **Effort**: 2-3 hours

### v1.0.0 (Target: Q4 2025)

#### Documentation & Distribution

- [ ] **Complete Documentation**
  - [ ] Detailed usage guide
  - [ ] Configuration options
  - [ ] Troubleshooting guide
  - [ ] Demo GIF/video
  - **Effort**: 3-4 hours

- [ ] **Release Infrastructure**
  - [ ] GitHub Actions CI/CD
  - [ ] Multi-platform binary releases
  - [ ] Homebrew formula
  - [ ] Package manager distribution
  - **Effort**: 2-3 hours

- [ ] **Community Features**
  - [ ] Contributing guidelines
  - [ ] Issue templates
  - [ ] Code of conduct
  - [ ] Security policy
  - **Effort**: 1-2 hours

---

## Future Versions (v1.1+)

### Advanced Features

- [ ] **Coverage Integration**
  - [ ] `go test -coverprofile` support
  - [ ] Coverage heat-map visualization
  - [ ] Coverage trend tracking

- [ ] **Git Integration**
  - [ ] "Changed files only" mode
  - [ ] Pre-commit hook integration
  - [ ] Branch-aware testing

- [ ] **Configuration System**
  - [ ] YAML/JSON config files
  - [ ] Custom key bindings
  - [ ] Theme customization
  - [ ] Command customization

- [ ] **Notifications & Integrations**
  - [ ] Desktop notifications on completion
  - [ ] Slack/Discord webhook support
  - [ ] VS Code extension integration
  - [ ] Terminal multiplexer integration

### Performance & Scalability

- [ ] **Large Project Support**
  - [ ] Efficient handling of large test suites
  - [ ] Incremental result loading
  - [ ] Memory optimization

- [ ] **Advanced Filtering**
  - [ ] Regex-based filtering
  - [ ] Saved filter presets
  - [ ] Tag-based test selection

---

## Current Development Priorities

### Sprint 1 (Next 2 weeks)

1. **Format Tab Enhancement** - File list and selection
2. **Benchmark Visualization** - Visual diff indicators
3. **Vet Tab** - `staticcheck` integration

### Sprint 2 (Following 2 weeks)

1. **Error Handling** - Better grouping and navigation
2. **Testing** - Comprehensive test coverage
3. **Documentation** - Usage guide and demo

### Sprint 3 (Month 2)

1. **UI Polish** - Themes and help system
2. **Release Infrastructure** - CI/CD and distribution
3. **Community** - Contributing guidelines

---

## How to Contribute

We welcome contributions to any of these roadmap items! Here's how you can help:

1. **Check Current Issues**: Look for open issues labeled with roadmap items
2. **Propose Features**: Open an issue to discuss new features before implementing
3. **Submit PRs**: Follow our contribution guidelines when submitting code
4. **Test & Report**: Use GoTUI in your projects and report bugs or UX issues

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines.

---

## Changelog

### Recent Updates

- **2025-07-07**: Initial roadmap created
- **2025-07-07**: v1.0.0-alpha architecture completed
- **2025-07-07**: Core tab functionality implemented

---

*This roadmap is a living document and will be updated as development progresses.*
