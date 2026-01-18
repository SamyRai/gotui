package tabs

import (
	"context"
	"goutui/internal/editor"
	"goutui/internal/runner"
	"goutui/internal/style"
	"goutui/internal/tui/components"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// VetProblem represents a problem found by vet/lint tools
type VetProblem struct {
	File    string
	Line    int
	Column  int
	Message string
	Package string
	Raw     string
	Tool    string // "govet" or "staticcheck"
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
	if vp.Tool == "staticcheck" {
		return "[SC] " + vp.Message
	}
	return "[VET] " + vp.Message
}

// Description returns the description for the list item
func (vp VetProblem) Description() string {
	return editor.FormatFileLocation(vp.File, vp.Line, vp.Column)
}

// VetRunner manages vet/lint execution
type VetRunner struct {
	ctx                context.Context
	width              int
	height             int
	runner             *runner.CommandRunner
	runnerStaticcheck  *runner.CommandRunner
	list               list.Model
	actionBar          components.ActionBar
	problems           []VetProblem
	running            bool
	currentTool        string // "govet" or "staticcheck"
	hasStaticcheck     bool
	status             string
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

	// Check if staticcheck is available first
	hasStaticcheck := false
	if _, err := exec.LookPath("staticcheck"); err == nil {
		hasStaticcheck = true
	}

	vr := &VetRunner{
		ctx:               ctx,
		runner:            runner.NewCommandRunner(ctx),
		runnerStaticcheck: runner.NewCommandRunner(ctx),
		list:              l,
		actionBar:         components.NewActionBar(),
		problems:          []VetProblem{},
		running:           false,
		currentTool:       "",
		hasStaticcheck:    hasStaticcheck,
		status:            "Ready to run vet",
	}
	vr.updateActionBar()
	return vr
}

// updateActionBar updates the action bar based on current state
func (vr *VetRunner) updateActionBar() {
	vr.actionBar.Clear()
	
	if vr.running {
		vr.actionBar.AddAction(components.Action{
			Key:         "s",
			Label:       "Stop",
			Description: "Stop vet/staticcheck",
			Primary:     true,
		})
	} else {
		label := "Run Vet"
		if vr.hasStaticcheck {
			label = "Run Vet + Lint"
		}
		vr.actionBar.AddAction(components.Action{
			Key:         "r",
			Label:       label,
			Description: "Run go vet and staticcheck",
			Primary:     true,
		})
		if len(vr.problems) > 0 {
			vr.actionBar.AddAction(components.Action{
				Key:         "o",
				Label:       "Open File",
				Description: "Open selected problem file in editor",
				Primary:     false,
			})
			vr.actionBar.AddAction(components.Action{
				Key:         "enter",
				Label:       "Open File",
				Description: "Open selected problem file in editor",
				Primary:     false,
			})
		}
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
	vr.actionBar.SetWidth(width)
	// Account for header, status, and action bar (calculate dynamically)
	actionBarHeight := vr.actionBar.Height()
	vr.list.SetWidth(width)
	vr.list.SetHeight(height - 4 - actionBarHeight)
}

// hasStaticcheckAvailable checks if staticcheck is available
func (vr *VetRunner) hasStaticcheckAvailable() bool {
	_, err := exec.LookPath("staticcheck")
	return err == nil
}

// runVet executes go vet and optionally staticcheck
func (vr *VetRunner) runVet() tea.Cmd {
	vr.running = true
	vr.currentTool = "govet"
	vr.status = "Running vet..."
	vr.problems = []VetProblem{}
	vr.updateList()
	vr.updateActionBar()

	// Start with go vet
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
		Tool:    "govet",
	}
}

// parseStaticcheckLine parses a staticcheck output line
// Example: file.go:10:5: message (SA1000)
var staticcheckLineRegex = regexp.MustCompile(`^(.+\.go):(\d+):(\d+):\s*(.+)$`)

func (vr *VetRunner) parseStaticcheckLine(line string) *VetProblem {
	matches := staticcheckLineRegex.FindStringSubmatch(line)
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
		Tool:    "staticcheck",
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

// handleStaticcheckOutput processes output from staticcheck
func (vr *VetRunner) handleStaticcheckOutput(output runner.CommandOutput) tea.Cmd {
	if problem := vr.parseStaticcheckLine(output.Line); problem != nil {
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
	vr.updateActionBar()
}

// updateStatus updates the status message based on results
func (vr *VetRunner) updateStatus() {
	if len(vr.problems) == 0 {
		tools := "vet"
		if vr.hasStaticcheck {
			tools += " and staticcheck"
		}
		vr.status = "No issues found with " + tools
	} else {
		vetCount := 0
		staticcheckCount := 0
		for _, problem := range vr.problems {
			if problem.Tool == "govet" {
				vetCount++
			} else if problem.Tool == "staticcheck" {
				staticcheckCount++
			}
		}

		parts := []string{}
		if vetCount > 0 {
			parts = append(parts, strconv.Itoa(vetCount)+" vet")
		}
		if staticcheckCount > 0 {
			parts = append(parts, strconv.Itoa(staticcheckCount)+" staticcheck")
		}

		vr.status = lipgloss.NewStyle().Foreground(style.WarningColor).Render(
			strings.Join(parts, ", ") + " issues found")
	}
}

// openSelectedFile opens the currently selected problem file
func (vr *VetRunner) openSelectedFile() tea.Cmd {
	if selected, ok := vr.list.SelectedItem().(VetProblem); ok {
		return func() tea.Msg {
			err := editor.OpenAtLine(selected.File, selected.Line, selected.Column)
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
			var cmd tea.Cmd
			if vr.currentTool == "govet" {
				cmd = vr.handleVetOutput(msg)
			} else if vr.currentTool == "staticcheck" {
				cmd = vr.handleStaticcheckOutput(msg)
			}
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}

	case runner.CommandFinished:
		if vr.running {
			if vr.currentTool == "govet" {
				// Go vet finished, start staticcheck if available
				if vr.hasStaticcheck {
					vr.currentTool = "staticcheck"
					vr.status = "Running staticcheck..."
					cmds = append(cmds, vr.runnerStaticcheck.Run("staticcheck", "./..."))
				} else {
					// No staticcheck, finish
					vr.running = false
					vr.updateStatus()
					vr.updateActionBar()
				}
			} else if vr.currentTool == "staticcheck" {
				// Staticcheck finished, we're done
				vr.running = false
				vr.updateStatus()
				vr.updateActionBar()
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
	var parts []string

	// Action bar at the top
	actionBarView := vr.actionBar.View()
	if actionBarView != "" {
		parts = append(parts, actionBarView, "")
	}

	// Header
	header := style.HeaderStyle.Render("Vet & Lint")
	parts = append(parts, header)

	// Tools info
	tools := "go vet"
	if vr.hasStaticcheck {
		tools += " + staticcheck"
	}
	toolsInfo := style.SubtleStyle.Render("Tools: " + tools)
	parts = append(parts, toolsInfo)

	// Status
	statusBar := style.StatusStyle.Render(vr.status)
	parts = append(parts, statusBar, "")

	// List
	listView := vr.list.View()
	if listView == "" && actionBarView != "" {
		parts = append(parts, style.SubtleStyle.Render("Press a key above to run vet and lint"))
	} else {
		parts = append(parts, listView)
	}

	return strings.Join(parts, "\n")
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
	hints := []string{"r: run vet"}
	if vr.hasStaticcheck {
		hints[0] = "r: run vet+lint"
	}
	hints = append(hints, "o/enter: open file", "↑↓: navigate")
	return hints
}

// Cleanup performs cleanup
func (vr VetRunner) Cleanup() {
	if vr.runner != nil {
		vr.runner.Stop()
	}
	if vr.runnerStaticcheck != nil {
		vr.runnerStaticcheck.Stop()
	}
}
