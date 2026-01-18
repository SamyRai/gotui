package logger

import (
	"bytes"
	"log/slog"
	"os"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		logLevel slog.Level
		json     bool
	}{
		{
			name: "development config",
			config: DevelopmentConfig(),
			logLevel: slog.LevelDebug,
			json: false,
		},
		{
			name: "production config",
			config: ProductionConfig(),
			logLevel: slog.LevelInfo,
			json: true,
		},
		{
			name: "default config",
			config: DefaultConfig(),
			logLevel: slog.LevelInfo,
			json: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := New(tt.config)
			if logger == nil {
				t.Error("Expected non-nil logger")
			}

			// We can't easily test the internal handler configuration,
			// but we can verify the logger is created successfully
		})
	}
}

func TestConfigStruct(t *testing.T) {
	config := Config{
		Level:      slog.LevelWarn,
		Output:     &bytes.Buffer{},
		DebugFile:  "test.log",
		EnableFile: true,
		AddSource:  true,
		JSON:       true,
	}

	if config.Level != slog.LevelWarn {
		t.Errorf("Expected level %v, got %v", slog.LevelWarn, config.Level)
	}

	if config.DebugFile != "test.log" {
		t.Errorf("Expected debug file 'test.log', got %q", config.DebugFile)
	}

	if !config.EnableFile {
		t.Error("Expected EnableFile to be true")
	}

	if !config.AddSource {
		t.Error("Expected AddSource to be true")
	}

	if !config.JSON {
		t.Error("Expected JSON to be true")
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.Level != slog.LevelInfo {
		t.Errorf("Expected default level to be Info, got %v", config.Level)
	}

	if config.DebugFile != "" {
		t.Errorf("Expected empty debug file, got %q", config.DebugFile)
	}

	if config.EnableFile {
		t.Error("Expected EnableFile to be false by default")
	}

	if config.AddSource {
		t.Error("Expected AddSource to be false by default")
	}

	if config.JSON {
		t.Error("Expected JSON to be false by default")
	}
}

func TestDevelopmentConfig(t *testing.T) {
	config := DevelopmentConfig()

	if config.Level != slog.LevelDebug {
		t.Errorf("Expected development level to be Debug, got %v", config.Level)
	}

	if config.DebugFile != "debug.log" {
		t.Errorf("Expected debug file 'debug.log', got %q", config.DebugFile)
	}

	if !config.EnableFile {
		t.Error("Expected EnableFile to be true in development")
	}

	if !config.AddSource {
		t.Error("Expected AddSource to be true in development")
	}

	if config.JSON {
		t.Error("Expected JSON to be false in development")
	}
}

func TestProductionConfig(t *testing.T) {
	config := ProductionConfig()

	if config.Level != slog.LevelInfo {
		t.Errorf("Expected production level to be Info, got %v", config.Level)
	}

	if config.DebugFile != "" {
		t.Errorf("Expected empty debug file in production, got %q", config.DebugFile)
	}

	if config.EnableFile {
		t.Error("Expected EnableFile to be false in production")
	}

	if config.AddSource {
		t.Error("Expected AddSource to be false in production")
	}

	if !config.JSON {
		t.Error("Expected JSON to be true in production")
	}
}

func TestSetGlobalLogger(t *testing.T) {
	// Save original logger
	originalLogger := slog.Default()

	// Create a test logger
	var buf bytes.Buffer
	testLogger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{}))

	// Set it as global
	SetGlobalLogger(testLogger)

	// Test that it's set (we can't easily verify this without accessing internal state)
	// But we can verify no panic occurs

	// Restore original logger
	slog.SetDefault(originalLogger)
}

func TestLoggerCreationWithFileOutput(t *testing.T) {
	// Create a temporary file for testing
	tempDir := t.TempDir()
	logFile := tempDir + "/test.log"

	config := Config{
		Level:      slog.LevelInfo,
		Output:     &bytes.Buffer{}, // Primary output
		DebugFile:  logFile,
		EnableFile: true,
		AddSource:  false,
		JSON:       false,
	}

	logger := New(config)
	if logger == nil {
		t.Error("Expected non-nil logger")
	}

	// Check if the log file was created
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Error("Expected log file to be created")
	}
}

func TestLoggerCreationWithInvalidFile(t *testing.T) {
	config := Config{
		Level:      slog.LevelInfo,
		Output:     &bytes.Buffer{},
		DebugFile:  "/invalid/path/that/does/not/exist/test.log",
		EnableFile: true,
		AddSource:  false,
		JSON:       false,
	}

	// This should not panic, even with an invalid file path
	logger := New(config)
	if logger == nil {
		t.Error("Expected non-nil logger even with invalid file path")
	}
}

func TestJSONHandler(t *testing.T) {
	var buf bytes.Buffer
	config := Config{
		Level:      slog.LevelInfo,
		Output:     &buf,
		DebugFile:  "",
		EnableFile: false,
		AddSource:  false,
		JSON:       true,
	}

	logger := New(config)

	// Log a test message
	logger.Info("test message", "key", "value")

	output := buf.String()

	// Check if output contains JSON-like structure
	if !strings.Contains(output, `"msg"`) || !strings.Contains(output, `"key"`) {
		t.Errorf("Expected JSON output, got: %s", output)
	}
}

func TestTextHandler(t *testing.T) {
	var buf bytes.Buffer
	config := Config{
		Level:      slog.LevelInfo,
		Output:     &buf,
		DebugFile:  "",
		EnableFile: false,
		AddSource:  false,
		JSON:       false,
	}

	logger := New(config)

	// Log a test message
	logger.Info("test message", "key", "value")

	output := buf.String()

	// Check if output contains text-like structure
	if !strings.Contains(output, "test message") || !strings.Contains(output, "key=value") {
		t.Errorf("Expected text output, got: %s", output)
	}
}