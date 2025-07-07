package tui

import (
	"context"
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
	ctx        context.Context
	cancel     context.CancelFunc
	tabBar     components.TabBar
	testTab    tabs.TabInterface
	fmtTab     tabs.TabInterface
	buildTab   tabs.TabInterface
	benchTab   tabs.TabInterface
	vetTab     tabs.TabInterface
	width      int
	height     int
	ready      bool
	focusOnTabBar bool // true = tab bar focused, false = tab content focused
	notification  runner.NotificationMsg // For displaying global messages
}

// MainKeyMap defines global key bindings
type MainKeyMap struct {
	Quit        key.Binding
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
		ctx:      ctx,
		cancel:   cancel,
		tabBar:   components.NewTabBar(),
		testTab:  tabs.NewTestRunner(ctx),
		fmtTab:   tabs.NewFmtDiff(ctx),
		buildTab: tabs.NewBuildRunner(ctx),
		benchTab: tabs.NewBenchmarkRunner(ctx),
		vetTab:   tabs.NewVetRunner(ctx),
		ready:    false,
		focusOnTabBar: true, // Start with tab bar focused
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

		contentHeight := m.height - 4 // Account for tab bar and status bar
		m.testTab.SetSize(m.width, contentHeight)
		m.fmtTab.SetSize(m.width, contentHeight)
		m.buildTab.SetSize(m.width, contentHeight)
		m.benchTab.SetSize(m.width, contentHeight)
		m.vetTab.SetSize(m.width, contentHeight)

	case tea.KeyMsg:
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

	return lipgloss.JoinVertical(
		lipgloss.Left,
		notificationView, // Display at the top
		mainContent,
		statusBar,
	)
}

// renderNotification renders the global notification bar.
func (m Model) renderNotification() string {
	var style lipgloss.Style
	switch m.notification.Type {
	case runner.SuccessNotification:
		style = lipgloss.NewStyle().Background(lipgloss.Color("#28A745"))
	case runner.ErrorNotification:
		style = lipgloss.NewStyle().Background(lipgloss.Color("#DC3545"))
	default:
		style = lipgloss.NewStyle().Background(lipgloss.Color("#007BFF"))
	}

	return style.Width(m.width).Padding(0, 1).Render(m.notification.Message)
}

// renderStatusBar renders the bottom status bar
func (m Model) renderStatusBar() string {
	activeTab := m.getCurrentTab()
	if activeTab == nil {
		return ""
	}

	// Get status from active tab
	status := activeTab.GetStatus()

	// Add key hints
	hints := []string{
		"q: quit",
		"r: refresh",
		"tab: next",
		"t/b/f/v/c: jump to tab",
	}

	// Get tab-specific hints
	tabHints := activeTab.GetKeyHints()
	hints = append(hints, tabHints...)

	hintText := lipgloss.NewStyle().
		Foreground(style.SubtleColor).
		Render(" | " + lipgloss.JoinHorizontal(lipgloss.Left, hints...))

	// Combine status and hints
	statusText := style.StatusStyle.Render(status)

	// Layout with status on left and hints on right
	statusBar := lipgloss.NewStyle().
		Width(m.width).
		Background(style.SelectedColor).
		Render(
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				statusText,
				hintText,
			),
		)

	return statusBar
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