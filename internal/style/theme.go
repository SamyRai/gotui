package style

import "github.com/charmbracelet/lipgloss"

// Color scheme for the TUI
var (
	// Primary colors
	PrimaryColor   = lipgloss.Color("#00D4AA")
	SecondaryColor = lipgloss.Color("#7C3AED")
	AccentColor    = lipgloss.Color("#F59E0B")
	
	// Status colors
	SuccessColor = lipgloss.Color("#10B981")
	ErrorColor   = lipgloss.Color("#EF4444")
	WarningColor = lipgloss.Color("#F59E0B")
	InfoColor    = lipgloss.Color("#3B82F6")
	
	// UI colors
	BorderColor     = lipgloss.Color("#374151")
	SelectedColor   = lipgloss.Color("#1F2937")
	BackgroundColor = lipgloss.Color("#111827")
	TextColor       = lipgloss.Color("#F9FAFB")
	SubtleColor     = lipgloss.Color("#9CA3AF")
	
	// Tab colors
	ActiveTabColor   = PrimaryColor
	InactiveTabColor = SubtleColor
)

// Common styles
var (
	// Tab styles
	TabStyle = lipgloss.NewStyle().
		Padding(0, 2).
		Bold(true)
	
	ActiveTabStyle = TabStyle.Copy().
		Foreground(ActiveTabColor).
		Background(SelectedColor).
		Border(lipgloss.ThickBorder(), false, false, true, false).
		BorderForeground(ActiveTabColor)
	
	InactiveTabStyle = TabStyle.Copy().
		Foreground(InactiveTabColor)
	
	// Border styles
	BorderStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BorderColor).
		Padding(1)
	
	// Status styles
	StatusStyle = lipgloss.NewStyle().
		Foreground(TextColor).
		Background(SelectedColor).
		Padding(0, 1).
		Bold(true)
	
	// List item styles
	ListItemStyle = lipgloss.NewStyle().
		PaddingLeft(2)
	
	SelectedListItemStyle = ListItemStyle.Copy().
		Background(SelectedColor).
		Foreground(ActiveTabColor)
	
	// Header styles
	HeaderStyle = lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true).
		Padding(0, 1)
	
	// Error styles
	ErrorStyle = lipgloss.NewStyle().
		Foreground(ErrorColor).
		Bold(true)
	
	SuccessStyle = lipgloss.NewStyle().
		Foreground(SuccessColor).
		Bold(true)
	
	WarningStyle = lipgloss.NewStyle().
		Foreground(WarningColor).
		Bold(true)
	
	// Log viewer styles
	LogStyle = lipgloss.NewStyle().
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BorderColor)
	
	// Diff viewer styles
	DiffAddStyle = lipgloss.NewStyle().
		Foreground(SuccessColor).
		Background(lipgloss.Color("#0F2419"))
	
	DiffRemoveStyle = lipgloss.NewStyle().
		Foreground(ErrorColor).
		Background(lipgloss.Color("#2D1B1B"))
	
	DiffContextStyle = lipgloss.NewStyle().
		Foreground(SubtleColor)
	
	// Subtle style for less prominent text
	SubtleStyle = lipgloss.NewStyle().
		Foreground(SubtleColor)
)

// Status icons
const (
	PassIcon     = "✔"
	FailIcon     = "✖"
	SkipIcon     = "⏸"
	RunningIcon  = "⏳"
	PendingIcon  = "○"
	ExpandedIcon = "▼"
	CollapsedIcon = "▶"
)

// Utility functions for styling
func RenderCounter(label string, count int, style lipgloss.Style) string {
	if count == 0 {
		return ""
	}
	return style.Render(lipgloss.JoinHorizontal(lipgloss.Left,
		label, ": ", lipgloss.NewStyle().Bold(true).Render(string(rune(count+'0')))))
}

func RenderStatus(status string, icon string, style lipgloss.Style) string {
	return style.Render(lipgloss.JoinHorizontal(lipgloss.Left, icon, " ", status))
}
