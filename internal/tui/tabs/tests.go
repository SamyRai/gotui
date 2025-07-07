package tabs

import (
	"context"
	"goutui/internal/runner"
	"goutui/internal/style"
	"goutui/internal/tui/components"
	"goutui/internal/util"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TestRunner manages the test execution tab
type TestRunner struct {
	ctx        context.Context
	runner     *runner.CommandRunner
	parser     *runner.GoTestParser
	treeList   components.TreeList
	logViewer  components.LogViewer
	width      int
	height     int
	splitView  bool
	failOnly   bool
	running    bool
	lastRun    time.Time
}

// TestRunnerKeyMap defines key bindings for the test runner
type TestRunnerKeyMap struct {
	Run         key.Binding
	Stop        key.Binding
	ToggleSplit key.Binding
	ToggleFail  key.Binding
	OpenFile    key.Binding
	Filter      key.Binding
}

// DefaultTestRunnerKeyMap returns default key bindings
func DefaultTestRunnerKeyMap() TestRunnerKeyMap {
	return TestRunnerKeyMap{
		Run: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "run tests"),
		),
		Stop: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "stop"),
		),
		ToggleSplit: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "toggle split"),
		),
		ToggleFail: key.NewBinding(
			key.WithKeys("f"),
			key.WithHelp("f", "show failures only"),
		),
		OpenFile: key.NewBinding(
			key.WithKeys("o"),
			key.WithHelp("o", "open in editor"),
		),
		Filter: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "filter"),
		),
	}
}

// NewTestRunner creates a new test runner
func NewTestRunner(ctx context.Context) *TestRunner {
	return &TestRunner{
		ctx:       ctx,
		runner:    runner.NewCommandRunner(ctx),
		parser:    runner.NewGoTestParser(),
		treeList:  components.NewTreeList(),
		logViewer: components.NewLogViewer("Test Output"),
		splitView: false,
		failOnly:  false,
		running:   false,
	}
}

// Init initializes the test runner
func (tr TestRunner) Init() tea.Cmd {
	return nil
}

// SetSize sets the dimensions of the test runner
func (tr *TestRunner) SetSize(width, height int) {
	tr.width = width
	tr.height = height

	if tr.splitView {
		// Split view: tree on left, logs on right
		treeWidth := width / 2
		logWidth := width - treeWidth

		tr.treeList.SetSize(treeWidth, height)
		tr.logViewer.SetSize(logWidth, height)
	} else {
		// Single view: tree takes full space
		tr.treeList.SetSize(width, height)
		tr.logViewer.SetSize(width, height/2)
	}
}

// Update handles messages
func (tr *TestRunner) Update(msg tea.Msg) (TabInterface, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		keyMap := DefaultTestRunnerKeyMap()
		switch {
		case key.Matches(msg, keyMap.Run):
			if !tr.running {
				return tr, tr.runTests()
			}
		case key.Matches(msg, keyMap.Stop):
			if tr.running {
				tr.runner.Stop()
				tr.running = false
			}
		case key.Matches(msg, keyMap.ToggleSplit):
			tr.splitView = !tr.splitView
			tr.SetSize(tr.width, tr.height)
		case key.Matches(msg, keyMap.ToggleFail):
			tr.failOnly = !tr.failOnly
			tr.updateTreeList()
		case key.Matches(msg, keyMap.OpenFile):
			return tr, tr.openSelectedFile()
		}

	case runner.CommandStarted:
		tr.running = true
		tr.parser = runner.NewGoTestParser()
		tr.treeList.Clear()
		tr.logViewer.Clear()

	case runner.CommandFinished:
		tr.running = false
		tr.parser.Finish()
		tr.updateTreeList()
		tr.lastRun = time.Now()

	case runner.CommandOutput:
		return tr, tr.handleCommandOutput(msg)
	}

	// Update child components
	var cmd tea.Cmd
	tr.treeList, cmd = tr.treeList.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	tr.logViewer, cmd = tr.logViewer.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return tr, tea.Batch(cmds...)
}

// View renders the test runner
func (tr TestRunner) View() string {
	if tr.splitView {
		// Split view: tree on left, logs on right
		treeView := tr.treeList.View()
		logView := tr.logViewer.View()

		return lipgloss.JoinHorizontal(
			lipgloss.Top,
			style.BorderStyle.Width(tr.width/2).Render(treeView),
			style.BorderStyle.Width(tr.width-tr.width/2).Render(logView),
		)
	}

	// Single view: just the tree
	return style.BorderStyle.Width(tr.width).Render(tr.treeList.View())
}

// Refresh triggers a test run
func (tr TestRunner) Refresh() tea.Cmd {
	if !tr.running {
		return tr.runTests()
	}
	return nil
}

// GetStatus returns the current status
func (tr TestRunner) GetStatus() string {
	if tr.running {
		return lipgloss.NewStyle().Foreground(style.InfoColor).Render("Running tests...")
	}

	summary := tr.parser.GetSummary()
	if summary == nil {
		return lipgloss.NewStyle().Foreground(style.SubtleColor).Render("No tests run")
	}

	counts := summary.Total
	parts := []string{}

	if counts.Pass > 0 {
		parts = append(parts, style.SuccessStyle.Render("✔ "+string(rune(counts.Pass+'0'))))
	}
	if counts.Fail > 0 {
		parts = append(parts, style.ErrorStyle.Render("✖ "+string(rune(counts.Fail+'0'))))
	}
	if counts.Skip > 0 {
		parts = append(parts, style.WarningStyle.Render("⏸ "+string(rune(counts.Skip+'0'))))
	}

	status := strings.Join(parts, " ")

	if !tr.lastRun.IsZero() {
		duration := summary.Duration.Truncate(time.Millisecond)
		status += lipgloss.NewStyle().Foreground(style.SubtleColor).Render(" ("+duration.String()+")")
	}

	return status
}

// GetKeyHints returns key binding hints
func (tr TestRunner) GetKeyHints() []string {
	hints := []string{
		"r: run",
		"enter: split view",
		"f: fail only",
	}

	if tr.running {
		hints = append(hints, "s: stop")
	}

	return hints
}

// Cleanup performs cleanup
func (tr TestRunner) Cleanup() {
	if tr.runner != nil {
		tr.runner.Stop()
	}
}

// runTests executes go test -json ./...
func (tr *TestRunner) runTests() tea.Cmd {
	return tr.runner.Run("go", "test", "-json", "-v", "./...")
}

// handleCommandOutput processes a line of output from go test -json
func (tr *TestRunner) handleCommandOutput(msg runner.CommandOutput) tea.Cmd {
	// Parse the event
	event, err := tr.parser.ParseLine(msg.Line)
	if err != nil || event == nil {
		return nil
	}

	// Update tree and logs
	tr.updateTreeList()
	tr.updateLogViewer(event)
	return nil
}

// updateTreeList updates the tree list from the current test summary
func (tr *TestRunner) updateTreeList() {
	summary := tr.parser.GetSummary()
	tr.treeList.Clear()
	for pkgName, pkg := range summary.Packages {
		pkgNode := &components.TreeNode{
			ID:     pkgName,
			Label:  pkgName,
			Icon:   "📦",
			Status: pkg.Status,
			Level:  0,
		}
		for testName, test := range pkg.Tests {
			icon := ""
			switch test.Status {
			case "pass":
				icon = style.PassIcon
			case "fail":
				icon = style.FailIcon
			case "skip":
				icon = style.SkipIcon
			default:
				icon = style.PendingIcon
			}
			testNode := &components.TreeNode{
				ID:     pkgName + "/" + testName,
				Label:  testName,
				Icon:   icon,
				Status: test.Status,
				Level:  1,
				Data:   test,
			}
			pkgNode.Children = append(pkgNode.Children, testNode)
		}
		tr.treeList.AddNode(pkgNode)
	}
}

// updateLogViewer updates the log viewer for the selected test
func (tr *TestRunner) updateLogViewer(event *runner.GoTestEvent) {
	selected := tr.treeList.GetSelectedNode()
	if selected == nil || selected.Level != 1 {
		tr.logViewer.SetContent("")
		return
	}
	test, ok := selected.Data.(*runner.TestCase)
	if !ok {
		tr.logViewer.SetContent("")
		return
	}
	tr.logViewer.SetContent(strings.Join(test.Output, "\n"))
}

// openSelectedFile opens the currently selected test file in editor
func (tr *TestRunner) openSelectedFile() tea.Cmd {
	selected := tr.treeList.GetSelectedNode()
	if selected == nil {
		return nil
	}

	var fileName string
	var lineNum int = 1

	if selected.Level == 1 { // Test case
		if test, ok := selected.Data.(*runner.TestCase); ok {
			// Try to extract file and line from test output
			for _, line := range test.Output {
				if strings.Contains(line, ".go:") {
					parts := strings.Fields(line)
					for _, part := range parts {
						if strings.Contains(part, ".go:") {
							fileParts := strings.Split(part, ":")
							if len(fileParts) >= 2 {
								fileName = fileParts[0]
								if num, err := strconv.Atoi(fileParts[1]); err == nil {
									lineNum = num
								}
								break
							}
						}
					}
					if fileName != "" {
						break
					}
				}
			}
		}
	} else { // Package
		// Open the package directory
		fileName = selected.ID
	}

	if fileName != "" {
		return func() tea.Msg {
			err := util.OpenInEditor(fileName, lineNum, 1)
			if err != nil {
				return runner.CommandOutput{Line: "Error opening file: " + err.Error()}
			}
			return runner.CommandOutput{Line: "Opened " + fileName}
		}
	}

	return nil
}
