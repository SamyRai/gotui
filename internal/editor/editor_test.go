package editor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetConfig(t *testing.T) {
	// Save original environment
	originalEditor := os.Getenv("EDITOR")
	originalVisual := os.Getenv("VISUAL")
	defer func() {
		os.Setenv("EDITOR", originalEditor)
		os.Setenv("VISUAL", originalVisual)
	}()

	tests := []struct {
		name         string
		editorEnv    string
		visualEnv    string
		expectedCmd  string
		expectedArgs []string
	}{
		{
			name:         "EDITOR environment variable set",
			editorEnv:    "vim",
			visualEnv:    "",
			expectedCmd:  "vim",
			expectedArgs: []string{},
		},
		{
			name:         "VISUAL environment variable set when EDITOR is not",
			editorEnv:    "",
			visualEnv:    "nano",
			expectedCmd:  "nano",
			expectedArgs: []string{},
		},
		{
			name:         "VISUAL takes precedence over EDITOR",
			editorEnv:    "vim",
			visualEnv:    "nano",
			expectedCmd:  "nano",
			expectedArgs: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("EDITOR", tt.editorEnv)
			os.Setenv("VISUAL", tt.visualEnv)

			config := GetConfig()

			if config.Command != tt.expectedCmd {
				t.Errorf("Expected command %q, got %q", tt.expectedCmd, config.Command)
			}

			if len(config.Args) != len(tt.expectedArgs) {
				t.Errorf("Expected args length %d, got %d", len(tt.expectedArgs), len(config.Args))
			}
		})
	}
}

func TestParseFileLocation(t *testing.T) {
	tests := []struct {
		name         string
		location     string
		expectedFile string
		expectedLine  int
		expectedCol  int
		expectError  bool
	}{
		{
			name:         "file only",
			location:     "main.go",
			expectedFile: "main.go",
			expectedLine:  1,
			expectedCol:  1,
			expectError:  false,
		},
		{
			name:         "file with line",
			location:     "main.go:42",
			expectedFile: "main.go",
			expectedLine:  42,
			expectedCol:  1,
			expectError:  false,
		},
		{
			name:         "file with line and column",
			location:     "main.go:42:10",
			expectedFile: "main.go",
			expectedLine:  42,
			expectedCol:  10,
			expectError:  false,
		},
		{
			name:         "empty string",
			location:     "",
			expectedFile: "",
			expectedLine:  0,
			expectedCol:  0,
			expectError:  true,
		},
		{
			name:         "invalid line number",
			location:     "main.go:abc",
			expectedFile: "main.go",
			expectedLine:  1,
			expectedCol:  1,
			expectError:  false,
		},
		{
			name:         "negative line number",
			location:     "main.go:-1",
			expectedFile: "main.go",
			expectedLine:  1,
			expectedCol:  1,
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, line, col, err := ParseFileLocation(tt.location)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if file != tt.expectedFile {
				t.Errorf("Expected file %q, got %q", tt.expectedFile, file)
			}

			if line != tt.expectedLine {
				t.Errorf("Expected line %d, got %d", tt.expectedLine, line)
			}

			if col != tt.expectedCol {
				t.Errorf("Expected column %d, got %d", tt.expectedCol, col)
			}
		})
	}
}

func TestFormatFileLocation(t *testing.T) {
	tests := []struct {
		name     string
		file     string
		line     int
		column   int
		expected string
	}{
		{
			name:     "file only",
			file:     "main.go",
			line:     1,
			column:   1,
			expected: "main.go",
		},
		{
			name:     "file with line only",
			file:     "main.go",
			line:     42,
			column:   1,
			expected: "main.go:42",
		},
		{
			name:     "file with line and column",
			file:     "main.go",
			line:     42,
			column:   10,
			expected: "main.go:42:10",
		},
		{
			name:     "zero line and column",
			file:     "main.go",
			line:     0,
			column:   0,
			expected: "main.go",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatFileLocation(tt.file, tt.line, tt.column)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestDetectProjectRoot(t *testing.T) {
	// Create a temporary directory structure
	tempDir, err := os.MkdirTemp("", "gotui_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create a subdirectory
	subDir := filepath.Join(tempDir, "subdir")
	err = os.MkdirAll(subDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	// Create project indicators
	indicators := []string{"go.mod", "go.sum", ".git", "Makefile"}
	for _, indicator := range indicators {
		filePath := filepath.Join(tempDir, indicator)
		file, err := os.Create(filePath)
		if err != nil {
			t.Fatal(err)
		}
		file.Close()
	}

	tests := []struct {
		name        string
		startPath   string
		expectError bool
	}{
		{
			name:        "start from subdirectory",
			startPath:   subDir,
			expectError: false,
		},
		{
			name:        "start from project root",
			startPath:   tempDir,
			expectError: false,
		},
		{
			name:        "start from non-existent path",
			startPath:   "/non/existent/path",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := DetectProjectRoot(tt.startPath)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Normalize paths for comparison
			expected, err := filepath.Abs(tempDir)
			if err != nil {
				t.Fatal(err)
			}

			resultAbs, err := filepath.Abs(result)
			if err != nil {
				t.Fatal(err)
			}

			if resultAbs != expected {
				t.Errorf("Expected project root %q, got %q", expected, resultAbs)
			}
		})
	}
}

func TestDetectProjectRoot_NoIndicators(t *testing.T) {
	// Create a temporary directory with no project indicators
	tempDir, err := os.MkdirTemp("", "gotui_test_no_indicators")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// The function should return the original path if no indicators are found
	result, err := DetectProjectRoot(tempDir)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected, err := filepath.Abs(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	resultAbs, err := filepath.Abs(result)
	if err != nil {
		t.Fatal(err)
	}

	if resultAbs != expected {
		t.Errorf("Expected %q, got %q", expected, resultAbs)
	}
}

// TestOpenAtLine_NonExistentFile tests error handling for non-existent files
func TestOpenAtLine_NonExistentFile(t *testing.T) {
	err := OpenAtLine("/non/existent/file.go", 1, 1)
	if err == nil {
		t.Error("Expected error for non-existent file")
	}

	if !strings.Contains(err.Error(), "file does not exist") {
		t.Errorf("Expected 'file does not exist' error, got: %v", err)
	}
}

// TestOpenDirectory_NonExistentDirectory tests error handling for non-existent directories
func TestOpenDirectory_NonExistentDirectory(t *testing.T) {
	err := OpenDirectory("/non/existent/directory")
	if err == nil {
		t.Error("Expected error for non-existent directory")
	}

	if !strings.Contains(err.Error(), "directory does not exist") {
		t.Errorf("Expected 'directory does not exist' error, got: %v", err)
	}
}