package tabs

import (
	"context"
	"fmt"
	"goutui/internal/discovery"
	"goutui/internal/editor"
	"goutui/internal/runner"
	"goutui/internal/style"
	"goutui/internal/tui/components"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// DiscoveryCompletedMsg indicates that test discovery has completed
type DiscoveryCompletedMsg struct {
	TestFiles []discovery.FileMatch
	Error     error
}


// TestRunner manages the test execution tab
type TestRunner struct {
	ctx            context.Context
	runner         *runner.CommandRunner
	parser         *runner.GoTestParser
	discoverer     *discovery.Discoverer
	treeList       components.TreeList
	logViewer      components.LogViewer
	actionBar      components.ActionBar
	width          int
	height         int
	splitView       bool
	failOnly        bool
	running         bool
	lastRun         time.Time
	availableTests  []discovery.FileMatch
	showingTests    bool // true when showing test files, false when showing results
	discoveryDone   bool // true when test discovery has completed
	status          string
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
	tr := &TestRunner{
		ctx:            ctx,
		runner:         runner.NewCommandRunner(ctx),
		parser:         runner.NewGoTestParser(),
		discoverer:     discovery.NewDiscoverer(ctx, "."),
		treeList:       components.NewTreeList(),
		logViewer:      components.NewLogViewer("Test Output"),
		actionBar:      components.NewActionBar(),
		splitView:      false,
		failOnly:       false,
		running:        false,
		availableTests: []discovery.FileMatch{},
		showingTests:   true, // Start by showing available tests
		discoveryDone:  false, // Discovery hasn't started yet
	}
	tr.updateActionBar()
	return tr
}

// Init initializes the test runner and discovers available tests
func (tr *TestRunner) Init() tea.Cmd {
	return tr.discoverTests()
}

// SetSize sets the dimensions of the test runner
func (tr *TestRunner) SetSize(width, height int) {
	tr.width = width
	tr.height = height
	tr.actionBar.SetWidth(width)

	// Calculate action bar height dynamically based on actual content
	actionBarHeight := tr.actionBar.Height()
	// Account for action bar height + 1 line spacing between action bar and content
	// The spacing is added in View() when joining actionBarView with content
	spacingHeight := 1
	if actionBarHeight == 0 {
		spacingHeight = 0 // No spacing if no action bar
	}
	contentHeight := height - actionBarHeight - spacingHeight
	if contentHeight < 10 { // Minimum usable height
		contentHeight = height
	}

	if tr.splitView {
		// Split view: tree on left, logs on right
		treeWidth := width / 2
		logWidth := width - treeWidth

		tr.treeList.SetSize(treeWidth, contentHeight)
		tr.logViewer.SetSize(logWidth, contentHeight)
	} else {
		// Single view: tree takes full space
		tr.treeList.SetSize(width, contentHeight)
		tr.logViewer.SetSize(width, contentHeight/2)
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
			tr.updateActionBar()
		case key.Matches(msg, keyMap.ToggleFail):
			if !tr.showingTests {
				tr.failOnly = !tr.failOnly
				tr.updateTreeList()
				tr.updateActionBar()
			}
		case key.Matches(msg, keyMap.OpenFile):
			return tr, tr.openSelectedFile()
		case key.Matches(msg, keyMap.Run):
			if tr.showingTests {
				selected := tr.treeList.GetSelectedNode()
				if selected != nil && selected.ID == "run_all" {
					tr.showingTests = false
					return tr, tr.runTests()
				}
			} else if !tr.running {
				return tr, tr.runTests()
			}
		}

	case runner.CommandStarted:
		tr.running = true
		tr.parser = runner.NewGoTestParser()
		tr.treeList.Clear()
		tr.logViewer.Clear()
		tr.updateActionBar()

	case runner.CommandFinished:
		tr.running = false
		tr.parser.Finish()
		tr.updateTreeList()
		tr.lastRun = time.Now()
		tr.showingTests = false // Now showing test results
		tr.updateActionBar()

	case runner.CommandOutput:
		return tr, tr.handleCommandOutput(msg)

	case DiscoveryCompletedMsg:
		return tr, tr.handleDiscoveryCompleted(msg)
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
	// Show action bar at the top
	actionBarView := tr.actionBar.View()
	
	// Determine content to show
	var content string
	hasContent := false
	
	// Check if we should show test files or results
	shouldShowTestFiles := tr.showingTests || (!tr.running && tr.parser.GetSummary() == nil)
	
	if shouldShowTestFiles {
		// Show test files discovery view
		if !tr.discoveryDone {
			// Show loading message while discovery is in progress
			discoveringMsg := lipgloss.NewStyle().
				Foreground(style.InfoColor).
				Padding(2).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(style.BorderColor).
				Width(tr.width - 4).
				Render(
					lipgloss.JoinVertical(
						lipgloss.Left,
						style.HeaderStyle.Render("Discovering Test Files..."),
						"",
						"Scanning for *_test.go files in the current directory.",
					),
				)
			content = discoveringMsg
		} else {
			// Discovery is done, show results
			treeView := tr.treeList.View()
			if treeView == "" || len(tr.availableTests) == 0 {
				// Show helpful message when no tests found
				noTestsMsg := lipgloss.NewStyle().
					Foreground(style.InfoColor).
					Padding(2).
					Border(lipgloss.RoundedBorder()).
					BorderForeground(style.BorderColor).
					Width(tr.width - 4).
					Render(
						lipgloss.JoinVertical(
							lipgloss.Left,
							style.HeaderStyle.Render("No Test Files Found"),
							"",
							"Test files are typically named *_test.go",
							"",
							"You can still run tests with the 'r' key above.",
							"This will execute 'go test ./...' in the current directory.",
						),
					)
				content = noTestsMsg
			} else {
				content = style.BorderStyle.Width(tr.width).Render(treeView)
			}
		}
		hasContent = true
	} else if tr.splitView {
		// Split view: tree on left, logs on right
		treeView := tr.treeList.View()
		logView := tr.logViewer.View()

		content = lipgloss.JoinHorizontal(
			lipgloss.Top,
			style.BorderStyle.Width(tr.width/2).Render(treeView),
			style.BorderStyle.Width(tr.width-tr.width/2).Render(logView),
		)
		hasContent = true
	} else if tr.parser.GetSummary() != nil {
		// Single view: just the tree with results
		content = style.BorderStyle.Width(tr.width).Render(tr.treeList.View())
		hasContent = true
	}

	// Combine action bar and content
	if actionBarView != "" && !hasContent {
		// Show action bar prominently when no content
		return lipgloss.JoinVertical(
			lipgloss.Left,
			"",
			actionBarView,
			"",
			style.SubtleStyle.Render("Press a key above to get started"),
		)
	} else if actionBarView != "" {
		// Show action bar above content
		return lipgloss.JoinVertical(
			lipgloss.Left,
			actionBarView,
			"",
			content,
		)
	}

	return content
}

// updateActionBar updates the action bar based on current state
func (tr *TestRunner) updateActionBar() {
	tr.actionBar.Clear()
	
	if tr.running {
		tr.actionBar.AddAction(components.Action{
			Key:         "s",
			Label:       "Stop",
			Description: "Stop running tests",
			Primary:     true,
		})
	} else if tr.showingTests || (tr.discoveryDone && len(tr.availableTests) > 0 && tr.parser.GetSummary() == nil) {
		// Showing test files - can run tests
		tr.actionBar.AddAction(components.Action{
			Key:         "r",
			Label:       "Run All Tests",
			Description: "Run all discovered tests",
			Primary:     true,
		})
		tr.actionBar.AddAction(components.Action{
			Key:         "enter",
			Label:       "Open File",
			Description: "Open selected test file in editor",
			Primary:     false,
		})
	} else if tr.parser.GetSummary() != nil {
		// Has results - can run again or interact
		tr.actionBar.AddAction(components.Action{
			Key:         "r",
			Label:       "Run Tests",
			Description: "Run tests again",
			Primary:     true,
		})
		tr.actionBar.AddAction(components.Action{
			Key:         "enter",
			Label:       "Toggle Split",
			Description: "Toggle split view with logs",
			Primary:     false,
		})
		tr.actionBar.AddAction(components.Action{
			Key:         "f",
			Label:       "Failures Only",
			Description: "Show only failing tests",
			Primary:     false,
		})
		tr.actionBar.AddAction(components.Action{
			Key:         "o",
			Label:       "Open File",
			Description: "Open selected test in editor",
			Primary:     false,
		})
	} else {
		// Initial state - no tests discovered yet
		tr.actionBar.AddAction(components.Action{
			Key:         "r",
			Label:       "Run Tests",
			Description: "Run tests in current directory",
			Primary:     true,
		})
	}
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

	if !tr.discoveryDone {
		return lipgloss.NewStyle().Foreground(style.InfoColor).Render("Discovering test files...")
	}

	if tr.showingTests || (len(tr.availableTests) > 0 && tr.parser.GetSummary() == nil) {
		return lipgloss.NewStyle().Foreground(style.InfoColor).Render(
			strconv.Itoa(len(tr.availableTests)) + " test files found")
	}

	summary := tr.parser.GetSummary()
	if summary == nil {
		return lipgloss.NewStyle().Foreground(style.SubtleColor).Render("No tests run")
	}

	counts := summary.Total
	parts := []string{}

	if counts.Pass > 0 {
		parts = append(parts, style.RenderStatusWithText(
			fmt.Sprintf("%d", counts.Pass),
			style.PassIcon,
			style.PassText,
			style.SuccessStyle))
	}
	if counts.Fail > 0 {
		parts = append(parts, style.RenderStatusWithText(
			fmt.Sprintf("%d", counts.Fail),
			style.FailIcon,
			style.FailText,
			style.ErrorStyle))
	}
	if counts.Skip > 0 {
		parts = append(parts, style.RenderStatusWithText(
			fmt.Sprintf("%d", counts.Skip),
			style.SkipIcon,
			style.SkipText,
			style.WarningStyle))
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
	if tr.showingTests || (tr.discoveryDone && tr.availableTests != nil && len(tr.availableTests) > 0 && tr.parser.GetSummary() == nil) {
		// Showing test files
		hints := []string{
			"r: run selected",
			"enter: open file",
		}
		if tr.running {
			hints = append(hints, "s: stop")
		}
		return hints
	}

	// Showing test results
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

// discoverTests finds and displays available test files using the generic discovery package
func (tr *TestRunner) discoverTests() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		// Use pattern-based discovery - go list doesn't include TestFiles by default
		// Pattern matching is more reliable for finding test files
		matcher, err := discovery.NewPatternMatcher(
			"test",
			`.*_test\.go$`, // Match files ending with _test.go
			"",              // No content pattern needed
		)
		if err != nil {
			return DiscoveryCompletedMsg{
				TestFiles: nil,
				Error:     fmt.Errorf("error creating pattern matcher: %w", err),
			}
		}

		testFiles, err := tr.discoverer.DiscoverFilesByPattern(*matcher)
		return DiscoveryCompletedMsg{
			TestFiles: testFiles,
			Error:     err,
		}
	})
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

// handleDiscoveryCompleted processes the completion of test discovery
func (tr *TestRunner) handleDiscoveryCompleted(msg DiscoveryCompletedMsg) tea.Cmd {
	if msg.Error != nil {
		tr.status = "Error discovering tests: " + msg.Error.Error()
		tr.availableTests = []discovery.FileMatch{}
	} else {
		tr.availableTests = msg.TestFiles
		if len(msg.TestFiles) == 0 {
			tr.status = "No test files found - you can still run tests with 'r' key"
		} else {
			tr.status = "Found " + strconv.Itoa(len(msg.TestFiles)) + " test files"
		}
		tr.updateTestFileList()
	}

	tr.discoveryDone = true // Mark discovery as completed
	tr.updateActionBar()
	return nil
}

// updateTestFileList updates the tree list with available test files
func (tr *TestRunner) updateTestFileList() {
	tr.treeList.Clear()

	// Add a "Run All Tests" node at the top
	runAllNode := &components.TreeNode{
		ID:     "run_all",
		Label:  "▶ Run All Tests",
		Icon:   "🚀",
		Status: "ready",
		Level:  0,
	}
	tr.treeList.AddNode(runAllNode)

	// Group test files by package using the generic discovery method
	packageMap := tr.discoverer.GroupByPackage(tr.availableTests)

	// Add packages and their test files
	for pkgName, files := range packageMap {
		pkgNode := &components.TreeNode{
			ID:     "pkg:" + pkgName,
			Label:  pkgName,
			Icon:   "📦",
			Status: "info",
			Level:  0,
		}
		tr.treeList.AddNode(pkgNode)

		for _, file := range files {
			fileName := filepath.Base(file.Path)
			testNode := &components.TreeNode{
				ID:     "file:" + file.Path,
				Label:  fileName,
				Icon:   "📄",
				Status: "info",
				Level:  1,
			}
			tr.treeList.AddNode(testNode)
		}
	}
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
			err := editor.OpenAtLine(fileName, lineNum, 1)
			if err != nil {
				return runner.CommandOutput{Line: "Error opening file: " + err.Error()}
			}
			return runner.CommandOutput{Line: "Opened " + fileName}
		}
	}

	return nil
}
