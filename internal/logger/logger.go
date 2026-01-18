package logger

import (
	"io"
	"log/slog"
	"os"
)

// Config holds logger configuration
type Config struct {
	Level      slog.Level
	Output     io.Writer
	DebugFile  string
	EnableFile bool
	AddSource  bool
	JSON       bool
}

// New creates a new slog logger with the given configuration
func New(config Config) *slog.Logger {
	var output io.Writer = os.Stderr
	if config.Output != nil {
		output = config.Output
	}

	// Create file writer if enabled
	if config.EnableFile && config.DebugFile != "" {
		file, err := os.OpenFile(config.DebugFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			output = io.MultiWriter(output, file)
		}
	}

	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level:     config.Level,
		AddSource: config.AddSource,
	}

	if config.JSON {
		handler = slog.NewJSONHandler(output, opts)
	} else {
		handler = slog.NewTextHandler(output, opts)
	}

	return slog.New(handler)
}

// DefaultConfig returns a default logger configuration
func DefaultConfig() Config {
	return Config{
		Level:      slog.LevelInfo,
		Output:     os.Stderr,
		DebugFile:  "",
		EnableFile: false,
		AddSource:  false,
		JSON:       false,
	}
}

// DevelopmentConfig returns a logger configuration suitable for development
func DevelopmentConfig() Config {
	return Config{
		Level:      slog.LevelDebug,
		Output:     os.Stderr,
		DebugFile:  "debug.log",
		EnableFile: true,
		AddSource:  true,
		JSON:       false,
	}
}

// ProductionConfig returns a logger configuration suitable for production
func ProductionConfig() Config {
	return Config{
		Level:      slog.LevelInfo,
		Output:     os.Stderr,
		DebugFile:  "",
		EnableFile: false,
		AddSource:  false,
		JSON:       true,
	}
}

// SetGlobalLogger sets the global slog logger
func SetGlobalLogger(logger *slog.Logger) {
	slog.SetDefault(logger)
}

// init initializes the global logger with default configuration
func init() {
	SetGlobalLogger(New(DefaultConfig()))
}