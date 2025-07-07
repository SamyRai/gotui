# GoTUI

> Interactive terminal dashboard for Go development: test, benchmark, format, vet, and build — all in one beautiful TUI.

[![Go Reference](https://pkg.go.dev/badge/github.com/damirmukimov/gotui.svg)](https://pkg.go.dev/github.com/damirmukimov/gotui)
[![CI](https://github.com/damirmukimov/gotui/actions/workflows/ci.yml/badge.svg)](https://github.com/damirmukimov/gotui/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/damirmukimov/gotui)](https://goreportcard.com/report/github.com/damirmukimov/gotui)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

<!-- ![GoTUI Demo](docs/demo.gif) -->

GoTUI is a single-binary CLI tool that provides a rich terminal UI for common Go development tasks. Run it in any Go module directory to get live streaming results from tests, benchmarks, formatting, vetting, and builds.

---

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage](#usage)
- [Key Bindings](#key-bindings)
- [Development](#development)
- [Contributing](#contributing)
- [License](#license)

---

## Features

**Five integrated tabs with live streaming results:**

- 🧪 **Tests** (`T`): Stream `go test -json` results with real-time pass/fail counters
- ⚙️ **Benchmarks** (`B`): Sortable performance table with comparison to previous runs
- ✨ **Format** (`F`): Preview `gofmt` diffs and auto-format files
- 🔎 **Vet** (`V`): Live `go vet` and `staticcheck` results with file navigation
- 🔨 **Build** (`C`): Real-time compilation with error reporting

**Key capabilities:**

- **Live streaming output** from all Go commands
- **Interactive navigation** with vim-like keybindings
- **Rich UI components** powered by Charmbracelet Bubbletea
- **File integration** - open files directly in your `$EDITOR`
- **Error handling** with clear notifications and graceful recovery
- **Zero configuration** - works in any Go module directory

---

## Installation

### Go Install (Recommended)

```bash
go install github.com/damirmukimov/gotui/cmd/goutui@latest
```

### Build from Source

```bash
git clone https://github.com/damirmukimov/gotui.git
cd gotui
go build -o goutui ./cmd/goutui
```

### Download Binary

Download the latest release from [GitHub Releases](https://github.com/damirmukimov/gotui/releases).

---

## Quick Start

```bash
# Navigate to any Go project
cd your-go-project

# Launch GoTUI
goutui
```

That's it! GoTUI will automatically detect your Go module and present the dashboard.

---

## Usage

### Basic Navigation

- **Switch tabs**: `Tab` / `Shift+Tab` or direct jump with `T`, `B`, `F`, `V`, `C`
- **Navigate lists**: `j` / `k` or arrow keys
- **Expand/collapse**: `Enter`
- **Refresh current tab**: `r`
- **Quit**: `q` or `Ctrl+C`

### Tab-Specific Features

#### Tests Tab (T)

- Toggle fail-only filter: `f`
- Open failing test in editor: `o`
- Live search: `/`

#### Benchmarks Tab (B)

- Sort by different columns
- Compare with previous runs
- View detailed performance metrics

#### Format Tab (F)

- Auto-format selected file: `a`
- Auto-format all files: `A`
- Preview diffs before applying

#### Vet Tab (V)

- Navigate to problematic lines: `o`
- View detailed lint messages
- Integrate with `staticcheck` if available

#### Build Tab (C)

- View compilation errors
- Navigate to error locations
- Real-time build status

---

## Key Bindings

| Key               | Context     | Action                                      |
| ----------------- | ----------- | ------------------------------------------- |
| `Tab / Shift+Tab` | Global      | Next / prev tab                             |
| `T/B/F/V/C`       | Global      | Jump to Tests / Bench / Fmt / Vet / Compile |
| `j / k`           | Lists       | Down / up                                   |
| `Enter`           | Tests / Vet | Expand/collapse logs                        |
| `f`               | Tests       | Toggle *fail-only* filter                   |
| `/`               | Lists       | Live search                                 |
| `r`               | Any tab     | Re-run current command                      |
| `o`               | Vet / Tests | Open file@line in `$EDITOR`                |
| `a / A`           | FmtDiff     | Auto-fmt selected / all                     |
| `q / Ctrl+C`      | Global      | Quit                                        |

---

## Development

### Prerequisites

- Go 1.24 or later
- A terminal with 256-color support

### Setup

```bash
git clone https://github.com/damirmukimov/gotui.git
cd gotui
go mod download
```

### Running

```bash
# Run from source
go run ./cmd/goutui

# Run tests
go test ./...

# Run linter
go vet ./...
```

### Project Structure

```text
gotui/
├── cmd/goutui/           # Main application entry point
├── internal/
│   ├── tui/              # Bubbletea models and components
│   │   ├── model.go      # Root model (tab router)
│   │   ├── tabs/         # Individual tab implementations
│   │   └── components/   # Reusable UI components
│   ├── runner/           # Command execution and parsing
│   ├── style/            # Lipgloss themes and styling
│   └── util/             # Utilities (editor integration, etc.)
└── docs/                 # Documentation and examples
```

### Architecture

GoTUI follows clean architecture principles:

- **Separation of concerns**: UI, business logic, and command execution are clearly separated
- **Event-driven**: Uses Bubbletea's message-passing architecture
- **Composable components**: Reusable UI components with consistent interfaces
- **Graceful error handling**: Robust error layer with user-friendly notifications

---

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Quick Contribution Guide

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes and add tests
4. Ensure tests pass: `go test ./...`
5. Ensure linting passes: `go vet ./...`
6. Commit your changes: `git commit -m 'feat: add amazing feature'`
7. Push to the branch: `git push origin feature/amazing-feature`
8. Open a Pull Request

### Development Roadmap

See [ROADMAP.md](ROADMAP.md) for planned features and development priorities.

---

## License

MIT License - see [LICENSE](LICENSE) file for details.

---

## Acknowledgments

Built with love using:

- [Charmbracelet Bubbletea](https://github.com/charmbracelet/bubbletea) - The TUI framework
- [Charmbracelet Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [Charmbracelet Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling

---

*Made with ❤️ for the Go community*
