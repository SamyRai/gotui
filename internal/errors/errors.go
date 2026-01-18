package errors

import (
	"errors"
	"fmt"
)

// Common error types for the GoTUI application
var (
	// File system errors
	ErrFileNotFound     = errors.New("file not found")
	ErrFileNotReadable  = errors.New("file not readable")
	ErrDirectoryNotFound = errors.New("directory not found")
	ErrPermissionDenied = errors.New("permission denied")

	// Command execution errors
	ErrCommandFailed     = errors.New("command execution failed")
	ErrCommandNotFound   = errors.New("command not found")
	ErrCommandTimeout    = errors.New("command timed out")
	ErrCommandCancelled  = errors.New("command was cancelled")

	// Parsing errors
	ErrInvalidFormat    = errors.New("invalid format")
	ErrParsingFailed    = errors.New("parsing failed")

	// Configuration errors
	ErrInvalidConfig    = errors.New("invalid configuration")
	ErrConfigNotFound   = errors.New("configuration file not found")

	// UI errors
	ErrUIInitialization = errors.New("UI initialization failed")
	ErrUIRendering      = errors.New("UI rendering failed")

	// Editor errors
	ErrEditorNotFound   = errors.New("editor not found")
	ErrEditorFailed     = errors.New("editor failed to start")
)

// WrappedError provides context around an error
type WrappedError struct {
	Op   string // Operation that failed
	Err  error  // Underlying error
	Msg  string // Additional context
	Path string // File path if relevant
}

func (e *WrappedError) Error() string {
	if e.Path != "" {
		return fmt.Sprintf("%s %s: %s (%v)", e.Op, e.Path, e.Msg, e.Err)
	}
	if e.Msg != "" {
		return fmt.Sprintf("%s: %s (%v)", e.Op, e.Msg, e.Err)
	}
	return fmt.Sprintf("%s: %v", e.Op, e.Err)
}

func (e *WrappedError) Unwrap() error {
	return e.Err
}

// Wrap wraps an error with additional context
func Wrap(op string, err error) error {
	if err == nil {
		return nil
	}
	return &WrappedError{
		Op:  op,
		Err: err,
	}
}

// WrapWithMsg wraps an error with operation and message
func WrapWithMsg(op string, err error, msg string) error {
	if err == nil {
		return nil
	}
	return &WrappedError{
		Op:  op,
		Err: err,
		Msg: msg,
	}
}

// WrapWithPath wraps an error with operation, message, and file path
func WrapWithPath(op string, err error, msg, path string) error {
	if err == nil {
		return nil
	}
	return &WrappedError{
		Op:   op,
		Err:  err,
		Msg:  msg,
		Path: path,
	}
}

// IsTemporary checks if an error is temporary and the operation could be retried
func IsTemporary(err error) bool {
	return errors.Is(err, ErrCommandTimeout) || errors.Is(err, ErrCommandCancelled)
}

// IsNotFound checks if an error indicates a resource was not found
func IsNotFound(err error) bool {
	return errors.Is(err, ErrFileNotFound) || errors.Is(err, ErrDirectoryNotFound) ||
		   errors.Is(err, ErrCommandNotFound) || errors.Is(err, ErrConfigNotFound)
}

// IsPermission checks if an error indicates a permission issue
func IsPermission(err error) bool {
	return errors.Is(err, ErrPermissionDenied)
}

// New creates a new error with the given message
func New(msg string) error {
	return errors.New(msg)
}

// Join joins multiple errors into a single error
func Join(errs ...error) error {
	return errors.Join(errs...)
}

// As finds the first error in err's chain that matches target
func As(err error, target any) bool {
	return errors.As(err, target)
}

// Is reports whether any error in err's chain matches target
func Is(err, target error) bool {
	return errors.Is(err, target)
}