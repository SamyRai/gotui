package components

import (
	"fmt"
	"goutui/internal/style"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// LogEntry represents a single log entry
type LogEntry struct {
	Level     string
	Message   string
	Timestamp string
	Source    string
	Raw       string
}

// LogViewer manages scrollable log display with filtering
type LogViewer struct {
	viewport    viewport.Model
	logs        []LogEntry
	filter      string
	showFilter  bool
	title       string
	borderStyle lipgloss.Style
}

// LogViewerKeyMap defines key bindings for log viewer
type LogViewerKeyMap struct {
	Up       key.Binding
	Down     key.Binding
	PageUp   key.Binding
	PageDown key.Binding
	Home     key.Binding
	End      key.Binding
	Filter   key.Binding
	Copy     key.Binding
}

// DefaultLogViewerKeyMap returns default key bindings
func DefaultLogViewerKeyMap() LogViewerKeyMap {
	return LogViewerKeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup"),
			key.WithHelp("pgup", "page up"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("pgdown"),
			key.WithHelp("pgdown", "page down"),
		),
		Home: key.NewBinding(
			key.WithKeys("home", "g"),
			key.WithHelp("home/g", "top"),
		),
		End: key.NewBinding(
			key.WithKeys("end", "G"),
			key.WithHelp("end/G", "bottom"),
		),
		Filter: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "filter"),
		),
		Copy: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "copy"),
		),
	}
}

// NewLogViewer creates a new log viewer
func NewLogViewer(title string) LogViewer {
	vp := viewport.New(40, 10)
	vp.Style = style.LogStyle
	
	return LogViewer{
		viewport:    vp,
		logs:        make([]LogEntry, 0),
		title:       title,
		borderStyle: style.BorderStyle,
	}
}

// SetSize sets the dimensions of the log viewer
func (lv *LogViewer) SetSize(width, height int) {
	lv.viewport.Width = width - 2  // Account for border
	lv.viewport.Height = height - 3 // Account for border and title
}

// SetTitle sets the title of the log viewer
func (lv *LogViewer) SetTitle(title string) {
	lv.title = title
}

// AddLog adds a new log entry
func (lv *LogViewer) AddLog(entry LogEntry) {
	lv.logs = append(lv.logs, entry)
	lv.updateContent()
}

// AddRawLog adds a raw log message
func (lv *LogViewer) AddRawLog(message string) {
	entry := LogEntry{
		Level:   "info",
		Message: message,
		Raw:     message,
	}
	lv.AddLog(entry)
}

// Clear clears all logs
func (lv *LogViewer) Clear() {
	lv.logs = make([]LogEntry, 0)
	lv.updateContent()
}

// SetFilter sets a filter for log entries
func (lv *LogViewer) SetFilter(filter string) {
	lv.filter = filter
	lv.showFilter = filter != ""
	lv.updateContent()
}

// ClearFilter clears the current filter
func (lv *LogViewer) ClearFilter() {
	lv.filter = ""
	lv.showFilter = false
	lv.updateContent()
}

// SetContent replaces the log viewer content with a single log entry (for direct display)
func (lv *LogViewer) SetContent(content string) {
	lv.logs = []LogEntry{{
		Level:   "info",
		Message: content,
		Raw:     content,
	}}
	lv.updateContent()
}

// updateContent updates the viewport content based on current logs and filter
func (lv *LogViewer) updateContent() {
	var content []string
	
	for _, log := range lv.logs {
		if lv.matchesFilter(log) {
			formatted := lv.formatLogEntry(log)
			content = append(content, formatted)
		}
	}
	
	lv.viewport.SetContent(strings.Join(content, "\n"))
	
	// Auto-scroll to bottom for new content
	lv.viewport.GotoBottom()
}

// matchesFilter checks if a log entry matches the current filter
func (lv *LogViewer) matchesFilter(entry LogEntry) bool {
	if lv.filter == "" {
		return true
	}
	
	filter := strings.ToLower(lv.filter)
	return strings.Contains(strings.ToLower(entry.Message), filter) ||
		strings.Contains(strings.ToLower(entry.Level), filter) ||
		strings.Contains(strings.ToLower(entry.Source), filter)
}

// formatLogEntry formats a log entry for display with improved accessibility
func (lv *LogViewer) formatLogEntry(entry LogEntry) string {
	var levelStyle lipgloss.Style
	var levelIcon string
	var levelText string
	
	switch strings.ToLower(entry.Level) {
	case "error", "fail":
		levelStyle = style.ErrorStyle
		levelIcon = style.GetStatusIcon(style.FailIcon, style.FailText)
		levelText = "ERROR"
	case "warn", "warning":
		levelStyle = style.WarningStyle
		levelIcon = "⚠"
		levelText = "WARN"
	case "info":
		levelStyle = lipgloss.NewStyle().Foreground(style.InfoColor).Bold(true)
		levelIcon = "ℹ"
		levelText = "INFO"
	case "debug":
		levelStyle = lipgloss.NewStyle().Foreground(style.SubtleColor)
		levelIcon = "🔍"
		levelText = "DEBUG"
	case "pass", "ok":
		levelStyle = style.SuccessStyle
		levelIcon = style.GetStatusIcon(style.PassIcon, style.PassText)
		levelText = "PASS"
	default:
		levelStyle = lipgloss.NewStyle().Foreground(style.TextColor)
		levelIcon = ""
		levelText = strings.ToUpper(entry.Level)
	}
	
	// If we have structured data, format it nicely with icons
	if entry.Source != "" || entry.Timestamp != "" {
		parts := []string{}
		
		if entry.Timestamp != "" {
			parts = append(parts, 
				lipgloss.NewStyle().Foreground(style.SubtleColor).Render(entry.Timestamp))
		}
		
		if entry.Level != "" {
			levelDisplay := levelIcon
			if levelDisplay == "" {
				levelDisplay = levelText
			} else {
				levelDisplay = levelIcon + " " + levelText
			}
			parts = append(parts, levelStyle.Render(levelDisplay))
		}
		
		if entry.Source != "" {
			parts = append(parts, 
				lipgloss.NewStyle().Foreground(style.SubtleColor).Render("["+entry.Source+"]"))
		}
		
		prefix := strings.Join(parts, " ")
		return lipgloss.JoinHorizontal(lipgloss.Left, prefix, " ", entry.Message)
	}
	
	// For raw messages, apply styling with level indicator
	if entry.Level != "" {
		prefix := levelIcon
		if prefix != "" {
			prefix = prefix + " "
		}
		return levelStyle.Render(prefix + entry.Message)
	}
	
	return entry.Message
}

// Update handles log viewer messages
func (lv LogViewer) Update(msg tea.Msg) (LogViewer, tea.Cmd) {
	var cmd tea.Cmd
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		keyMap := DefaultLogViewerKeyMap()
		switch {
		case key.Matches(msg, keyMap.Up):
			lv.viewport.LineUp(1)
		case key.Matches(msg, keyMap.Down):
			lv.viewport.LineDown(1)
		case key.Matches(msg, keyMap.PageUp):
			lv.viewport.ViewUp()
		case key.Matches(msg, keyMap.PageDown):
			lv.viewport.ViewDown()
		case key.Matches(msg, keyMap.Home):
			lv.viewport.GotoTop()
		case key.Matches(msg, keyMap.End):
			lv.viewport.GotoBottom()
		}
	}
	
	lv.viewport, cmd = lv.viewport.Update(msg)
	return lv, cmd
}

// View renders the log viewer
func (lv LogViewer) View() string {
	// Build title with filter indicator
	title := lv.title
	if lv.showFilter {
		title += " (filtered: " + lv.filter + ")"
	}
	
	// Add log count with proper number formatting
	visibleCount := lv.getVisibleLogCount()
	totalCount := len(lv.logs)
	
	if lv.showFilter {
		title += lipgloss.NewStyle().
			Foreground(style.SubtleColor).
			Render(fmt.Sprintf(" [%d/%d]", visibleCount, totalCount))
	} else {
		title += lipgloss.NewStyle().
			Foreground(style.SubtleColor).
			Render(fmt.Sprintf(" [%d]", totalCount))
	}
	
	titleBar := style.HeaderStyle.Render(title)
	
	// Render viewport with border
	content := lv.borderStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			titleBar,
			lv.viewport.View(),
		),
	)
	
	return content
}

// getVisibleLogCount returns the number of logs matching the current filter
func (lv *LogViewer) getVisibleLogCount() int {
	if !lv.showFilter {
		return len(lv.logs)
	}
	
	count := 0
	for _, log := range lv.logs {
		if lv.matchesFilter(log) {
			count++
		}
	}
	return count
}

// GetLogs returns all log entries
func (lv *LogViewer) GetLogs() []LogEntry {
	return lv.logs
}

// ScrollToBottom scrolls to the bottom of the log
func (lv *LogViewer) ScrollToBottom() {
	lv.viewport.GotoBottom()
}
