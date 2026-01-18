package editor

import (
	"fmt"
	"goutui/internal/errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Config holds configuration for opening files in editors
type Config struct {
	Command string
	Args    []string
}

// GetConfig returns the appropriate editor configuration
func GetConfig() Config {
	editor := os.Getenv("VISUAL") // VISUAL takes precedence over EDITOR
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}
	if editor == "" {
		// Try to detect common editors
		for _, cmd := range []string{"code", "subl", "atom", "vim", "nvim", "nano"} {
			if _, err := exec.LookPath(cmd); err == nil {
				editor = cmd
				break
			}
		}
	}
	if editor == "" {
		editor = "vim" // Fallback
	}

	return Config{
		Command: editor,
		Args:    []string{},
	}
}

// OpenAtLine opens a file at a specific line and column in the user's editor
func OpenAtLine(filename string, line int, column int) error {
	// Convert to absolute path
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return errors.WrapWithPath("absolute path", err, "failed to get absolute path", filename)
	}

	// Check if file exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return errors.WrapWithPath("file check", errors.ErrFileNotFound, "file does not exist", absPath)
	}

	config := GetConfig()
	editor := config.Command

	var cmd *exec.Cmd

	// Handle different editors
	switch {
	case strings.Contains(editor, "code") || strings.HasSuffix(editor, "code"):
		// Visual Studio Code
		cmd = exec.Command(editor, "-g", fmt.Sprintf("%s:%d:%d", absPath, line, column))
	case strings.Contains(editor, "subl") || strings.HasSuffix(editor, "subl"):
		// Sublime Text
		cmd = exec.Command(editor, fmt.Sprintf("%s:%d:%d", absPath, line, column))
	case strings.Contains(editor, "atom"):
		// Atom
		cmd = exec.Command(editor, fmt.Sprintf("%s:%d:%d", absPath, line, column))
	case strings.Contains(editor, "vim") || strings.Contains(editor, "nvim"):
		// Vim/Neovim
		cmd = exec.Command(editor, fmt.Sprintf("+%d", line), absPath)
	case strings.Contains(editor, "emacs"):
		// Emacs
		cmd = exec.Command(editor, fmt.Sprintf("+%d:%d", line, column), absPath)
	case strings.Contains(editor, "nano"):
		// Nano
		cmd = exec.Command(editor, fmt.Sprintf("+%d,%d", line, column), absPath)
	case strings.Contains(editor, "gedit"):
		// Gedit
		cmd = exec.Command(editor, fmt.Sprintf("+%d", line), absPath)
	default:
		// Generic fallback - just open the file
		cmd = exec.Command(editor, absPath)
	}

	// Set up the command to run in the background
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Start() // Use Start() instead of Run() to not block
}

// OpenFile opens a file without specific line/column
func OpenFile(filename string) error {
	return OpenAtLine(filename, 1, 1)
}

// OpenDirectory opens a directory in the editor
func OpenDirectory(dirname string) error {
	// Convert to absolute path
	absPath, err := filepath.Abs(dirname)
	if err != nil {
		return errors.WrapWithPath("absolute path", err, "failed to get absolute path", dirname)
	}

	// Check if directory exists
	if info, err := os.Stat(absPath); os.IsNotExist(err) || !info.IsDir() {
		return errors.WrapWithPath("directory check", errors.ErrDirectoryNotFound, "directory does not exist", absPath)
	}

	config := GetConfig()
	editor := config.Command

	var cmd *exec.Cmd

	// Handle different editors
	switch {
	case strings.Contains(editor, "code") || strings.HasSuffix(editor, "code"):
		// Visual Studio Code
		cmd = exec.Command(editor, absPath)
	case strings.Contains(editor, "subl") || strings.HasSuffix(editor, "subl"):
		// Sublime Text
		cmd = exec.Command(editor, absPath)
	case strings.Contains(editor, "atom"):
		// Atom
		cmd = exec.Command(editor, absPath)
	default:
		// For terminal editors, we can't really open a directory
		// So we'll open the current directory file
		cmd = exec.Command(editor, absPath)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Start()
}

// DetectProjectRoot tries to find the project root directory
func DetectProjectRoot(startPath string) (string, error) {
	// Check if the start path exists
	if _, err := os.Stat(startPath); os.IsNotExist(err) {
		return "", errors.WrapWithPath("path check", errors.ErrDirectoryNotFound, "start path does not exist", startPath)
	}

	currentPath, err := filepath.Abs(startPath)
	if err != nil {
		return "", errors.WrapWithPath("absolute path", err, "failed to get absolute path", startPath)
	}

	// Look for common project indicators
	indicators := []string{
		"go.mod",
		"go.sum",
		".git",
		"Makefile",
		"README.md",
		"package.json",
		"Cargo.toml",
	}

	for {
		for _, indicator := range indicators {
			if _, err := os.Stat(filepath.Join(currentPath, indicator)); err == nil {
				return currentPath, nil
			}
		}

		parent := filepath.Dir(currentPath)
		if parent == currentPath {
			// Reached the root directory
			break
		}
		currentPath = parent
	}

	// If no project root found, return the original path
	return startPath, nil
}

// ParseFileLocation parses a file location string like "file.go:10:5"
func ParseFileLocation(location string) (file string, line int, column int, err error) {
	if location == "" {
		return "", 0, 0, errors.WrapWithMsg("parse location", errors.ErrInvalidFormat, "empty location string")
	}

	parts := strings.Split(location, ":")
	if len(parts) == 0 || parts[0] == "" {
		return "", 0, 0, errors.WrapWithMsg("parse location", errors.ErrInvalidFormat, "invalid location format")
	}

	file = parts[0]
	line = 1
	column = 1

	if len(parts) > 1 {
		if l := parseInt(parts[1]); l > 0 {
			line = l
		}
	}

	if len(parts) > 2 {
		if c := parseInt(parts[2]); c > 0 {
			column = c
		}
	}

	return file, line, column, nil
}

// parseInt safely parses an integer from a string
func parseInt(s string) int {
	var result int
	for _, r := range s {
		if r >= '0' && r <= '9' {
			result = result*10 + int(r-'0')
		} else {
			break
		}
	}
	return result
}

// FormatFileLocation formats a file location for display
func FormatFileLocation(file string, line int, column int) string {
	if line <= 1 && column <= 1 {
		return file
	}
	if column <= 1 {
		return fmt.Sprintf("%s:%d", file, line)
	}
	return fmt.Sprintf("%s:%d:%d", file, line, column)
}