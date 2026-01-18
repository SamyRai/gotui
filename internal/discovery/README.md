# Discovery Package

A generic, pattern-based file and package discovery system for Go projects.

## Overview

The discovery package provides a generic way to discover Go packages and files matching custom patterns. It doesn't know about specific file types (tests, benchmarks, etc.) - instead, each package defines its own patterns.

## Key Principles

1. **Generic**: No knowledge of specific file types or patterns
2. **Pattern-Based**: Uses regex patterns for file names and content
3. **Reusable**: Any package can define its own discovery patterns
4. **Go Standard**: Uses `go list` for package discovery (the standard Go way)

## Usage Examples

### Discovering Test Files

```go
import "goutui/internal/discovery"

ctx := context.Background()
d := discovery.NewDiscoverer(ctx, ".")

// Define test file pattern: files ending with _test.go
matcher, _ := discovery.NewPatternMatcher(
    "test",
    `.*_test\.go$`,  // File name pattern
    "",              // No content pattern needed
)

testFiles, err := d.DiscoverFilesByPattern(*matcher)
```

### Discovering Benchmark Functions

```go
// First find test files
fileMatcher, _ := discovery.NewPatternMatcher("test", `.*_test\.go$`, "")
testFiles, _ := d.DiscoverFilesByPattern(*fileMatcher)

// Then find Benchmark functions in those files
contentMatcher, _ := discovery.NewPatternMatcher(
    "benchmark",
    "",                          // No file pattern
    `^func\s+(Benchmark\w+)\s*\(`, // Content pattern
)

benchmarks, _ := d.FindMatchesInFiles(testFiles, *contentMatcher)
```

### Discovering Custom Patterns

```go
// Find all .go files in internal/ directory
matcher, _ := discovery.NewPatternMatcher(
    "internal",
    `^internal/.*\.go$`,
    "",
)

files, _ := d.DiscoverFilesByPattern(*matcher)
```

## API Reference

### Types

- `Package`: Represents a Go package discovered via `go list`
- `FileMatch`: A file that matches a pattern
- `PatternMatcher`: Defines patterns for file names and/or content
- `MatchResult`: A match found in file content

### Methods

- `NewDiscoverer(ctx, root)`: Create a discoverer for a root directory
- `DiscoverPackages()`: Discover all Go packages using `go list`
- `DiscoverFilesByPattern(matcher)`: Find files matching a pattern
- `FindMatchesInFiles(files, matcher)`: Search for content patterns in files
- `GroupByPackage(files)`: Group files by package
- `FilterFiles(files, packagePattern)`: Filter files by package pattern

## Best Practices

1. **Use `go list` for packages**: The discovery package uses `go list` which is the standard Go way
2. **Define patterns per use case**: Each package (tests, benchmarks, etc.) defines its own patterns
3. **Combine file and content patterns**: Use file patterns for discovery, content patterns for detailed analysis
4. **Reuse discoverer instances**: Create one discoverer per root directory and reuse it
