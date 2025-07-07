package components

import (
	"goutui/internal/style"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Tab represents a single tab in the tab bar
type Tab struct {
	Name    string
	Icon    string
	Key     string
	Tooltip string
}

// TabBar manages navigation between different tabs
type TabBar struct {
	tabs      []Tab
	activeTab int
	width     int
}

// TabBarKeyMap defines key bindings for tab navigation
type TabBarKeyMap struct {
	NextTab     key.Binding
	PrevTab     key.Binding
	JumpToTests key.Binding
	JumpToBench key.Binding
	JumpToFmt   key.Binding
	JumpToVet   key.Binding
	JumpToBuild key.Binding
}

// DefaultTabBarKeyMap returns the default key bindings for tab navigation
func DefaultTabBarKeyMap() TabBarKeyMap {
	return TabBarKeyMap{
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

// NewTabBar creates a new tab bar with predefined tabs
// IMPORTANT: The order of this slice determines both the UI and tab switching order.
// Keep this in sync with the order in Model (model.go) for correct behavior.
func NewTabBar() TabBar {
	tabs := []Tab{
		{Name: "Tests", Icon: "🧪", Key: "T", Tooltip: "Run tests (go test)"},
		{Name: "Fmt", Icon: "✨", Key: "F", Tooltip: "Format check (gofmt)"},
		{Name: "Build", Icon: "🔨", Key: "C", Tooltip: "Compile (go build)"},
		{Name: "Bench", Icon: "⚙️", Key: "B", Tooltip: "Run benchmarks (go test -bench)"},
		{Name: "Vet", Icon: "🔎", Key: "V", Tooltip: "Vet & lint (go vet)"},
	}

	return TabBar{
		tabs:      tabs,
		activeTab: 0,
		width:     80,
	}
}

// SetWidth sets the width of the tab bar
func (tb *TabBar) SetWidth(width int) {
	tb.width = width
}

// ActiveTab returns the index of the currently active tab
func (tb TabBar) ActiveTab() int {
	return tb.activeTab
}

// SetActiveTab sets the active tab by index
func (tb *TabBar) SetActiveTab(index int) {
	if index >= 0 && index < len(tb.tabs) {
		tb.activeTab = index
	}
}

// NextTab moves to the next tab (wrapping around)
func (tb *TabBar) NextTab() {
	tb.activeTab = (tb.activeTab + 1) % len(tb.tabs)
}

// PrevTab moves to the previous tab (wrapping around)
func (tb *TabBar) PrevTab() {
	tb.activeTab = (tb.activeTab - 1 + len(tb.tabs)) % len(tb.tabs)
}

// JumpToTab jumps to a specific tab by key
func (tb *TabBar) JumpToTab(tabKey string) {
	for i, tab := range tb.tabs {
		if tab.Key == tabKey {
			tb.activeTab = i
			break
		}
	}
}

// Update handles tab bar messages
func (tb TabBar) Update(msg tea.Msg) (TabBar, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		keyMap := DefaultTabBarKeyMap()
		switch {
		case key.Matches(msg, keyMap.NextTab):
			tb.NextTab()
		case key.Matches(msg, keyMap.PrevTab):
			tb.PrevTab()
		case key.Matches(msg, keyMap.JumpToTests):
			tb.JumpToTab("T")
		case key.Matches(msg, keyMap.JumpToBench):
			tb.JumpToTab("B")
		case key.Matches(msg, keyMap.JumpToFmt):
			tb.JumpToTab("F")
		case key.Matches(msg, keyMap.JumpToVet):
			tb.JumpToTab("V")
		case key.Matches(msg, keyMap.JumpToBuild):
			tb.JumpToTab("C")
		}
	}
	return tb, nil
}

// View renders the tab bar
// If focus is true, the tab bar is visually highlighted
func (tb TabBar) View(focus bool) string {
	var tabs []string
	for i, tab := range tb.tabs {
		var tabStyle lipgloss.Style
		if i == tb.activeTab {
			tabStyle = style.ActiveTabStyle
		} else {
			tabStyle = style.InactiveTabStyle
		}
		tabContent := lipgloss.JoinHorizontal(
			lipgloss.Left,
			tab.Icon,
			" ",
			tab.Name,
		)
		tabs = append(tabs, tabStyle.Render(tabContent))
	}

	tabRow := lipgloss.JoinHorizontal(lipgloss.Left, tabs...)

	// Add padding to center the tabs
	availableWidth := tb.width
	tabRowWidth := lipgloss.Width(tabRow)
	if tabRowWidth < availableWidth {
		padding := (availableWidth - tabRowWidth) / 2
		tabRow = lipgloss.NewStyle().
			PaddingLeft(padding).
			Render(tabRow)
	}

	// Add a border at the bottom, highlight if focused
	var border lipgloss.Style
	if focus {
		border = lipgloss.NewStyle().
			Foreground(style.PrimaryColor).
			Bold(true)
	} else {
		border = lipgloss.NewStyle().
			Foreground(style.BorderColor)
	}

	borderLine := border.Render(lipgloss.NewStyle().Width(tb.width).Render("─"))

	return lipgloss.JoinVertical(lipgloss.Left, tabRow, borderLine)
}

// GetActiveTabName returns the name of the currently active tab
func (tb TabBar) GetActiveTabName() string {
	if tb.activeTab >= 0 && tb.activeTab < len(tb.tabs) {
		return tb.tabs[tb.activeTab].Name
	}
	return ""
}
