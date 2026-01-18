package errors

import (
	"errors"
	"strings"
	"testing"
)

func TestWrappedError_Error(t *testing.T) {
	tests := []struct {
		name     string
		wrapped  *WrappedError
		expected string
	}{
		{
			name: "operation and error only",
			wrapped: &WrappedError{
				Op:  "test operation",
				Err: errors.New("underlying error"),
			},
			expected: "test operation: underlying error",
		},
		{
			name: "operation, message, and error",
			wrapped: &WrappedError{
				Op:  "file read",
				Err: errors.New("permission denied"),
				Msg: "failed to open file",
			},
			expected: "file read: failed to open file (permission denied)",
		},
		{
			name: "operation, message, path, and error",
			wrapped: &WrappedError{
				Op:   "file open",
				Err:  errors.New("not found"),
				Msg:  "file missing",
				Path: "/tmp/test.txt",
			},
			expected: "file open /tmp/test.txt: file missing (not found)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.wrapped.Error()
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestWrappedError_Unwrap(t *testing.T) {
	underlying := errors.New("underlying error")
	wrapped := &WrappedError{
		Op:  "test",
		Err: underlying,
	}

	if wrapped.Unwrap() != underlying {
		t.Error("Unwrap should return the underlying error")
	}
}

func TestWrap(t *testing.T) {
	tests := []struct {
		name     string
		op       string
		err      error
		expected bool // whether result should be nil
	}{
		{
			name:     "wrap non-nil error",
			op:       "test op",
			err:      errors.New("test error"),
			expected: false,
		},
		{
			name:     "wrap nil error",
			op:       "test op",
			err:      nil,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Wrap(tt.op, tt.err)
			if (result == nil) != tt.expected {
				t.Errorf("Expected nil result: %v, got: %v", tt.expected, result == nil)
			}
			if result != nil {
				wrapped, ok := result.(*WrappedError)
				if !ok {
					t.Error("Expected WrappedError")
				}
				if wrapped.Op != tt.op {
					t.Errorf("Expected op %q, got %q", tt.op, wrapped.Op)
				}
			}
		})
	}
}

func TestWrapWithMsg(t *testing.T) {
	err := errors.New("underlying")
	result := WrapWithMsg("operation", err, "additional context")

	if result == nil {
		t.Error("Expected non-nil error")
	}

	wrapped, ok := result.(*WrappedError)
	if !ok {
		t.Error("Expected WrappedError")
	}

	if wrapped.Op != "operation" {
		t.Errorf("Expected op 'operation', got %q", wrapped.Op)
	}

	if wrapped.Msg != "additional context" {
		t.Errorf("Expected msg 'additional context', got %q", wrapped.Msg)
	}
}

func TestWrapWithPath(t *testing.T) {
	err := errors.New("underlying")
	result := WrapWithPath("operation", err, "context", "/path/to/file")

	if result == nil {
		t.Error("Expected non-nil error")
	}

	wrapped, ok := result.(*WrappedError)
	if !ok {
		t.Error("Expected WrappedError")
	}

	if wrapped.Path != "/path/to/file" {
		t.Errorf("Expected path '/path/to/file', got %q", wrapped.Path)
	}
}

func TestIsTemporary(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "timeout error",
			err:      ErrCommandTimeout,
			expected: true,
		},
		{
			name:     "cancelled error",
			err:      ErrCommandCancelled,
			expected: true,
		},
		{
			name:     "other error",
			err:      ErrCommandFailed,
			expected: false,
		},
		{
			name:     "wrapped timeout error",
			err:      Wrap("test", ErrCommandTimeout),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsTemporary(tt.err)
			if result != tt.expected {
				t.Errorf("Expected IsTemporary %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestIsNotFound(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "file not found",
			err:      ErrFileNotFound,
			expected: true,
		},
		{
			name:     "directory not found",
			err:      ErrDirectoryNotFound,
			expected: true,
		},
		{
			name:     "command not found",
			err:      ErrCommandNotFound,
			expected: true,
		},
		{
			name:     "config not found",
			err:      ErrConfigNotFound,
			expected: true,
		},
		{
			name:     "other error",
			err:      ErrCommandFailed,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNotFound(tt.err)
			if result != tt.expected {
				t.Errorf("Expected IsNotFound %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestIsPermission(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "permission denied",
			err:      ErrPermissionDenied,
			expected: true,
		},
		{
			name:     "other error",
			err:      ErrCommandFailed,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPermission(tt.err)
			if result != tt.expected {
				t.Errorf("Expected IsPermission %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestErrorConstants(t *testing.T) {
	// Test that all error constants are properly defined and not nil
	errorConstants := []error{
		ErrFileNotFound,
		ErrFileNotReadable,
		ErrDirectoryNotFound,
		ErrPermissionDenied,
		ErrCommandFailed,
		ErrCommandNotFound,
		ErrCommandTimeout,
		ErrCommandCancelled,
		ErrInvalidFormat,
		ErrParsingFailed,
		ErrInvalidConfig,
		ErrConfigNotFound,
		ErrUIInitialization,
		ErrUIRendering,
		ErrEditorNotFound,
		ErrEditorFailed,
	}

	for _, err := range errorConstants {
		if err == nil {
			t.Error("Error constant should not be nil")
		}
		if err.Error() == "" {
			t.Error("Error constant should have a non-empty error message")
		}
	}
}

func TestJoin(t *testing.T) {
	err1 := errors.New("error 1")
	err2 := errors.New("error 2")

	joined := Join(err1, err2)
	if joined == nil {
		t.Error("Joined error should not be nil")
	}

	errStr := joined.Error()
	if !strings.Contains(errStr, "error 1") || !strings.Contains(errStr, "error 2") {
		t.Errorf("Joined error should contain both error messages: %s", errStr)
	}
}

func TestAsAndIs(t *testing.T) {
	wrappedErr := Wrap("test", ErrFileNotFound)

	// Test Is function
	if !Is(wrappedErr, ErrFileNotFound) {
		t.Error("Is should return true for wrapped ErrFileNotFound")
	}

	if Is(wrappedErr, ErrCommandFailed) {
		t.Error("Is should return false for different error")
	}

	// Test As function
	var target *WrappedError
	if !As(wrappedErr, &target) {
		t.Error("As should succeed for WrappedError")
	}

	if target == nil || target.Op != "test" {
		t.Error("As should correctly unwrap to WrappedError")
	}
}