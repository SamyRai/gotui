package main

import (
	"fmt"
	"log"
	"os"

	"goutui/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}
}

func run() error {
	// TODO: Add proper logging to a file.
	// For now, we can disable the default logger to avoid cluttering the TUI.
	// log.SetOutput(io.Discard)

	p := tea.NewProgram(tui.NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		return err
	}
	return nil
}
