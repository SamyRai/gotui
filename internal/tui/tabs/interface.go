package tabs

import (
	tea "github.com/charmbracelet/bubbletea"
)

// TabInterface defines the common interface for all tab models
type TabInterface interface {
	// Update handles Bubbletea messages
	Update(msg tea.Msg) (TabInterface, tea.Cmd)
	
	// View renders the tab content
	View() string
	
	// Init initializes the tab
	Init() tea.Cmd
	
	// SetSize sets the dimensions available to the tab
	SetSize(width, height int)
	
	// Refresh triggers a refresh of the tab's data
	Refresh() tea.Cmd
	
	// GetStatus returns the current status text for the status bar
	GetStatus() string
	
	// GetKeyHints returns tab-specific key binding hints
	GetKeyHints() []string
	
	// Cleanup performs any necessary cleanup when the tab is closed
	Cleanup()
}
