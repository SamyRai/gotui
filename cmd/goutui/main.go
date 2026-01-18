package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"goutui/internal/logger"
	"goutui/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if err := run(); err != nil {
		slog.Error("Error running program", "error", err)
		os.Exit(1)
	}
}

func run() error {
	// Configure logging based on environment
	if os.Getenv("GOTUI_DEBUG") == "1" || os.Getenv("DEBUG") == "1" {
		logger.SetGlobalLogger(logger.New(logger.DevelopmentConfig()))
		slog.Debug("Starting GoTUI in debug mode")
	} else {
		logger.SetGlobalLogger(logger.New(logger.ProductionConfig()))
	}

	slog.Info("Initializing GoTUI")

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create the TUI program
	p := tea.NewProgram(tui.NewModel(), tea.WithAltScreen())

	// Run the program in a goroutine
	errChan := make(chan error, 1)
	go func() {
		_, err := p.Run()
		errChan <- err
	}()

	// Wait for either program completion or signal
	select {
	case err := <-errChan:
		if err != nil {
			slog.Error("Program exited with error", "error", err)
			return err
		}
		slog.Info("Program exited successfully")
		return nil
	case sig := <-sigChan:
		slog.Info("Received signal, shutting down gracefully", "signal", sig)
		p.Quit()
		return nil
	}
}
