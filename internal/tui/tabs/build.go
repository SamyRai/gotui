package tabs

import (
	"context"
	"goutui/internal/tui/components"
	"goutui/internal/style"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// BuildRunner manages build execution
type BuildRunner struct {
	ctx       context.Context
	width     int
	height    int
	actionBar components.ActionBar
}

// NewBuildRunner creates a new build runner
func NewBuildRunner(ctx context.Context) *BuildRunner {
	br := &BuildRunner{
		ctx:       ctx,
		actionBar: components.NewActionBar(),
	}
	br.updateActionBar()
	return br
}

// updateActionBar updates the action bar
func (br *BuildRunner) updateActionBar() {
	br.actionBar.Clear()
	br.actionBar.AddAction(components.Action{
		Key:         "r",
		Label:       "Build",
		Description: "Build the project",
		Primary:     true,
	})
	br.actionBar.AddAction(components.Action{
		Key:         "c",
		Label:       "Clean",
		Description: "Clean build artifacts",
		Primary:     false,
	})
}

// Init initializes the build runner
func (br BuildRunner) Init() tea.Cmd {
	return nil
}

// SetSize sets the dimensions
func (br *BuildRunner) SetSize(width, height int) {
	br.width = width
	br.height = height
	br.actionBar.SetWidth(width)
}

// Update handles messages
func (br *BuildRunner) Update(msg tea.Msg) (TabInterface, tea.Cmd) {
	return br, nil
}

// View renders the build runner
func (br BuildRunner) View() string {
	actionBarView := br.actionBar.View()
	
	var parts []string
	if actionBarView != "" {
		parts = append(parts, actionBarView, "")
	}
	
	header := style.HeaderStyle.Render("Build")
	parts = append(parts, header, "")
	
	parts = append(parts, style.SubtleStyle.Render("Build functionality coming soon!"))
	
	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

// Refresh triggers a refresh
func (br BuildRunner) Refresh() tea.Cmd {
	return nil
}

// GetStatus returns the current status
func (br BuildRunner) GetStatus() string {
	return "Build ready"
}

// GetKeyHints returns key binding hints
func (br BuildRunner) GetKeyHints() []string {
	return []string{"r: build", "c: clean"}
}

// Cleanup performs cleanup
func (br BuildRunner) Cleanup() {
}
