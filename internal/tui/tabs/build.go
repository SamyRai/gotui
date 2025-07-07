package tabs

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
)

// BuildRunner manages build execution
type BuildRunner struct {
	ctx    context.Context
	width  int
	height int
}

// NewBuildRunner creates a new build runner
func NewBuildRunner(ctx context.Context) *BuildRunner {
	return &BuildRunner{
		ctx: ctx,
	}
}

// Init initializes the build runner
func (br BuildRunner) Init() tea.Cmd {
	return nil
}

// SetSize sets the dimensions
func (br *BuildRunner) SetSize(width, height int) {
	br.width = width
	br.height = height
}

// Update handles messages
func (br *BuildRunner) Update(msg tea.Msg) (TabInterface, tea.Cmd) {
	return br, nil
}

// View renders the build runner
func (br BuildRunner) View() string {
	return "Build tab - Coming soon!"
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
