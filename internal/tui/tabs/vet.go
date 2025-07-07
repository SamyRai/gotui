package tabs

import (
	"context"
	"goutui/internal/runner"
	"goutui/internal/style"
	"goutui/internal/util"
	"regexp"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// VetProblem represents a problem found by go vet
type VetProblem struct {
	File    string
	Line    int
	Column  int
	Message string
	Package string
	Raw     string
}

// String returns a formatted string representation
func (vp VetProblem) String() string {
	return vp.Message
}

// FilterValue returns the value to filter on
func (vp VetProblem) FilterValue() string {
	return vp.File + " " + vp.Message
}

// Title returns the title for the list item
func (vp VetProblem) Title() string {
	return vp.Message
}

// Description returns the description for the list item
func (vp VetProblem) Description() string {
	return util.FormatFileLocation(vp.File, vp.Line, vp.Column)
}

// VetRunner manages vet/lint execution
type VetRunner struct {
	ctx      context.Context
	width    int
	height   int
	runner   *runner.CommandRunner
	list     list.Model
	problems []VetProblem
	running  bool
	status   string
}

// VetKeyMap defines key bindings
type VetKeyMap struct {
	Run      key.Binding
	OpenFile key.Binding
	Refresh  key.Binding
}

// DefaultVetKeyMap returns default key bindings
func DefaultVetKeyMap() VetKeyMap {
	return VetKeyMap{
		Run: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "run vet"),
		),
		OpenFile: key.NewBinding(
			key.WithKeys("o", "enter"),
			key.WithHelp("o/enter", "open file"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("ctrl+r"),
			key.WithHelp("ctrl+r", "refresh"),
		),
	}
}

// NewVetRunner creates a new vet runner
func NewVetRunner(ctx context.Context) *VetRunner {
	// Create list for problems
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Vet Issues"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = style.HeaderStyle
	l.Styles.PaginationStyle = style.SubtleStyle
	l.Styles.HelpStyle = style.SubtleStyle

	return &VetRunner{
		ctx:    ctx,
		runner: runner.NewCommandRunner(ctx),
		list:   l,
		status: "Ready to run vet",
	}
}

// Init initializes the vet runner
func (vr VetRunner) Init() tea.Cmd {
	return nil
}

// SetSize sets the dimensions
func (vr *VetRunner) SetSize(width, height int) {
	vr.width = width
	vr.height = height
	vr.list.SetWidth(width)
	vr.list.SetHeight(height - 4) // Account for header and status
}

// runVet executes go vet ./...
func (vr *VetRunner) runVet() tea.Cmd {
	vr.running = true
	vr.status = "Running vet..."
	vr.problems = []VetProblem{}
	vr.updateList()
	return vr.runner.Run("go", "vet", "./...")
}

// parseVetLine parses a vet output line
// Example: ./file.go:10:5: error message
var vetLineRegex = regexp.MustCompile(`^(.+\.go):(\d+):(\d+):\s*(.+)$`)

func (vr *VetRunner) parseVetLine(line string) *VetProblem {
	matches := vetLineRegex.FindStringSubmatch(line)
	if len(matches) < 5 {
		return nil
	}

	lineNum, _ := strconv.Atoi(matches[2])
	colNum, _ := strconv.Atoi(matches[3])

	return &VetProblem{
		File:    matches[1],
		Line:    lineNum,
		Column:  colNum,
		Message: matches[4],
		Raw:     line,
	}
}

// handleVetOutput processes output from go vet
func (vr *VetRunner) handleVetOutput(output runner.CommandOutput) tea.Cmd {
	if problem := vr.parseVetLine(output.Line); problem != nil {
		vr.problems = append(vr.problems, *problem)
		vr.updateList()
	}
	return nil
}

// updateList updates the list with current problems
func (vr *VetRunner) updateList() {
	items := make([]list.Item, len(vr.problems))
	for i, problem := range vr.problems {
		items[i] = problem
	}
	vr.list.SetItems(items)
}

// openSelectedFile opens the currently selected problem file
func (vr *VetRunner) openSelectedFile() tea.Cmd {
	if selected, ok := vr.list.SelectedItem().(VetProblem); ok {
		return func() tea.Msg {
			err := util.OpenInEditor(selected.File, selected.Line, selected.Column)
			if err != nil {
				return runner.CommandOutput{Line: "Error opening file: " + err.Error()}
			}
			return runner.CommandOutput{Line: "Opened " + selected.File}
		}
	}
	return nil
}

// Update handles messages
func (vr *VetRunner) Update(msg tea.Msg) (TabInterface, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !vr.running {
			keyMap := DefaultVetKeyMap()
			switch {
			case key.Matches(msg, keyMap.Run):
				return vr, vr.runVet()
			case key.Matches(msg, keyMap.OpenFile):
				return vr, vr.openSelectedFile()
			}
		}

	case runner.CommandOutput:
		if vr.running {
			cmd := vr.handleVetOutput(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}

	case runner.CommandFinished:
		if vr.running {
			vr.running = false
			if len(vr.problems) == 0 {
				vr.status = "No issues found"
			} else {
				vr.status = lipgloss.NewStyle().Foreground(style.WarningColor).Render(
					strconv.Itoa(len(vr.problems)) + " issues found")
			}
		}
	}

	// Update list
	var cmd tea.Cmd
	vr.list, cmd = vr.list.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return vr, tea.Batch(cmds...)
}

// View renders the vet runner
func (vr VetRunner) View() string {
	var content strings.Builder

	// Header
	header := style.HeaderStyle.Render("Vet & Lint")
	content.WriteString(header + "\n")

	// Status
	statusBar := style.StatusStyle.Render(vr.status)
	content.WriteString(statusBar + "\n")

	// List
	content.WriteString(vr.list.View())

	return content.String()
}

// Refresh triggers a refresh
func (vr VetRunner) Refresh() tea.Cmd {
	return vr.runVet()
}

// GetStatus returns the current status
func (vr VetRunner) GetStatus() string {
	return vr.status
}

// GetKeyHints returns key binding hints
func (vr VetRunner) GetKeyHints() []string {
	if vr.running {
		return []string{"running..."}
	}
	return []string{"r: run vet", "o/enter: open file", "↑↓: navigate"}
}

// Cleanup performs cleanup
func (vr VetRunner) Cleanup() {
	if vr.runner != nil {
		vr.runner.Stop()
	}
}
