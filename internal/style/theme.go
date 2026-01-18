package style

import (
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	// noColor indicates if colors should be disabled (NO_COLOR env var)
	noColor = false
	
	// darkBackground indicates if terminal has dark background
	darkBackground = true
)

func init() {
	// Check for NO_COLOR environment variable (accessibility best practice)
	if _, exists := os.LookupEnv("NO_COLOR"); exists {
		noColor = true
	}
	
	// Detect terminal background color using lipgloss
	// Default to dark background if detection fails
	darkBackground = lipgloss.HasDarkBackground()
}

// getColor returns a color or empty string if colors are disabled
func getColor(color lipgloss.Color) lipgloss.Color {
	if noColor {
		return ""
	}
	return color
}

// Color scheme for the TUI
var (
	// Primary colors - high contrast for accessibility
	PrimaryColor   = lipgloss.Color("#00D4AA")
	SecondaryColor = lipgloss.Color("#7C3AED")
	AccentColor    = lipgloss.Color("#F59E0B")
	
	// Status colors - high contrast, colorblind-friendly alternatives
	SuccessColor = lipgloss.Color("#10B981")
	ErrorColor   = lipgloss.Color("#EF4444")
	WarningColor = lipgloss.Color("#F59E0B")
	InfoColor    = lipgloss.Color("#3B82F6")
	
	// UI colors - optimized for dark terminals
	BorderColor     = lipgloss.Color("#374151")
	SelectedColor   = lipgloss.Color("#1F2937")
	BackgroundColor = lipgloss.Color("#111827")
	TextColor       = lipgloss.Color("#F9FAFB")
	SubtleColor     = lipgloss.Color("#9CA3AF")
	
	// Tab colors
	ActiveTabColor   = PrimaryColor
	InactiveTabColor = SubtleColor
)

// Common styles - with NO_COLOR support
var (
	// Tab styles
	TabStyle = lipgloss.NewStyle().
		Padding(0, 2).
		Bold(true)
	
	ActiveTabStyle = TabStyle.Copy().
		Foreground(getColor(ActiveTabColor)).
		Background(getColor(SelectedColor)).
		Border(lipgloss.ThickBorder(), false, false, true, false).
		BorderForeground(getColor(ActiveTabColor))
	
	InactiveTabStyle = TabStyle.Copy().
		Foreground(getColor(InactiveTabColor))
	
	// Border styles
	BorderStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(getColor(BorderColor)).
		Padding(1)
	
	// Status styles - high contrast for accessibility
	StatusStyle = lipgloss.NewStyle().
		Foreground(getColor(TextColor)).
		Background(getColor(SelectedColor)).
		Padding(0, 1).
		Bold(true)
	
	// List item styles
	ListItemStyle = lipgloss.NewStyle().
		PaddingLeft(2)
	
	SelectedListItemStyle = ListItemStyle.Copy().
		Background(getColor(SelectedColor)).
		Foreground(getColor(ActiveTabColor)).
		Bold(true) // Add bold for better visibility without color
	
	// Header styles
	HeaderStyle = lipgloss.NewStyle().
		Foreground(getColor(PrimaryColor)).
		Bold(true).
		Padding(0, 1)
	
	// Error styles - with text indicators for accessibility
	ErrorStyle = lipgloss.NewStyle().
		Foreground(getColor(ErrorColor)).
		Bold(true)
	
	SuccessStyle = lipgloss.NewStyle().
		Foreground(getColor(SuccessColor)).
		Bold(true)
	
	WarningStyle = lipgloss.NewStyle().
		Foreground(getColor(WarningColor)).
		Bold(true)
	
	// Log viewer styles
	LogStyle = lipgloss.NewStyle().
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(getColor(BorderColor))
	
	// Diff viewer styles
	DiffAddStyle = lipgloss.NewStyle().
		Foreground(getColor(SuccessColor)).
		Background(getColor(lipgloss.Color("#0F2419"))).
		Bold(true) // Bold for better visibility
	
	DiffRemoveStyle = lipgloss.NewStyle().
		Foreground(getColor(ErrorColor)).
		Background(getColor(lipgloss.Color("#2D1B1B"))).
		Bold(true) // Bold for better visibility
	
	DiffContextStyle = lipgloss.NewStyle().
		Foreground(getColor(SubtleColor))
	
	// Subtle style for less prominent text
	SubtleStyle = lipgloss.NewStyle().
		Foreground(getColor(SubtleColor))
	
	// Focus indicator style - visible even without color
	FocusStyle = lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(getColor(PrimaryColor)).
		Bold(true)
	
	// Notification styles
	NotificationSuccessStyle = lipgloss.NewStyle().
		Background(getColor(SuccessColor)).
		Foreground(getColor(TextColor)).
		Bold(true).
		Padding(0, 1)
	
	NotificationErrorStyle = lipgloss.NewStyle().
		Background(getColor(ErrorColor)).
		Foreground(getColor(TextColor)).
		Bold(true).
		Padding(0, 1)
	
	NotificationInfoStyle = lipgloss.NewStyle().
		Background(getColor(InfoColor)).
		Foreground(getColor(TextColor)).
		Bold(true).
		Padding(0, 1)
)

// Status icons - with text fallbacks for accessibility
const (
	PassIcon     = "✔"
	FailIcon     = "✖"
	SkipIcon     = "⏸"
	RunningIcon  = "⏳"
	PendingIcon  = "○"
	ExpandedIcon = "▼"
	CollapsedIcon = "▶"
	
	// Text-based alternatives for accessibility
	PassText     = "[PASS]"
	FailText     = "[FAIL]"
	SkipText     = "[SKIP]"
	RunningText  = "[RUNNING]"
	PendingText  = "[PENDING]"
	ExpandedText = "[EXPANDED]"
	CollapsedText = "[COLLAPSED]"
)

// GetStatusIcon returns icon with text label for accessibility
func GetStatusIcon(icon, text string) string {
	if noColor {
		return text
	}
	return icon + " " + text
}

// GetStatusIconOnly returns just the icon or text alternative
func GetStatusIconOnly(icon, text string) string {
	if noColor {
		return text
	}
	return icon
}

// Utility functions for styling
func RenderCounter(label string, count int, style lipgloss.Style) string {
	if count == 0 {
		return ""
	}
	countStr := ""
	if count < 10 {
		countStr = string(rune(count + '0'))
	} else {
		countStr = strings.Repeat("9", count/9) + string(rune(count%9+'0'))
	}
	// Better formatting for numbers
	countStr = formatNumber(count)
	return style.Render(lipgloss.JoinHorizontal(lipgloss.Left,
		label, ": ", lipgloss.NewStyle().Bold(true).Render(countStr)))
}

func formatNumber(n int) string {
	if n < 10 {
		return string(rune(n + '0'))
	}
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(
		strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(
			strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(
				strings.ReplaceAll("0123456789", "0", ""),
				"1", ""), "2", ""), "3", ""), "4", ""), "5", ""),
			"6", ""), "7", ""), "8", ""), "9", "")
}

func RenderStatus(status string, icon string, style lipgloss.Style) string {
	// Always include text label for accessibility
	displayIcon := GetStatusIconOnly(icon, "")
	if displayIcon != "" {
		return style.Render(lipgloss.JoinHorizontal(lipgloss.Left, displayIcon, " ", status))
	}
	return style.Render(status)
}

// RenderStatusWithText renders status with both icon and text label
func RenderStatusWithText(status string, icon string, textLabel string, style lipgloss.Style) string {
	displayIcon := GetStatusIcon(icon, textLabel)
	return style.Render(lipgloss.JoinHorizontal(lipgloss.Left, displayIcon, " ", status))
}
