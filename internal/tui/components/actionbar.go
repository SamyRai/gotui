package components

import (
	"goutui/internal/style"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Action represents a single action button
type Action struct {
	Key         string
	Label       string
	Description string
	Primary     bool // Primary actions are more prominent
}

// ActionBar displays a responsive list of action buttons
type ActionBar struct {
	actions []Action
	width   int
}

// NewActionBar creates a new action bar
func NewActionBar() ActionBar {
	return ActionBar{
		actions: []Action{},
		width:   80,
	}
}

// SetWidth sets the width of the action bar
func (ab *ActionBar) SetWidth(width int) {
	ab.width = width
}

// SetActions sets the actions to display
func (ab *ActionBar) SetActions(actions []Action) {
	ab.actions = actions
}

// AddAction adds a single action
func (ab *ActionBar) AddAction(action Action) {
	ab.actions = append(ab.actions, action)
}

// Clear removes all actions
func (ab *ActionBar) Clear() {
	ab.actions = []Action{}
}

// View renders the action bar with responsive button layout
func (ab ActionBar) View() string {
	if len(ab.actions) == 0 {
		return ""
	}

	// Group actions by primary/secondary
	var primaryActions []Action
	var secondaryActions []Action

	for _, action := range ab.actions {
		if action.Primary {
			primaryActions = append(primaryActions, action)
		} else {
			secondaryActions = append(secondaryActions, action)
		}
	}

	var buttons []string

	// Render primary actions with prominent styling
	for _, action := range primaryActions {
		button := ab.renderButton(action, true)
		buttons = append(buttons, button)
	}

	// Render secondary actions
	for _, action := range secondaryActions {
		button := ab.renderButton(action, false)
		buttons = append(buttons, button)
	}

	// Calculate how many buttons fit per row based on width
	// Each button needs ~20-25 chars (key + label + padding)
	buttonsPerRow := ab.width / 25
	if buttonsPerRow < 1 {
		buttonsPerRow = 1
	}

	// Split buttons into rows
	var rows []string
	for i := 0; i < len(buttons); i += buttonsPerRow {
		end := i + buttonsPerRow
		if end > len(buttons) {
			end = len(buttons)
		}
		row := lipgloss.JoinHorizontal(lipgloss.Left, buttons[i:end]...)
		rows = append(rows, row)
	}

	// Add header
	header := style.HeaderStyle.Copy().
		Foreground(style.AccentColor).
		Render("Available Actions:")

	// Join rows with spacing
	content := lipgloss.JoinVertical(lipgloss.Left, rows...)

	// Combine header and buttons
	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		content,
	)
}

// renderButton renders a single action button
func (ab ActionBar) renderButton(action Action, primary bool) string {
	// Build button content: [Key] Label
	keyPart := lipgloss.NewStyle().
		Bold(true).
		Foreground(style.TextColor).
		Background(style.PrimaryColor).
		Padding(0, 1).
		Render("[" + strings.ToUpper(action.Key) + "]")

	labelPart := lipgloss.NewStyle().
		Bold(primary).
		Foreground(style.TextColor).
		Render(action.Label)

	buttonContent := lipgloss.JoinHorizontal(lipgloss.Left, keyPart, " ", labelPart)

	// Style the button
	var buttonStyle lipgloss.Style
	if primary {
		buttonStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(style.PrimaryColor).
			Padding(0, 1).
			Margin(0, 1)
	} else {
		buttonStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(style.BorderColor).
			Padding(0, 1).
			Margin(0, 1)
	}

	return buttonStyle.Render(buttonContent)
}

// GetActionByKey returns the action for a given key
func (ab ActionBar) GetActionByKey(key string) *Action {
	for i := range ab.actions {
		if strings.EqualFold(ab.actions[i].Key, key) {
			return &ab.actions[i]
		}
	}
	return nil
}

// Height calculates and returns the height of the action bar in lines
// This uses the same logic as View() to ensure consistency
func (ab ActionBar) Height() int {
	if len(ab.actions) == 0 {
		return 0
	}

	// Calculate how many buttons fit per row based on width
	// Each button needs ~20-25 chars (key + label + padding)
	buttonsPerRow := ab.width / 25
	if buttonsPerRow < 1 {
		buttonsPerRow = 1
	}

	// Calculate number of rows needed
	numButtons := len(ab.actions)
	numRows := (numButtons + buttonsPerRow - 1) / buttonsPerRow // Ceiling division
	if numRows < 1 {
		numRows = 1
	}

	// Height = header (1 line) + empty line (1 line) + button rows
	return 1 + 1 + numRows
}
