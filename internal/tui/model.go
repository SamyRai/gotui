package tui

import (
	"context"
	"strings"
	"goutui/internal/runner"
	"goutui/internal/style"
	"goutui/internal/tui/components"
	"goutui/internal/tui/tabs"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model represents the main TUI model
// focusOnTabBar: true if the tab bar is focused, false if the tab content is focused
// This enables clear navigation and visual feedback.
type Model struct {
	ctx           context.Context
	cancel        context.CancelFunc
	tabBar        components.TabBar
	testTab       tabs.TabInterface
	fmtTab        tabs.TabInterface
	buildTab      tabs.TabInterface
	benchTab      tabs.TabInterface
	vetTab        tabs.TabInterface
	width         int
	height        int
	ready         bool
	focusOnTabBar bool // true = tab bar focused, false = tab content focused
	showHelp      bool // true if help overlay is shown
	notification  runner.NotificationMsg // For displaying global messages
}

// MainKeyMap defines global key bindings
type MainKeyMap struct {
	Quit        key.Binding
	Help        key.Binding
	Refresh     key.Binding
	NextTab     key.Binding
	PrevTab     key.Binding
	JumpToTests key.Binding
	JumpToBench key.Binding
	JumpToFmt   key.Binding
	JumpToVet   key.Binding
	JumpToBuild key.Binding
}

// DefaultMainKeyMap returns the default key bindings
func DefaultMainKeyMap() MainKeyMap {
	return MainKeyMap{
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q/ctrl+c", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh"),
		),
		NextTab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next tab"),
		),
		PrevTab: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "previous tab"),
		),
		JumpToTests: key.NewBinding(
			key.WithKeys("t"),
			key.WithHelp("t", "tests"),
		),
		JumpToBench: key.NewBinding(
			key.WithKeys("b"),
			key.WithHelp("b", "benchmarks"),
		),
		JumpToFmt: key.NewBinding(
			key.WithKeys("f"),
			key.WithHelp("f", "format"),
		),
		JumpToVet: key.NewBinding(
			key.WithKeys("v"),
			key.WithHelp("v", "vet"),
		),
		JumpToBuild: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "compile"),
		),
	}
}

// NewModel creates a new main model
// IMPORTANT: The order of the tab fields below must match the order in TabBar (tabbar.go)
// for correct tab switching and rendering.
func NewModel() Model {
	ctx, cancel := context.WithCancel(context.Background())
	return Model{
		ctx:           ctx,
		cancel:        cancel,
		tabBar:        components.NewTabBar(),
		testTab:       tabs.NewTestRunner(ctx),
		fmtTab:        tabs.NewFmtDiff(ctx),
		buildTab:      tabs.NewBuildRunner(ctx),
		benchTab:      tabs.NewBenchmarkRunner(ctx),
		vetTab:        tabs.NewVetRunner(ctx),
		ready:         false,
		focusOnTabBar: true, // Start with tab bar focused
		showHelp:      false,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.WindowSize(),
		m.testTab.Init(),
		m.benchTab.Init(),
		m.fmtTab.Init(),
		m.vetTab.Init(),
		m.buildTab.Init(),
	)
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true

		// Set dimensions for all components
		m.tabBar.SetWidth(m.width)

		// Account for tab bar (2 lines) and status bar (1 line) = 3 lines total
		// Also account for potential spacing between elements
		contentHeight := m.height - 3
		if contentHeight < 10 {
			contentHeight = m.height - 1 // Minimum space, keep at least status bar
		}
		m.testTab.SetSize(m.width, contentHeight)
		m.fmtTab.SetSize(m.width, contentHeight)
		m.buildTab.SetSize(m.width, contentHeight)
		m.benchTab.SetSize(m.width, contentHeight)
		m.vetTab.SetSize(m.width, contentHeight)

	case tea.KeyMsg:
		// Handle help overlay first (global)
		if m.showHelp {
			if msg.String() == "esc" || msg.String() == "?" {
				m.showHelp = false
				return m, nil
			}
		}

		keyMap := DefaultMainKeyMap()
		if m.focusOnTabBar {
			switch msg.String() {
			case "left":
				m.tabBar.PrevTab()
			case "right":
				m.tabBar.NextTab()
			case "enter", "down":
				m.focusOnTabBar = false // Move focus to tab content
			default:
				switch {
				case key.Matches(msg, keyMap.NextTab):
					m.tabBar.NextTab()
				case key.Matches(msg, keyMap.PrevTab):
					m.tabBar.PrevTab()
				case key.Matches(msg, keyMap.JumpToTests):
					m.tabBar.JumpToTab("T")
				case key.Matches(msg, keyMap.JumpToBench):
					m.tabBar.JumpToTab("B")
				case key.Matches(msg, keyMap.JumpToFmt):
					m.tabBar.JumpToTab("F")
				case key.Matches(msg, keyMap.JumpToVet):
					m.tabBar.JumpToTab("V")
				case key.Matches(msg, keyMap.JumpToBuild):
					m.tabBar.JumpToTab("C")
				case key.Matches(msg, keyMap.Help):
					m.showHelp = !m.showHelp
				case key.Matches(msg, keyMap.Quit):
					m.cleanup()
					return m, tea.Quit
				}
			}
		} else {
			// Tab content focused
			switch msg.String() {
			case "esc", "up":
				m.focusOnTabBar = true // Move focus to tab bar
			default:
				switch {
				case key.Matches(msg, keyMap.Help):
					m.showHelp = !m.showHelp
				case key.Matches(msg, keyMap.Quit):
					m.cleanup()
					return m, tea.Quit
				case key.Matches(msg, keyMap.Refresh):
					return m, m.getCurrentTab().Refresh()
				}
			}
		}

	case runner.CommandStarted:
		// Handle command started events

	case runner.CommandFinished:
		// Handle command finished events

	case runner.CommandOutput:
		// Handle command output events

	case runner.NotificationMsg:
		m.notification = msg
		return m, nil // Or trigger a timed clear
	}

	// Update tab bar
	var cmd tea.Cmd
	m.tabBar, cmd = m.tabBar.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	// Update the active tab
	activeTab := m.getCurrentTab()
	if activeTab != nil {
		var cmd tea.Cmd
		updated, cmd := activeTab.Update(msg)
		switch m.tabBar.ActiveTab() {
		case 0:
			m.testTab = updated
		case 1:
			m.fmtTab = updated
		case 2:
			m.buildTab = updated
		case 3:
			m.benchTab = updated
		case 4:
			m.vetTab = updated
		}
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

// View renders the main view
func (m Model) View() string {
	if !m.ready {
		return style.StatusStyle.Render("Initializing...")
	}

	// Render tab bar with focus indication
	tabBarView := m.tabBar.View(m.focusOnTabBar)

	// Render active tab content
	var content string
	switch m.tabBar.ActiveTab() {
	case 0:
		content = m.testTab.View()
	case 1:
		content = m.fmtTab.View()
	case 2:
		content = m.buildTab.View()
	case 3:
		content = m.benchTab.View()
	case 4:
		content = m.vetTab.View()
	default:
		content = style.ErrorStyle.Render("Unknown tab")
	}

	// Render status bar
	statusBar := m.renderStatusBar()

	// Render notification if it exists
	var notificationView string
	if m.notification.Message != "" {
		notificationView = m.renderNotification()
	}

	// Combine all parts
	mainContent := lipgloss.JoinVertical(
		lipgloss.Left,
		tabBarView,
		content,
	)

	// Render help overlay if requested
	var helpView string
	if m.showHelp {
		helpView = m.renderHelp()
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		notificationView, // Display at the top
		helpView,         // Help overlay (when shown)
		mainContent,
		statusBar,
	)
}

// renderNotification renders the global notification bar with improved accessibility.
func (m Model) renderNotification() string {
	var notificationStyle lipgloss.Style
	var prefix string
	
	switch m.notification.Type {
	case runner.SuccessNotification:
		notificationStyle = style.NotificationSuccessStyle
		prefix = style.GetStatusIcon(style.PassIcon, style.PassText) + " "
	case runner.ErrorNotification:
		notificationStyle = style.NotificationErrorStyle
		prefix = style.GetStatusIcon(style.FailIcon, style.FailText) + " "
	default:
		notificationStyle = style.NotificationInfoStyle
		prefix = "ℹ "
	}

	// Add text prefix for accessibility (visible even without color)
	message := prefix + m.notification.Message
	return notificationStyle.Width(m.width).Render(message)
}

// renderStatusBar renders the bottom status bar with improved UX
func (m Model) renderStatusBar() string {
	activeTab := m.getCurrentTab()
	if activeTab == nil {
		return ""
	}

	// Get status from active tab
	status := activeTab.GetStatus()

	// Add focus indicator
	focusIndicator := ""
	if m.focusOnTabBar {
		focusIndicator = " [TAB BAR] "
	} else {
		focusIndicator = " [CONTENT] "
	}

	// Add key hints - more organized and readable
	hints := []string{
		"q: quit",
		"r: refresh",
		"?: help",
	}

	// Add navigation hints based on focus
	if m.focusOnTabBar {
		hints = append(hints, "←→: navigate", "↓/Enter: focus")
	} else {
		hints = append(hints, "↑/Esc: focus tabs")
	}

	// Get tab-specific hints
	tabHints := activeTab.GetKeyHints()
	if len(tabHints) > 0 {
		hints = append(hints, tabHints...)
	}

	// Format hints with proper spacing between each hint
	// Join hints with " • " separator for better readability
	hintParts := make([]string, 0, len(hints))
	for _, hint := range hints {
		hintParts = append(hintParts, hint)
	}
	hintText := lipgloss.NewStyle().
		Foreground(style.SubtleColor).
		Render(" • " + strings.Join(hintParts, " • "))

	// Combine status, focus indicator, and hints
	statusText := style.StatusStyle.Render(status + focusIndicator)

	// Calculate available width for hints
	statusWidth := lipgloss.Width(statusText)
	padding := 2 // Left and right padding
	availableWidth := m.width - statusWidth - padding
	
	// Truncate hints if they don't fit
	hintWidth := lipgloss.Width(hintText)
	if hintWidth > availableWidth {
		// Truncate hint text to fit
		truncated := hintText
		for lipgloss.Width(truncated) > availableWidth && len(hintParts) > 1 {
			// Remove last hint
			hintParts = hintParts[:len(hintParts)-1]
			truncated = lipgloss.NewStyle().
				Foreground(style.SubtleColor).
				Render(" • " + strings.Join(hintParts, " • "))
		}
		hintText = truncated
	}

	// Layout with status on left and hints on right
	statusBar := lipgloss.NewStyle().
		Width(m.width).
		Background(style.SelectedColor).
		Padding(0, 1).
		Render(
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				statusText,
				lipgloss.NewStyle().
					Width(m.width-statusWidth-padding).
					Align(lipgloss.Right).
					Render(hintText),
			),
		)

	return statusBar
}

// renderHelp renders the help overlay with improved design
func (m Model) renderHelp() string {
	helpContent := []string{
		style.HeaderStyle.Render("GoTUI - Terminal UI for Go Development"),
		"",
		style.HeaderStyle.Copy().Foreground(style.AccentColor).Render("Global Shortcuts:"),
		"  Tab/Shift+Tab    Switch between tabs",
		"  ←/→              Navigate tabs (when tab bar focused)",
		"  t/b/f/v/c        Jump to Tests/Benchmarks/Format/Vet/Build",
		"  ?                Toggle this help",
		"  r                Refresh current tab",
		"  q/Ctrl+C         Quit",
		"",
		style.HeaderStyle.Copy().Foreground(style.AccentColor).Render("Navigation:"),
		"  ↑/↓/j/k          Navigate lists",
		"  Enter            Open file in editor / Toggle split view",
		"  Space            Expand/collapse items",
		"  Esc              Return to tab navigation",
		"  ↓/Enter          Focus content (from tab bar)",
		"  ↑/Esc            Focus tab bar (from content)",
		"",
		style.HeaderStyle.Copy().Foreground(style.InfoColor).Render("Current Tab: " + m.getCurrentTabName()),
		"",
		style.SubtleStyle.Render("Press ? or Esc to close this help"),
	}

	// Create a bordered box for the help content with better styling
	boxWidth := 85
	if m.width-8 < 85 {
		boxWidth = m.width - 8
	}
	if boxWidth < 50 {
		boxWidth = 50
	}
	boxHeight := len(helpContent) + 4

	helpText := lipgloss.JoinVertical(lipgloss.Left, helpContent...)

	// Create a semi-transparent overlay effect
	box := lipgloss.NewStyle().
		Width(boxWidth).
		Height(boxHeight).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(style.PrimaryColor).
		Padding(1, 2).
		Background(style.BackgroundColor).
		Render(helpText)

	// Center the box horizontally and vertically
	verticalPadding := (m.height - boxHeight) / 2
	if verticalPadding < 0 {
		verticalPadding = 0
	}
	
	return lipgloss.JoinVertical(
		lipgloss.Center,
		strings.Repeat("\n", verticalPadding),
		lipgloss.NewStyle().
			Width(m.width).
			Align(lipgloss.Center).
			Render(box),
	)
}

// getCurrentTabName returns the name of the currently active tab
func (m Model) getCurrentTabName() string {
	switch m.tabBar.ActiveTab() {
	case 0:
		return "Tests"
	case 1:
		return "Format"
	case 2:
		return "Build"
	case 3:
		return "Benchmarks"
	case 4:
		return "Vet & Lint"
	default:
		return "Unknown"
	}
}

// getCurrentTab returns the currently active tab
func (m Model) getCurrentTab() tabs.TabInterface {
	switch m.tabBar.ActiveTab() {
	case 0:
		return m.testTab
	case 1:
		return m.fmtTab
	case 2:
		return m.buildTab
	case 3:
		return m.benchTab
	case 4:
		return m.vetTab
	default:
		return nil
	}
}

// cleanup performs cleanup when exiting
func (m Model) cleanup() {
	if m.cancel != nil {
		m.cancel()
	}

	// Cleanup individual tabs
	m.testTab.Cleanup()
	m.benchTab.Cleanup()
	m.fmtTab.Cleanup()
	m.vetTab.Cleanup()
	m.buildTab.Cleanup()
}