package discovery

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewDiscoverer(t *testing.T) {
	ctx := context.Background()
	d := NewDiscoverer(ctx, ".")
	if d == nil {
		t.Error("Expected non-nil discoverer")
	}
	if d.root == "" {
		t.Error("Expected non-empty root")
	}

	// Test with empty root - should default to current directory
	d2 := NewDiscoverer(ctx, "")
	if d2.root == "" {
		t.Error("Expected non-empty root for empty string")
	}
}

func TestDiscoverPackages(t *testing.T) {
	ctx := context.Background()
	d := NewDiscoverer(ctx, "../../")

	packages, err := d.DiscoverPackages()
	if err != nil {
		t.Fatalf("Failed to discover packages: %v", err)
	}

	if len(packages) == 0 {
		t.Error("Expected at least one package")
	}

	// Check that we found the main package
	foundMain := false
	for _, pkg := range packages {
		if pkg.Path == "goutui" || pkg.Path == "." {
			foundMain = true
		}
		// Verify package structure
		if pkg.Dir == "" {
			t.Error("Package Dir should not be empty")
		}
	}

	if !foundMain {
		t.Log("Note: Main package not found, but this is OK if running from subdirectory")
	}
}

func TestDiscoverFilesByPattern(t *testing.T) {
	ctx := context.Background()
	// Use project root (go up from internal/discovery)
	d := NewDiscoverer(ctx, "../../")

	// Test pattern for test files
	matcher, err := NewPatternMatcher("test", `.*_test\.go$`, "")
	if err != nil {
		t.Fatalf("Failed to create pattern matcher: %v", err)
	}

	testFiles, err := d.DiscoverFilesByPattern(*matcher)
	if err != nil {
		t.Fatalf("Failed to discover test files: %v", err)
	}

	if len(testFiles) == 0 {
		t.Log("No test files found - this might be OK if running from wrong directory")
	}

	// Verify file match structure
	for _, file := range testFiles {
		if file.Path == "" {
			t.Error("FileMatch Path should not be empty")
		}
		if file.Package == "" {
			t.Error("FileMatch Package should not be empty")
		}
		if !filepath.IsAbs(file.Path) && !filepath.IsLocal(file.Path) {
			t.Errorf("FileMatch Path should be a valid path: %s", file.Path)
		}
	}
}

func TestFindMatchesInFiles(t *testing.T) {
	ctx := context.Background()
	d := NewDiscoverer(ctx, "../../")

	// First discover test files
	fileMatcher, err := NewPatternMatcher("test", `.*_test\.go$`, "")
	if err != nil {
		t.Fatalf("Failed to create file pattern matcher: %v", err)
	}

	testFiles, err := d.DiscoverFilesByPattern(*fileMatcher)
	if err != nil {
		t.Fatalf("Failed to discover test files: %v", err)
	}

	if len(testFiles) == 0 {
		t.Skip("No test files found, skipping content pattern test")
	}

	// Now search for Test functions in those files
	contentMatcher, err := NewPatternMatcher("test", "", `^func\s+(Test\w+)\s*\(`)
	if err != nil {
		t.Fatalf("Failed to create content pattern matcher: %v", err)
	}

	matches, err := d.FindMatchesInFiles(testFiles, *contentMatcher)
	if err != nil {
		t.Fatalf("Failed to find matches: %v", err)
	}

	// Verify match structure
	for _, match := range matches {
		if match.File == "" {
			t.Error("MatchResult File should not be empty")
		}
		if match.Package == "" {
			t.Error("MatchResult Package should not be empty")
		}
		if match.Line <= 0 {
			t.Errorf("MatchResult Line should be > 0, got %d", match.Line)
		}
		if match.Match == "" {
			t.Error("MatchResult Match should not be empty")
		}
	}
}

func TestGetPackageInfo(t *testing.T) {
	ctx := context.Background()
	d := NewDiscoverer(ctx, "../../")

	// Try to get info about the current package
	pkg, err := d.GetPackageInfo(".")
	if err != nil {
		t.Logf("Could not get package info for '.': %v (this is OK if not in a Go module)", err)
		return
	}

	if pkg == nil {
		t.Error("Expected non-nil package")
		return
	}

	if pkg.Dir == "" {
		t.Error("Package Dir should not be empty")
	}
}

func TestGroupByPackage(t *testing.T) {
	ctx := context.Background()
	d := NewDiscoverer(ctx, "../../")

	matcher, err := NewPatternMatcher("test", `.*_test\.go$`, "")
	if err != nil {
		t.Fatalf("Failed to create pattern matcher: %v", err)
	}

	testFiles, err := d.DiscoverFilesByPattern(*matcher)
	if err != nil {
		t.Fatalf("Failed to discover test files: %v", err)
	}

	groups := d.GroupByPackage(testFiles)

	if len(groups) == 0 && len(testFiles) > 0 {
		t.Error("Expected at least one group")
	}

	// Verify grouping
	for pkg, files := range groups {
		if pkg == "" {
			t.Error("Package name should not be empty")
		}
		if len(files) == 0 {
			t.Error("Package group should have at least one file")
		}
		for _, file := range files {
			if file.Package != pkg {
				t.Errorf("File package %q does not match group package %q", file.Package, pkg)
			}
		}
	}
}

func TestFilterFiles(t *testing.T) {
	ctx := context.Background()
	d := NewDiscoverer(ctx, "../../")

	matcher, err := NewPatternMatcher("test", `.*_test\.go$`, "")
	if err != nil {
		t.Fatalf("Failed to create pattern matcher: %v", err)
	}

	testFiles, err := d.DiscoverFilesByPattern(*matcher)
	if err != nil {
		t.Fatalf("Failed to discover test files: %v", err)
	}

	// Filter by a pattern that should match some packages
	filtered := d.FilterFiles(testFiles, ".*")

	if len(filtered) != len(testFiles) {
		t.Errorf("Filter with '.*' should return all files, got %d expected %d", len(filtered), len(testFiles))
	}

	// Filter with empty pattern should return all
	filtered2 := d.FilterFiles(testFiles, "")
	if len(filtered2) != len(testFiles) {
		t.Errorf("Filter with empty pattern should return all files, got %d expected %d", len(filtered2), len(testFiles))
	}
}

// TestDiscoverPackagesErrorHandling tests error handling in package discovery
func TestDiscoverPackagesErrorHandling(t *testing.T) {
	ctx := context.Background()

	// Test with non-existent directory
	d := NewDiscoverer(ctx, "/non/existent/directory")

	_, err := d.DiscoverPackages()
	if err == nil {
		t.Error("Expected error when discovering packages in non-existent directory")
	}

	// Test with invalid go list arguments (simulate by using a context that gets cancelled)
	ctx, cancel := context.WithCancel(ctx)
	cancel() // Cancel immediately

	d2 := NewDiscoverer(ctx, "../../")
	_, err = d2.DiscoverPackages()
	if err == nil {
		t.Error("Expected error when context is cancelled")
	}
}

// TestNewPatternMatcher tests pattern matcher creation
func TestNewPatternMatcher(t *testing.T) {
	tests := []struct {
		name            string
		matcherName     string
		filePattern     string
		contentPattern  string
		expectError     bool
	}{
		{
			name:           "valid file pattern only",
			matcherName:    "test",
			filePattern:    `.*\.go$`,
			contentPattern: "",
			expectError:    false,
		},
		{
			name:            "valid content pattern only",
			matcherName:     "test",
			filePattern:     "",
			contentPattern:  `func Test`,
			expectError:     false,
		},
		{
			name:            "both patterns",
			matcherName:     "test",
			filePattern:     `.*\.go$`,
			contentPattern:  `func Test`,
			expectError:     false,
		},
		{
			name:           "invalid file pattern",
			matcherName:    "test",
			filePattern:    `[invalid`,
			contentPattern: "",
			expectError:    true,
		},
		{
			name:            "invalid content pattern",
			matcherName:     "test",
			filePattern:     "",
			contentPattern:  `[invalid`,
			expectError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher, err := NewPatternMatcher(tt.matcherName, tt.filePattern, tt.contentPattern)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				if matcher != nil {
					t.Error("Expected nil matcher on error")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if matcher == nil {
					t.Error("Expected non-nil matcher")
				}
				if matcher.Name != tt.matcherName {
					t.Errorf("Expected name %q, got %q", tt.matcherName, matcher.Name)
				}
				if tt.filePattern != "" && matcher.FilePattern == nil {
					t.Error("Expected file pattern to be set")
				}
				if tt.contentPattern != "" && matcher.ContentPattern == nil {
					t.Error("Expected content pattern to be set")
				}
			}
		})
	}
}

// TestDiscoverFilesByPatternErrorHandling tests error handling in file discovery
func TestDiscoverFilesByPatternErrorHandling(t *testing.T) {
	ctx := context.Background()
	d := NewDiscoverer(ctx, "../../")

	// Test with invalid file pattern
	_, err := NewPatternMatcher("test", `[invalid`, "")
	if err == nil {
		t.Error("Expected error creating invalid pattern matcher")
	}

	// Test with matcher that has neither file nor content pattern
	emptyMatcher := &PatternMatcher{Name: "empty"}
	_, err = d.DiscoverFilesByPattern(*emptyMatcher)
	if err != nil {
		// Expected error for empty matcher - this is OK
		t.Logf("Empty matcher error (expected): %v", err)
	}

	// Test file system discovery with no file pattern
	_, err = d.discoverFilesByPatternFS(*emptyMatcher)
	if err == nil {
		t.Error("Expected error when file system discovery has no file pattern")
	}
}

// TestDeterminePackageFromPath tests package path determination
func TestDeterminePackageFromPath(t *testing.T) {
	ctx := context.Background()
	d := NewDiscoverer(ctx, "../../")

	// Test with a known file path
	testFile := filepath.Join(d.root, "internal", "discovery", "discovery.go")
	pkgPath := d.determinePackageFromPath(testFile)

	if pkgPath == "" {
		t.Error("Expected non-empty package path")
	}

	if !strings.Contains(pkgPath, "discovery") {
		t.Errorf("Expected package path to contain 'discovery', got %q", pkgPath)
	}
}

// TestFindMatchesInFilesErrorHandling tests error handling in content matching
func TestFindMatchesInFilesErrorHandling(t *testing.T) {
	ctx := context.Background()
	d := NewDiscoverer(ctx, "../../")

	// Create matcher with invalid content pattern
	_, err := NewPatternMatcher("test", "", `[invalid`)
	if err == nil {
		t.Error("Expected error creating invalid content pattern")
	}

	// Test with empty files list
	matcher, _ := NewPatternMatcher("test", "", `func Test`)
	emptyFiles := []FileMatch{}
	results, err := d.FindMatchesInFiles(emptyFiles, *matcher)
	if err != nil {
		t.Errorf("Unexpected error with empty files list: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("Expected empty results for empty files list, got %d", len(results))
	}
}

// TestWithTempDirectory tests discovery functionality with temporary directories
func TestWithTempDirectory(t *testing.T) {
	// Create temporary directory structure
	tempDir, err := os.MkdirTemp("", "discovery_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create some test Go files
	createTestGoFile(t, tempDir, "main.go", `package main

import "fmt"

func main() {
	fmt.Println("Hello")
}

func TestFunction() {
	// This should be found by content pattern
}`)

	createTestGoFile(t, tempDir, "utils.go", `package main

func Helper() {
	// Helper function
}`)

	createTestGoFile(t, tempDir, "utils_test.go", `package main

import "testing"

func TestHelper(t *testing.T) {
	Helper()
}`)

	// Create subdirectory with Go files
	subDir := filepath.Join(tempDir, "pkg")
	err = os.MkdirAll(subDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}

	createTestGoFile(t, subDir, "lib.go", `package pkg

func LibraryFunc() {
	// Library function
}`)

	ctx := context.Background()
	d := NewDiscoverer(ctx, tempDir)

	// Test package discovery (might fail in temp dir without go.mod)
	packages, err := d.DiscoverPackages()
	if err != nil {
		// This is expected in temp directories without proper Go modules
		t.Logf("Package discovery failed as expected in temp dir: %v", err)
	} else if len(packages) == 0 {
		t.Error("Expected to find at least one package")
	}

	// Test file discovery with file pattern
	fileMatcher, err := NewPatternMatcher("go_files", `.*\.go$`, "")
	if err != nil {
		t.Fatalf("Failed to create file matcher: %v", err)
	}

	goFiles, err := d.DiscoverFilesByPattern(*fileMatcher)
	if err != nil {
		t.Fatalf("Failed to discover Go files: %v", err)
	}

	expectedFiles := 4 // main.go, utils.go, utils_test.go, lib.go
	if len(goFiles) != expectedFiles {
		t.Errorf("Expected %d Go files, got %d", expectedFiles, len(goFiles))
	}

	// Test file discovery with test pattern
	testMatcher, err := NewPatternMatcher("test_files", `.*_test\.go$`, "")
	if err != nil {
		t.Fatalf("Failed to create test file matcher: %v", err)
	}

	testFiles, err := d.DiscoverFilesByPattern(*testMatcher)
	if err != nil {
		t.Fatalf("Failed to discover test files: %v", err)
	}

	if len(testFiles) != 1 {
		t.Errorf("Expected 1 test file, got %d", len(testFiles))
	}

	// Test content pattern matching
	contentMatcher, err := NewPatternMatcher("test_funcs", "", `^func Test\w+\(`)
	if err != nil {
		t.Fatalf("Failed to create content matcher: %v", err)
	}

	matches, err := d.FindMatchesInFiles(testFiles, *contentMatcher)
	if err != nil {
		t.Fatalf("Failed to find matches: %v", err)
	}

	if len(matches) != 1 {
		t.Errorf("Expected 1 test function match, got %d", len(matches))
	}

	if len(matches) > 0 {
		match := matches[0]
		if match.Match == "" {
			t.Error("Expected non-empty match text")
		}
		if !strings.Contains(match.Match, "TestHelper") {
			t.Errorf("Expected match to contain 'TestHelper', got %q", match.Match)
		}
	}
}

// TestGroupByPackageWithTempDir tests grouping functionality
func TestGroupByPackageWithTempDir(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "discovery_group_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create files in different "packages" (directories)
	createTestGoFile(t, tempDir, "main.go", "package main")
	createTestGoFile(t, tempDir, "main_test.go", "package main")

	subDir := filepath.Join(tempDir, "utils")
	os.MkdirAll(subDir, 0755)
	createTestGoFile(t, subDir, "helper.go", "package utils")
	createTestGoFile(t, subDir, "helper_test.go", "package utils")

	ctx := context.Background()
	d := NewDiscoverer(ctx, tempDir)

	matcher, err := NewPatternMatcher("all_go", `.*\.go$`, "")
	if err != nil {
		t.Fatalf("Failed to create matcher: %v", err)
	}

	files, err := d.DiscoverFilesByPattern(*matcher)
	if err != nil {
		t.Fatalf("Failed to discover files: %v", err)
	}

	groups := d.GroupByPackage(files)

	// We expect at least 2 groups (main and utils), but package detection might vary
	if len(groups) == 0 {
		t.Error("Expected at least one group")
	}

	totalFiles := 0
	for pkg, pkgFiles := range groups {
		if pkg == "" {
			t.Error("Package name should not be empty")
		}
		if len(pkgFiles) == 0 {
			t.Error("Package group should have files")
		}
		totalFiles += len(pkgFiles)
	}

	if totalFiles != len(files) {
		t.Errorf("Total grouped files %d doesn't match discovered files %d", totalFiles, len(files))
	}
}

// TestFilterFilesWithTempDir tests filtering functionality
func TestFilterFilesWithTempDir(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "discovery_filter_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	createTestGoFile(t, tempDir, "main.go", "package main")
	createTestGoFile(t, tempDir, "utils.go", "package main")
	createTestGoFile(t, tempDir, "main_test.go", "package main")

	ctx := context.Background()
	d := NewDiscoverer(ctx, tempDir)

	matcher, err := NewPatternMatcher("all_go", `.*\.go$`, "")
	if err != nil {
		t.Fatalf("Failed to create matcher: %v", err)
	}

	files, err := d.DiscoverFilesByPattern(*matcher)
	if err != nil {
		t.Fatalf("Failed to discover files: %v", err)
	}

	// Test filtering with regex pattern
	filtered := d.FilterFiles(files, ".*main.*")
	if len(filtered) != len(files) {
		t.Errorf("Expected all files to match '.*main.*' pattern, got %d filtered from %d", len(filtered), len(files))
	}

	// Test filtering with non-matching pattern
	filtered2 := d.FilterFiles(files, "nonexistent")
	if len(filtered2) != 0 {
		t.Errorf("Expected no files to match 'nonexistent' pattern, got %d", len(filtered2))
	}
}

// createTestGoFile creates a test Go file with given content
func createTestGoFile(t *testing.T, dir, filename, content string) {
	fullPath := filepath.Join(dir, filename)
	err := os.WriteFile(fullPath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file %s: %v", fullPath, err)
	}
}

// TestFileSystemDiscovery tests the file system discovery method specifically
func TestFileSystemDiscovery(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "fs_discovery_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files
	createTestGoFile(t, tempDir, "main.go", "package main")
	createTestGoFile(t, tempDir, "utils.go", "package main")
	createTestGoFile(t, tempDir, "main_test.go", "package main")
	createTestGoFile(t, tempDir, "data.txt", "not a go file")

	// Create hidden file (should be skipped)
	createTestGoFile(t, tempDir, ".hidden.go", "package main")

	// Create file in .git directory (should be skipped)
	gitDir := filepath.Join(tempDir, ".git")
	os.MkdirAll(gitDir, 0755)
	createTestGoFile(t, gitDir, "config.go", "package main")

	ctx := context.Background()
	d := NewDiscoverer(ctx, tempDir)

	// Test with valid pattern
	matcher, err := NewPatternMatcher("go_files", `.*\.go$`, "")
	if err != nil {
		t.Fatalf("Failed to create matcher: %v", err)
	}

	files, err := d.discoverFilesByPatternFS(*matcher)
	if err != nil {
		t.Fatalf("Failed to discover files via FS: %v", err)
	}

	// Should find 3 .go files (main.go, utils.go, main_test.go) but not .hidden.go or .git/config.go
	expectedFiles := 3
	if len(files) != expectedFiles {
		t.Errorf("Expected %d Go files, got %d", expectedFiles, len(files))
		for _, file := range files {
			t.Logf("Found file: %s", file.Path)
		}
	}

	// Verify file paths are absolute
	for _, file := range files {
		if !filepath.IsAbs(file.Path) {
			t.Errorf("Expected absolute path, got %s", file.Path)
		}
		if file.RelPath == "" {
			t.Errorf("Expected non-empty relative path for %s", file.Path)
		}
	}
}

// TestContentPatternMatching tests content pattern matching in detail
func TestContentPatternMatching(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "content_match_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test file with various function signatures
	testContent := `package main

import "testing"

// Test function at start of line
func TestMain(t *testing.T) {
	t.Log("main test")
}

// Test function with receiver
func (s *Suite) TestMethod(t *testing.T) {
	t.Log("method test")
}

// Regular function (should not match)
func RegularFunction() {
	fmt.Println("not a test")
}

// Test function with spaces (should not match due to comment filtering)
func NotATest(t *testing.T) {
	// This is not a test function despite having t *testing.T
}

func TestAnother() {
	// Another test function
}
`
	createTestGoFile(t, tempDir, "test_functions.go", testContent)

	ctx := context.Background()
	d := NewDiscoverer(ctx, tempDir)

	// Create file list manually since we know the file exists
	files := []FileMatch{{
		Path:    filepath.Join(tempDir, "test_functions.go"),
		Package: "main",
		RelPath: "test_functions.go",
	}}

	// Test finding Test functions (including methods with receivers)
	matcher, err := NewPatternMatcher("test_funcs", "", `^func\s+(?:\([^)]+\)\s+)?(Test\w+)\s*\(`)
	if err != nil {
		t.Fatalf("Failed to create content matcher: %v", err)
	}

	matches, err := d.FindMatchesInFiles(files, *matcher)
	if err != nil {
		t.Fatalf("Failed to find matches: %v", err)
	}

	// Should find 3 test functions: TestMain, (s *Suite) TestMethod, TestAnother
	expectedMatches := 3
	if len(matches) != expectedMatches {
		t.Errorf("Expected %d test function matches, got %d", expectedMatches, len(matches))
		for i, match := range matches {
			t.Logf("Match %d: %s at line %d", i, match.Match, match.Line)
		}
	}

	// Verify match details
	for _, match := range matches {
		if match.File != filepath.Join(tempDir, "test_functions.go") {
			t.Errorf("Expected file path %s, got %s", filepath.Join(tempDir, "test_functions.go"), match.File)
		}
		if match.Package != "main" {
			t.Errorf("Expected package 'main', got %q", match.Package)
		}
		if match.Line <= 0 {
			t.Errorf("Expected positive line number, got %d", match.Line)
		}
		if len(match.Groups) == 0 {
			t.Error("Expected captured groups from regex")
		}
		if match.Groups[0] == "" {
			t.Error("Expected first group to contain function name")
		}
	}
}

// TestIntegrationScenario tests a complete discovery workflow
func TestIntegrationScenario(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "integration_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a mock Go project structure
	createTestGoFile(t, tempDir, "main.go", `package main

import (
	"fmt"
	"example.com/project/utils"
)

func main() {
	fmt.Println("Hello")
	utils.Helper()
}`)

	utilsDir := filepath.Join(tempDir, "utils")
	os.MkdirAll(utilsDir, 0755)

	createTestGoFile(t, utilsDir, "helper.go", `package utils

// Helper provides utility functions
func Helper() {
	println("helper")
}`)

	createTestGoFile(t, utilsDir, "helper_test.go", `package utils

import "testing"

func TestHelper(t *testing.T) {
	Helper()
}

func TestAnother(t *testing.T) {
	// Another test
}

func BenchmarkHelper(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Helper()
	}
}`)

	ctx := context.Background()
	d := NewDiscoverer(ctx, tempDir)

	// Step 1: Discover packages
	_, err = d.DiscoverPackages()
	if err != nil {
		t.Logf("Package discovery failed (might be OK in temp dir): %v", err)
	}

	// Step 2: Discover test files
	testMatcher, err := NewPatternMatcher("test_files", `.*_test\.go$`, "")
	if err != nil {
		t.Fatalf("Failed to create test matcher: %v", err)
	}

	testFiles, err := d.DiscoverFilesByPattern(*testMatcher)
	if err != nil {
		t.Fatalf("Failed to discover test files: %v", err)
	}

	if len(testFiles) != 1 {
		t.Errorf("Expected 1 test file, got %d", len(testFiles))
	}

	// Step 3: Find test functions
	testFuncMatcher, err := NewPatternMatcher("test_funcs", "", `^func\s+(Test\w+)\s*\(`)
	if err != nil {
		t.Fatalf("Failed to create test function matcher: %v", err)
	}

	testFuncs, err := d.FindMatchesInFiles(testFiles, *testFuncMatcher)
	if err != nil {
		t.Fatalf("Failed to find test functions: %v", err)
	}

	if len(testFuncs) != 2 {
		t.Errorf("Expected 2 test functions, got %d", len(testFuncs))
	}

	// Step 4: Find benchmark functions
	benchMatcher, err := NewPatternMatcher("benchmarks", "", `^func\s+(Benchmark\w+)\s*\(`)
	if err != nil {
		t.Fatalf("Failed to create benchmark matcher: %v", err)
	}

	benchmarks, err := d.FindMatchesInFiles(testFiles, *benchMatcher)
	if err != nil {
		t.Fatalf("Failed to find benchmarks: %v", err)
	}

	if len(benchmarks) != 1 {
		t.Errorf("Expected 1 benchmark function, got %d", len(benchmarks))
	}

	// Step 5: Group by package
	groups := d.GroupByPackage(testFiles)
	if len(groups) == 0 {
		t.Error("Expected at least one package group")
	}

	// Step 6: Filter files
	filtered := d.FilterFiles(testFiles, ".*utils.*")
	if len(filtered) != len(testFiles) {
		t.Errorf("Expected all test files to match utils pattern, got %d filtered from %d", len(filtered), len(testFiles))
	}
}

// TestDiscovererRoot tests the GetRoot method
func TestDiscovererRoot(t *testing.T) {
	ctx := context.Background()

	// Test with absolute path
	absPath, err := filepath.Abs(".")
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	d := NewDiscoverer(ctx, absPath)
	if d.GetRoot() != absPath {
		t.Errorf("Expected root %q, got %q", absPath, d.GetRoot())
	}

	// Test with relative path
	d2 := NewDiscoverer(ctx, ".")
	root := d2.GetRoot()
	if root == "" {
		t.Error("Expected non-empty root")
	}

	// Root should be absolute
	if !filepath.IsAbs(root) {
		t.Errorf("Expected absolute root path, got %q", root)
	}
}
