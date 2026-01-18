# GoTUI - Terminal UI for Go Development

[![Go Version](https://img.shields.io/badge/Go-1.25+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/github/actions/workflow/status/damirmukimov/gotui/ci.yml)](https://github.com/damirmukimov/gotui/actions)

A modern terminal user interface for Go development that provides an integrated development environment in your terminal. Run tests, benchmarks, formatting, linting, and building with live output and intuitive navigation.

## ✨ Features

### 🚀 Core Functionality
- **Real-time Testing**: Live `go test -json` output with hierarchical test tree view
- **Performance Analysis**: Sortable benchmark results with performance comparisons
- **Code Formatting**: Interactive `gofmt` diff preview with selective application
- **Static Analysis**: `go vet` and `staticcheck` integration with file navigation
- **Build Monitoring**: Live compilation output with error grouping

### 🎨 User Experience
- **Tabbed Interface**: Clean tabbed navigation between different development tasks
- **Keyboard-Driven**: Full keyboard navigation with mnemonic shortcuts
- **Editor Integration**: Seamless integration with your preferred editor (`$EDITOR`)
- **Live Updates**: Real-time command output with ANSI color support
- **Responsive Design**: Adapts to terminal size changes

### 🛠️ Developer-Friendly
- **Hot Reloading**: Development mode with automatic rebuilds
- **Comprehensive Testing**: 90%+ test coverage with integration tests
- **Structured Logging**: Configurable logging with development/production modes
- **Error Handling**: Robust error handling with detailed context
- **Clean Architecture**: Modular design following Go best practices

## 📦 Installation

### Quick Install (Recommended)

```bash
# Clone the repository
git clone https://github.com/damirmukimov/gotui.git
cd gotui

# Build and install
make install
```

### Manual Build

```bash
# Prerequisites: Go 1.25+
go version

# Clone and build
git clone https://github.com/damirmukimov/gotui.git
cd gotui
go build -o goutui ./cmd/goutui

# Optional: Install globally
sudo mv goutui /usr/local/bin/
```

### Development Setup

```bash
# Install development dependencies
make setup-dev

# Run in development mode with hot reloading
make dev

# Or run tests continuously
make test
```

## 🚀 Usage

### Basic Usage

```bash
# Run in current Go project directory
goutui

# Run with debug logging
GOTUI_DEBUG=1 goutui

# Run with custom editor
EDITOR=code goutui
```

### Key Bindings

#### Global Navigation
- `Tab` / `Shift+Tab` - Switch between tabs
- `t` / `b` / `f` / `v` / `c` - Jump to Tests/Benchmarks/Format/Vet/Build tabs
- `q` / `Ctrl+C` - Quit application
- `r` - Refresh current tab
- `Esc` - Return to tab bar from content

#### Content Navigation
- `↑` / `↓` / `j` / `k` - Navigate lists
- `Enter` - Open selected item in editor
- `Space` - Expand/collapse tree items
- `PgUp` / `PgDown` - Page through content

### Tab Reference

#### 🧪 Tests Tab
- **Command**: `go test -json ./...`
- **Features**: Hierarchical test tree, pass/fail indicators, file navigation
- **Navigation**: Expand/collapse test suites, jump to failing tests

#### 📊 Benchmarks Tab
- **Command**: `go test -bench=. ./...`
- **Features**: Sortable performance table, historical comparisons
- **Navigation**: Sort by time, memory, or allocations

#### 🎨 Format Tab
- **Command**: `gofmt -d .`
- **Features**: Diff preview, selective formatting, editor integration
- **Actions**: Apply formatting to individual files or all files

#### 🔍 Vet Tab
- **Command**: `go vet ./...` + `staticcheck`
- **Features**: Issue categorization, file and line navigation
- **Integration**: Direct jump to problematic code

#### 🔨 Build Tab
- **Command**: `go build ./...`
- **Features**: Real-time compilation, error grouping
- **Navigation**: Jump to compilation errors

## 🏗️ Architecture

```
gotui/
├── cmd/goutui/           # Application entry point
├── internal/
│   ├── editor/           # Editor integration utilities
│   ├── errors/           # Error handling and types
│   ├── logger/           # Structured logging system
│   ├── runner/           # Command execution engine
│   ├── style/            # UI theming and styling
│   └── tui/
│       ├── components/   # Reusable UI components
│       ├── tabs/         # Tab implementations
│       └── model.go      # Main application model
├── .air.toml             # Hot reload configuration
├── .golangci.yml         # Linting configuration
└── Makefile              # Build and development tasks
```

### Key Components

- **Bubble Tea Framework**: Modern TUI framework with event-driven architecture
- **Structured Logging**: slog-based logging with configurable output formats
- **Error Handling**: Comprehensive error types with context preservation
- **Command Runner**: Async command execution with output streaming
- **Component Architecture**: Modular UI components with clean interfaces

## 🔧 Configuration

### Environment Variables

```bash
# Enable debug logging and additional output
export GOTUI_DEBUG=1

# Set preferred editor (defaults to $EDITOR or $VISUAL)
export EDITOR=code
export VISUAL=nvim

# Enable development mode logging
export DEBUG=1
```

### Development Configuration

The `.air.toml` file configures hot reloading for development:

```toml
[build]
cmd = "go build -o ./tmp/goutui ./cmd/goutui"
bin = "./tmp/goutui"

[screen]
clear_on_rebuild = true
```

## 🧪 Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-cover

# Run tests in verbose mode
make test-verbose

# Run tests with race detection
make test-race
```

### Test Coverage

Current coverage by package:
- `internal/logger`: 100%
- `internal/errors`: 86.4%
- `internal/editor`: 58.2%

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Quick Start for Contributors

```bash
# Fork and clone
git clone https://github.com/your-username/gotui.git
cd gotui

# Set up development environment
make setup-dev

# Run tests and linting
make ci

# Start development with hot reload
make dev
```

### Development Workflow

1. **Setup**: `make setup-dev`
2. **Develop**: `make dev` (with hot reloading)
3. **Test**: `make test` (continuous testing)
4. **Lint**: `make lint` (code quality)
5. **Build**: `make build` (production build)

## 📚 Documentation

- [Contributing Guide](CONTRIBUTING.md)
- [Development Roadmap](ROADMAP.md)
- [API Documentation](docs/) *(coming soon)*

## 🔄 Roadmap

See [ROADMAP.md](ROADMAP.md) for upcoming features and development priorities.

### Recent Milestones ✅
- [x] Core TUI architecture with Bubble Tea
- [x] All major tabs (Tests, Benchmarks, Format, Vet, Build)
- [x] Editor integration and file navigation
- [x] Comprehensive error handling and logging
- [x] Development tooling and hot reloading

### Upcoming Features 🚧
- [ ] Configuration file support
- [ ] Theme customization
- [ ] Coverage integration
- [ ] Git integration
- [ ] Plugin system

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

## 🙏 Acknowledgments

Built with ❤️ using:
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - The elegant TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - UI components
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Style definitions
- [Air](https://github.com/cosmtrek/air) - Hot reloading for development

---

**Happy Go Development!** 🎉