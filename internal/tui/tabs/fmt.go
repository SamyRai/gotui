package tabs

import (
	"context"
	"goutui/internal/runner"
	"goutui/internal/style"
	"goutui/internal/tui/components"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// FmtDiff manages format checking and diff display
type FmtDiff struct {
	ctx        context.Context
	width      int
	height     int
	runner     *runner.CommandRunner
	diffViewer components.DiffViewer
	running    bool
	hasChanges bool
	status     string
}

// FmtKeyMap defines key bindings
type FmtKeyMap struct {
	Check      key.Binding
	AutoFormat key.Binding
	OpenFile   key.Binding
	Refresh    key.Binding
}

// DefaultFmtKeyMap returns default key bindings
func DefaultFmtKeyMap() FmtKeyMap {
	return FmtKeyMap{
		Check: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "check format"),
		),
		AutoFormat: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "auto-format"),
		),
		OpenFile: key.NewBinding(
			key.WithKeys("o"),
			key.WithHelp("o", "open file"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("ctrl+r"),
			key.WithHelp("ctrl+r", "refresh"),
		),
	}
}

// NewFmtDiff creates a new fmt diff checker
func NewFmtDiff(ctx context.Context) *FmtDiff {
	return &FmtDiff{
		ctx:        ctx,
		runner:     runner.NewCommandRunner(ctx),
		diffViewer: components.NewDiffViewer("Format Diff"),
		status:     "Ready to check format",
	}
}

// Init initializes the fmt diff checker
func (fd FmtDiff) Init() tea.Cmd {
	return nil
}

// SetSize sets the dimensions
func (fd *FmtDiff) SetSize(width, height int) {
	fd.width = width
	fd.height = height
	// Account for header lines (title, status, spacing lines = ~5 lines)
	contentHeight := height - 5
	if contentHeight < 1 {
		contentHeight = 1
	}
	fd.diffViewer.SetSize(width, contentHeight)
}

// runFmtCheck executes gofmt -d ./...
func (fd *FmtDiff) runFmtCheck() tea.Cmd {
	fd.running = true
	fd.status = "Checking format..."
	return fd.runner.Run("gofmt", "-d", "./...")
}

// runAutoFormat executes gofmt -w ./...
func (fd *FmtDiff) runAutoFormat() tea.Cmd {
	fd.running = true
	fd.status = "Auto-formatting..."
	return fd.runner.Run("gofmt", "-w", "./...")
}

// handleFmtOutput processes output from gofmt -d
func (fd *FmtDiff) handleFmtOutput(output runner.CommandOutput) tea.Cmd {
	// Accumulate diff output for full display
	if output.Line != "" {
		currentDiff := fd.diffViewer.GetRawDiff()
		if currentDiff == "" {
			fd.diffViewer.SetDiff(output.Line)
		} else {
			fd.diffViewer.SetDiff(currentDiff + "\n" + output.Line)
		}
	}
	return nil
}

// Update handles messages
func (fd *FmtDiff) Update(msg tea.Msg) (TabInterface, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !fd.running {
			keyMap := DefaultFmtKeyMap()
			switch {
			case key.Matches(msg, keyMap.Check):
				return fd, fd.runFmtCheck()
			case key.Matches(msg, keyMap.AutoFormat):
				return fd, fd.runAutoFormat()
			case key.Matches(msg, keyMap.OpenFile):
				// TODO(medium, 4h): Implement file opening from diff
				// When user presses 'o', open the file at the diff location. Requires parsing diff output and integrating with editor logic.
				// Open file from diff (basic implementation)
				filePath := fd.diffViewer.GetFirstFilePathFromDiff()
				if filePath != "" {
					// Use util/editor.go logic to open file (pseudo-code)
					// util.OpenFileInEditor(filePath)
				}
			}
		}

	case runner.CommandOutput:
		if fd.running {
			cmd := fd.handleFmtOutput(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}

	case runner.CommandFinished:
		if fd.running {
			fd.running = false
			if msg.Result.ExitCode == 0 {
				fd.status = "No formatting changes needed"
				fd.hasChanges = false
			} else {
				fd.status = "Format changes detected"
				fd.hasChanges = true
			}
		}
	}

	// Update diff viewer
	var cmd tea.Cmd
	fd.diffViewer, cmd = fd.diffViewer.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return fd, tea.Batch(cmds...)
}

// View renders the fmt diff checker
func (fd FmtDiff) View() string {
	// Create header styled consistently with other tabs
	header := lipgloss.NewStyle().
		Foreground(style.AccentColor).
		Bold(true).
		Render("Format Checker")

	// Create status line
	statusText := fd.status
	if fd.running {
		statusText = lipgloss.NewStyle().
			Foreground(style.InfoColor).
			Render(statusText)
	}

	// Get clean content from diff viewer
	content := fd.diffViewer.GetContent()

	// Create the complete view using lipgloss layout
	return lipgloss.JoinVertical(
		lipgloss.Left,
		"",                    // Empty line for spacing
		header,                // Title
		"",                    // Empty line
		statusText,            // Status
		"",                    // Empty line
		content,               // Main content
	)
}

// Refresh triggers a refresh
func (fd FmtDiff) Refresh() tea.Cmd {
	return fd.runFmtCheck()
}

// GetStatus returns the current status
func (fd FmtDiff) GetStatus() string {
	return fd.status
}

// GetKeyHints returns key binding hints
func (fd FmtDiff) GetKeyHints() []string {
	if fd.running {
		return []string{"running..."}
	}
	return []string{"r: check format", "a: auto-format", "o: open file"}
}

// Cleanup performs cleanup
func (fd FmtDiff) Cleanup() {
	if fd.runner != nil {
		fd.runner.Stop()
	}
}
