# Contributing to GoTUI

Thank you for your interest in contributing to GoTUI! This document provides guidelines and information for contributors.

---

## Table of Contents

- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Contributing Guidelines](#contributing-guidelines)
- [Code Standards](#code-standards)
- [Testing](#testing)
- [Pull Request Process](#pull-request-process)
- [Issue Reporting](#issue-reporting)
- [Community](#community)

---

## Getting Started

### Prerequisites

- **Go 1.24+**: Ensure you have Go installed and updated
- **Git**: For version control
- **Terminal**: 256-color terminal support recommended
- **Editor**: Any editor with Go support (VS Code, Vim, etc.)

### Quick Setup

```bash
# Fork the repository on GitHub
# Clone your fork
git clone https://github.com/YOUR_USERNAME/gotui.git
cd gotui

# Add upstream remote
git remote add upstream https://github.com/damirmukimov/gotui.git

# Install dependencies
go mod download

# Verify everything works
go run ./cmd/goutui
```

---

## Development Setup

### Project Structure

```text
gotui/
├── cmd/goutui/           # Main application
├── internal/
│   ├── tui/              # UI components and models
│   ├── runner/           # Command execution
│   ├── style/            # Styling and themes
│   └── util/             # Utilities
├── docs/                 # Documentation
└── examples/             # Example configurations
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific test
go test ./internal/runner -v

# Run tests with race detection
go test -race ./...
```

### Linting and Formatting

```bash
# Format code
go fmt ./...

# Run vet
go vet ./...

# Run additional linters (if installed)
golangci-lint run
staticcheck ./...
```

---

## Contributing Guidelines

### Types of Contributions

We welcome various types of contributions:

- **Bug fixes**: Help us squash bugs
- **Feature development**: Implement new features from our roadmap
- **Documentation**: Improve docs, examples, and guides
- **Testing**: Add tests and improve coverage
- **UI/UX**: Enhance the user interface and experience
- **Performance**: Optimize code and improve efficiency

### Before You Start

1. **Check existing issues**: Look for related issues or discussions
2. **Create an issue**: For new features or significant changes
3. **Discuss first**: Reach out for guidance on large contributions
4. **Read the roadmap**: Check [ROADMAP.md](ROADMAP.md) for planned features

### Workflow

1. **Create a branch**: Use descriptive branch names

   ```bash
   git checkout -b feature/add-coverage-tab
   git checkout -b fix/test-parser-crash
   git checkout -b docs/update-installation
   ```

2. **Make changes**: Follow our code standards
3. **Test thoroughly**: Ensure all tests pass
4. **Commit with clear messages**: Use conventional commit format
5. **Push and create PR**: Submit for review

---

## Code Standards

### Go Style Guide

We follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) and these additional guidelines:

#### Naming Conventions

```go
// Good: Clear, descriptive names
type TestRunner struct {
    parser   *GoTestParser
    treeList components.TreeList
}

func (tr *TestRunner) UpdateTestResults(result TestResult) error {
    // Implementation
}

// Avoid: Abbreviated or unclear names
type TR struct {
    p  *GTP
    tl components.TL
}
```

#### Error Handling

```go
// Good: Explicit error handling
result, err := tr.parser.ParseLine(line)
if err != nil {
    return fmt.Errorf("failed to parse test line: %w", err)
}

// Avoid: Ignoring errors
result, _ := tr.parser.ParseLine(line)
```

#### Interface Design

```go
// Good: Small, focused interfaces
type TestParser interface {
    ParseLine(line string) (TestResult, error)
}

// Avoid: Large interfaces
type TestManager interface {
    ParseLine(line string) (TestResult, error)
    RunTests() error
    FormatOutput() string
    SaveResults() error
    // ... many more methods
}
```

### Bubbletea Patterns

#### Model Updates

```go
// Good: Clear message handling
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        return m.handleKeyPress(msg)
    case runner.CommandOutput:
        return m.handleCommandOutput(msg)
    default:
        return m, nil
    }
}

// Avoid: Huge switch statements in Update
```

#### Component Composition

```go
// Good: Composable components
type TestTab struct {
    treeList  components.TreeList
    logViewer components.LogViewer
    parser    *runner.TestParser
}

func (t *TestTab) View() string {
    return lipgloss.JoinVertical(
        lipgloss.Left,
        t.treeList.View(),
        t.logViewer.View(),
    )
}
```

### Documentation

- **Public functions**: Must have doc comments
- **Complex logic**: Add inline comments
- **TODOs**: Include context and priority
- **Examples**: Provide usage examples for public APIs

```go
// ParseTestResult parses a JSON line from go test -json output.
// It returns a TestResult struct or an error if parsing fails.
//
// Example:
//   result, err := parser.ParseTestResult(`{"Action":"pass","Test":"TestExample"}`)
//   if err != nil {
//       log.Fatal(err)
//   }
//   fmt.Printf("Test %s: %s\n", result.Name, result.Status)
func ParseTestResult(line string) (TestResult, error) {
    // Implementation
}
```

---

## Testing

### Test Coverage

We aim for **90%+ test coverage** across the codebase:

- **Unit tests**: Test individual functions and methods
- **Integration tests**: Test component interactions
- **UI tests**: Test Bubbletea model behavior

### Test Structure

```go
func TestTestParser_ParseLine(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected TestResult
        wantErr  bool
    }{
        {
            name:     "successful test parse",
            input:    `{"Action":"pass","Test":"TestExample","Package":"example"}`,
            expected: TestResult{Action: "pass", Test: "TestExample", Package: "example"},
            wantErr:  false,
        },
        {
            name:    "invalid JSON",
            input:   `{"Action":"pass"`,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            parser := NewTestParser()
            result, err := parser.ParseLine(tt.input)
            
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            
            assert.NoError(t, err)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### Benchmarks

Include benchmarks for performance-critical code:

```go
func BenchmarkTestParser_ParseLine(b *testing.B) {
    parser := NewTestParser()
    line := `{"Action":"pass","Test":"TestExample","Package":"example"}`
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := parser.ParseLine(line)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

---

## Pull Request Process

### Before Submitting

- [ ] All tests pass locally
- [ ] Code is formatted (`go fmt`)
- [ ] No linting errors (`go vet`)
- [ ] Documentation is updated
- [ ] Commit messages follow conventional format

### PR Description Template

```markdown
## Description
Brief description of the changes.

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests added/updated
- [ ] Integration tests added/updated
- [ ] Manual testing completed

## Screenshots (if applicable)
Add screenshots for UI changes.

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] No breaking changes (or clearly documented)
```

### Review Process

1. **Automated checks**: CI must pass
2. **Code review**: At least one maintainer review
3. **Testing**: Manual testing if needed
4. **Approval**: Maintainer approval required
5. **Merge**: Squash and merge preferred

---

## Issue Reporting

### Bug Reports

Use the bug report template and include:

- **Environment**: OS, terminal, Go version
- **Steps to reproduce**: Clear, minimal reproduction steps
- **Expected vs actual behavior**: What should happen vs what happens
- **Screenshots/logs**: If applicable
- **Version**: GoTUI version or commit hash

### Feature Requests

Use the feature request template and include:

- **Problem description**: What problem does this solve?
- **Proposed solution**: How should it work?
- **Alternatives considered**: Other approaches considered
- **Use cases**: Real-world usage scenarios

---

## Community

### Communication

- **GitHub Issues**: Bug reports and feature requests
- **GitHub Discussions**: General questions and ideas
- **Code Review**: Constructive feedback on PRs

### Code of Conduct

We follow the [Contributor Covenant](https://www.contributor-covenant.org/). Be respectful, inclusive, and constructive in all interactions.

### Recognition

Contributors are recognized in:

- **README acknowledgments**: Regular contributors
- **Release notes**: Feature contributions
- **GitHub contributors page**: All contributions tracked

---

## Questions?

- Check [existing issues](https://github.com/damirmukimov/gotui/issues)
- Start a [discussion](https://github.com/damirmukimov/gotui/discussions)
- Review the [roadmap](ROADMAP.md)

Thank you for contributing to GoTUI! 🚀
