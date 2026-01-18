package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// Package represents a discovered Go package
type Package struct {
	Path         string   `json:"ImportPath"`
	Dir          string   `json:"Dir"`
	GoFiles      []string `json:"GoFiles"`
	TestFiles    []string `json:"TestFiles"`
	XTestGoFiles []string `json:"XTestGoFiles"`
}

// FileMatch represents a file that matches a pattern
type FileMatch struct {
	Path    string // Full file path
	Package string // Package path
	RelPath string // Relative path from root
}

// PatternMatcher defines a pattern to match against file contents or names
type PatternMatcher struct {
	Name           string         // Name of the pattern (e.g., "test", "benchmark")
	FilePattern    *regexp.Regexp // Pattern to match file names (optional)
	ContentPattern *regexp.Regexp // Pattern to match file contents (optional)
}

// MatchResult represents a match found in a file
type MatchResult struct {
	File    string   // File path
	Package string   // Package path
	Line    int      // Line number (if content pattern matched)
	Match   string   // The matched text
	Groups  []string // Captured groups from regex
}

// Discoverer discovers Go packages and files matching patterns
type Discoverer struct {
	ctx  context.Context
	root string
}

// NewDiscoverer creates a new discoverer for the given root directory
func NewDiscoverer(ctx context.Context, root string) *Discoverer {
	if root == "" {
		root = "."
	}
	absRoot, err := filepath.Abs(root)
	if err != nil {
		absRoot = root
	}
	return &Discoverer{
		ctx:  ctx,
		root: absRoot,
	}
}

// DiscoverPackages discovers all Go packages in the root directory and subdirectories
// Uses `go list` command which is the standard way Go discovers packages
func (d *Discoverer) DiscoverPackages() ([]Package, error) {
	cmd := exec.CommandContext(d.ctx, "go", "list", "-json", "./...")
	cmd.Dir = d.root

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run go list: %w", err)
	}

	var packages []Package
	decoder := json.NewDecoder(strings.NewReader(string(output)))

	for decoder.More() {
		var pkg Package
		if err := decoder.Decode(&pkg); err != nil {
			return nil, fmt.Errorf("failed to decode package: %w", err)
		}
		packages = append(packages, pkg)
	}

	return packages, nil
}

// DiscoverFilesByPattern discovers files matching the given pattern
// pattern can be a file name pattern (e.g., "*_test.go") or a content pattern
func (d *Discoverer) DiscoverFilesByPattern(matcher PatternMatcher) ([]FileMatch, error) {
	// If we have a file pattern, use file system scanning (more reliable for test files)
	// go list doesn't include TestFiles by default, so FS scanning is necessary
	if matcher.FilePattern != nil {
		matches, err := d.discoverFilesByPatternFS(matcher)
		// Always return FS scan results, even if empty (means no matches found)
		// Only return error if there was an actual error during scanning
		if err != nil {
			// If FS scan fails, try package-based as fallback
			// But for test files, FS scan should work
		} else {
			return matches, nil
		}
	}

	// Fallback to package-based discovery for content patterns or if FS scan failed
	packages, err := d.DiscoverPackages()
	if err != nil {
		return nil, err
	}

	var matches []FileMatch
	seen := make(map[string]bool)

	for _, pkg := range packages {
		// Check all Go files in the package
		allFiles := append(pkg.GoFiles, pkg.TestFiles...)
		allFiles = append(allFiles, pkg.XTestGoFiles...)

		for _, file := range allFiles {
			fullPath := filepath.Join(pkg.Dir, file)
			if seen[fullPath] {
				continue
			}

			// Check file name pattern if provided
			if matcher.FilePattern != nil {
				fileName := filepath.Base(file)
				if !matcher.FilePattern.MatchString(fileName) {
					continue
				}
			}

			relPath, err := filepath.Rel(d.root, fullPath)
			if err != nil {
				relPath = file
			}

			matches = append(matches, FileMatch{
				Path:    fullPath,
				Package: pkg.Path,
				RelPath: relPath,
			})
			seen[fullPath] = true
		}
	}

	return matches, nil
}

// discoverFilesByPatternFS discovers files by scanning the file system
// This is more reliable for finding test files since go list doesn't include them by default
func (d *Discoverer) discoverFilesByPatternFS(matcher PatternMatcher) ([]FileMatch, error) {
	if matcher.FilePattern == nil {
		return nil, fmt.Errorf("file pattern is required for file system discovery")
	}

	var matches []FileMatch
	seen := make(map[string]bool)

	// Walk the directory tree
	err := filepath.Walk(d.root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip files we can't access
		}

		// Skip hidden directories and vendor, but allow .git and other common dirs
		if info.IsDir() {
			base := filepath.Base(path)
			// Skip hidden directories except current directory
			if strings.HasPrefix(base, ".") && base != "." {
				// Allow some common hidden directories that might contain code
				if base == ".git" || base == ".github" {
					return filepath.SkipDir
				}
				// Skip other hidden directories
				return filepath.SkipDir
			}
			// Skip vendor and node_modules
			if base == "vendor" || base == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}

		// Check if file matches pattern
		fileName := filepath.Base(path)
		if !matcher.FilePattern.MatchString(fileName) {
			return nil
		}

		// Skip hidden files (files starting with .)
		if strings.HasPrefix(fileName, ".") {
			return nil
		}

		// Skip if already seen
		if seen[path] {
			return nil
		}

		// Get relative path
		relPath, err := filepath.Rel(d.root, path)
		if err != nil {
			relPath = fileName
		}

		// Try to determine package from directory
		// This is a best-effort approach
		pkgPath := d.determinePackageFromPath(path)

		matches = append(matches, FileMatch{
			Path:    path,
			Package: pkgPath,
			RelPath: relPath,
		})
		seen[path] = true

		return nil
	})

	return matches, err
}

// determinePackageFromPath tries to determine the package path from a file path
func (d *Discoverer) determinePackageFromPath(filePath string) string {
	// Get the directory containing the file
	dir := filepath.Dir(filePath)

	// Try to use go list to get the actual package path
	// Use relative path from root for go list
	relDir, err := filepath.Rel(d.root, dir)
	if err != nil {
		relDir = dir
	}

	// Convert to Go import path format (use forward slashes)
	goPath := filepath.ToSlash(relDir)
	if goPath == "." || goPath == "" {
		goPath = "./"
	} else if !strings.HasPrefix(goPath, "./") {
		goPath = "./" + goPath
	}

	cmd := exec.CommandContext(d.ctx, "go", "list", "-f", "{{.ImportPath}}", goPath)
	cmd.Dir = d.root

	output, err := cmd.Output()
	if err == nil {
		pkgPath := strings.TrimSpace(string(output))
		// Check if it's a valid package path (not an error message)
		if pkgPath != "" && !strings.HasPrefix(pkgPath, "_") && !strings.Contains(pkgPath, "cannot find package") {
			return pkgPath
		}
	}

	// Fallback: construct package path from relative directory
	if relDir == "." || relDir == "" {
		return "main"
	}

	// Convert to package path format
	pkgPath := filepath.ToSlash(relDir)
	return pkgPath
}

// FindMatchesInFiles searches for content patterns in the given files
// Returns matches found in each file
func (d *Discoverer) FindMatchesInFiles(files []FileMatch, matcher PatternMatcher) ([]MatchResult, error) {
	if matcher.ContentPattern == nil {
		return nil, fmt.Errorf("content pattern is required")
	}

	var results []MatchResult

	for _, file := range files {
		matches, err := d.findMatchesInFile(file.Path, file.Package, matcher)
		if err != nil {
			continue // Skip files we can't read
		}
		results = append(results, matches...)
	}

	return results, nil
}

// findMatchesInFile searches for pattern matches in a single file
func (d *Discoverer) findMatchesInFile(filePath, packagePath string, matcher PatternMatcher) ([]MatchResult, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var results []MatchResult
	lines := strings.Split(string(content), "\n")

	for lineNum, line := range lines {
		// Skip comments and empty lines for function patterns
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "/*") || trimmed == "" {
			continue
		}

		matches := matcher.ContentPattern.FindStringSubmatch(line)
		if matches != nil {
			result := MatchResult{
				File:    filePath,
				Package: packagePath,
				Line:    lineNum + 1,
				Match:   matches[0],
				Groups:  matches[1:], // Captured groups
			}
			results = append(results, result)
		}
	}

	return results, nil
}

// NewPatternMatcher creates a new pattern matcher
func NewPatternMatcher(name string, filePattern, contentPattern string) (*PatternMatcher, error) {
	matcher := &PatternMatcher{
		Name: name,
	}

	if filePattern != "" {
		pattern, err := regexp.Compile(filePattern)
		if err != nil {
			return nil, fmt.Errorf("invalid file pattern: %w", err)
		}
		matcher.FilePattern = pattern
	}

	if contentPattern != "" {
		pattern, err := regexp.Compile(contentPattern)
		if err != nil {
			return nil, fmt.Errorf("invalid content pattern: %w", err)
		}
		matcher.ContentPattern = pattern
	}

	return matcher, nil
}

// GetPackageInfo returns detailed information about a specific package
func (d *Discoverer) GetPackageInfo(pkgPath string) (*Package, error) {
	cmd := exec.CommandContext(d.ctx, "go", "list", "-json", pkgPath)
	cmd.Dir = d.root

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get package info: %w", err)
	}

	var pkg Package
	if err := json.Unmarshal(output, &pkg); err != nil {
		return nil, fmt.Errorf("failed to decode package: %w", err)
	}

	return &pkg, nil
}

// FilterFiles filters files by package path pattern
func (d *Discoverer) FilterFiles(files []FileMatch, packagePattern string) []FileMatch {
	if packagePattern == "" {
		return files
	}

	var filtered []FileMatch
	pattern := regexp.MustCompile(packagePattern)

	for _, file := range files {
		if pattern.MatchString(file.Package) {
			filtered = append(filtered, file)
		}
	}

	return filtered
}

// GroupByPackage groups files by package
func (d *Discoverer) GroupByPackage(files []FileMatch) map[string][]FileMatch {
	groups := make(map[string][]FileMatch)

	for _, file := range files {
		groups[file.Package] = append(groups[file.Package], file)
	}

	return groups
}

// GetRoot returns the root directory of the discoverer
func (d *Discoverer) GetRoot() string {
	return d.root
}
