package components

import (
	"goutui/internal/style"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TreeNode represents a node in the tree
type TreeNode struct {
	ID       string
	Label    string
	Icon     string
	Status   string
	Level    int
	Expanded bool
	Children []*TreeNode
	Parent   *TreeNode
	Data     interface{} // Custom data for the node
}

// TreeList manages a hierarchical list with expand/collapse functionality
type TreeList struct {
	nodes         []*TreeNode
	flattenedView []*TreeNode
	selected      int
	height        int
	width         int
	offset        int
	filter        string
	showFiltered  bool
}

// TreeListKeyMap defines key bindings for tree navigation
type TreeListKeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Expand   key.Binding
	Collapse key.Binding
	Toggle   key.Binding
}

// DefaultTreeListKeyMap returns default key bindings
func DefaultTreeListKeyMap() TreeListKeyMap {
	return TreeListKeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Expand: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "expand"),
		),
		Collapse: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "collapse"),
		),
		Toggle: key.NewBinding(
			key.WithKeys("enter", " "),
			key.WithHelp("enter/space", "toggle"),
		),
	}
}

// NewTreeList creates a new tree list
func NewTreeList() TreeList {
	return TreeList{
		nodes:         make([]*TreeNode, 0),
		flattenedView: make([]*TreeNode, 0),
		selected:      0,
		height:        10,
		width:         40,
		offset:        0,
	}
}

// SetSize sets the dimensions of the tree list
func (tl *TreeList) SetSize(width, height int) {
	tl.width = width
	tl.height = height
}

// AddNode adds a root node to the tree
func (tl *TreeList) AddNode(node *TreeNode) {
	tl.nodes = append(tl.nodes, node)
	tl.rebuildFlattened()
}

// AddChildNode adds a child node to a parent
func (tl *TreeList) AddChildNode(parentID string, child *TreeNode) {
	parent := tl.findNodeByID(parentID)
	if parent != nil {
		child.Parent = parent
		child.Level = parent.Level + 1
		parent.Children = append(parent.Children, child)
		tl.rebuildFlattened()
	}
}

// findNodeByID recursively finds a node by ID
func (tl *TreeList) findNodeByID(id string) *TreeNode {
	return tl.findNodeByIDInList(tl.nodes, id)
}

func (tl *TreeList) findNodeByIDInList(nodes []*TreeNode, id string) *TreeNode {
	for _, node := range nodes {
		if node.ID == id {
			return node
		}
		if found := tl.findNodeByIDInList(node.Children, id); found != nil {
			return found
		}
	}
	return nil
}

// SetFilter sets a filter string for the tree
func (tl *TreeList) SetFilter(filter string) {
	tl.filter = filter
	tl.showFiltered = filter != ""
	tl.rebuildFlattened()
}

// ClearFilter clears the current filter
func (tl *TreeList) ClearFilter() {
	tl.filter = ""
	tl.showFiltered = false
	tl.rebuildFlattened()
}

// rebuildFlattened rebuilds the flattened view of the tree
func (tl *TreeList) rebuildFlattened() {
	tl.flattenedView = make([]*TreeNode, 0)
	
	if tl.showFiltered {
		tl.addFilteredNodes(tl.nodes)
	} else {
		tl.addVisibleNodes(tl.nodes)
	}
	
	// Ensure selected index is valid
	if tl.selected >= len(tl.flattenedView) {
		tl.selected = len(tl.flattenedView) - 1
	}
	if tl.selected < 0 {
		tl.selected = 0
	}
}

// addVisibleNodes adds visible nodes to the flattened view
func (tl *TreeList) addVisibleNodes(nodes []*TreeNode) {
	for _, node := range nodes {
		tl.flattenedView = append(tl.flattenedView, node)
		if node.Expanded && len(node.Children) > 0 {
			tl.addVisibleNodes(node.Children)
		}
	}
}

// addFilteredNodes adds nodes matching the filter
func (tl *TreeList) addFilteredNodes(nodes []*TreeNode) {
	for _, node := range nodes {
		if tl.matchesFilter(node) {
			tl.flattenedView = append(tl.flattenedView, node)
		}
		if len(node.Children) > 0 {
			tl.addFilteredNodes(node.Children)
		}
	}
}

// matchesFilter checks if a node matches the current filter
func (tl *TreeList) matchesFilter(node *TreeNode) bool {
	if tl.filter == "" {
		return true
	}
	return strings.Contains(strings.ToLower(node.Label), strings.ToLower(tl.filter)) ||
		strings.Contains(strings.ToLower(node.Status), strings.ToLower(tl.filter))
}

// GetSelectedNode returns the currently selected node
func (tl *TreeList) GetSelectedNode() *TreeNode {
	if tl.selected >= 0 && tl.selected < len(tl.flattenedView) {
		return tl.flattenedView[tl.selected]
	}
	return nil
}

// Update handles tree list messages
func (tl TreeList) Update(msg tea.Msg) (TreeList, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		keyMap := DefaultTreeListKeyMap()
		switch {
		case key.Matches(msg, keyMap.Up):
			tl.moveUp()
		case key.Matches(msg, keyMap.Down):
			tl.moveDown()
		case key.Matches(msg, keyMap.Expand):
			tl.expandSelected()
		case key.Matches(msg, keyMap.Collapse):
			tl.collapseSelected()
		case key.Matches(msg, keyMap.Toggle):
			tl.toggleSelected()
		}
	}
	return tl, nil
}

// moveUp moves selection up
func (tl *TreeList) moveUp() {
	if tl.selected > 0 {
		tl.selected--
		tl.ensureVisible()
	}
}

// moveDown moves selection down
func (tl *TreeList) moveDown() {
	if tl.selected < len(tl.flattenedView)-1 {
		tl.selected++
		tl.ensureVisible()
	}
}

// expandSelected expands the selected node
func (tl *TreeList) expandSelected() {
	node := tl.GetSelectedNode()
	if node != nil && len(node.Children) > 0 {
		node.Expanded = true
		tl.rebuildFlattened()
	}
}

// collapseSelected collapses the selected node
func (tl *TreeList) collapseSelected() {
	node := tl.GetSelectedNode()
	if node != nil && node.Expanded {
		node.Expanded = false
		tl.rebuildFlattened()
	}
}

// toggleSelected toggles the expanded state of the selected node
func (tl *TreeList) toggleSelected() {
	node := tl.GetSelectedNode()
	if node != nil && len(node.Children) > 0 {
		node.Expanded = !node.Expanded
		tl.rebuildFlattened()
	}
}

// ensureVisible ensures the selected item is visible in the viewport
func (tl *TreeList) ensureVisible() {
	if tl.selected < tl.offset {
		tl.offset = tl.selected
	} else if tl.selected >= tl.offset+tl.height {
		tl.offset = tl.selected - tl.height + 1
	}
}

// View renders the tree list
func (tl TreeList) View() string {
	if len(tl.flattenedView) == 0 {
		return lipgloss.NewStyle().Foreground(style.SubtleColor).Render("No items")
	}
	
	var lines []string
	start := tl.offset
	end := start + tl.height
	
	if end > len(tl.flattenedView) {
		end = len(tl.flattenedView)
	}
	
	for i := start; i < end; i++ {
		node := tl.flattenedView[i]
		line := tl.renderNode(node, i == tl.selected)
		lines = append(lines, line)
	}
	
	// Fill remaining lines if needed
	for len(lines) < tl.height {
		lines = append(lines, "")
	}
	
	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

// renderNode renders a single tree node with improved accessibility
func (tl TreeList) renderNode(node *TreeNode, selected bool) string {
	// Create indentation with better visual hierarchy
	indent := strings.Repeat("  ", node.Level)
	
	// Expansion indicator with text fallback
	var expandIcon string
	if len(node.Children) > 0 {
		if node.Expanded {
			expandIcon = style.GetStatusIcon(style.ExpandedIcon, style.ExpandedText)
		} else {
			expandIcon = style.GetStatusIcon(style.CollapsedIcon, style.CollapsedText)
		}
	} else {
		expandIcon = "  " // Two spaces for alignment
	}
	
	// Status icon with text fallback
	statusIcon := node.Icon
	if statusIcon == "" {
		statusIcon = style.GetStatusIcon(style.PendingIcon, style.PendingText)
	} else {
		// Ensure icon has text label for accessibility
		statusIcon = style.GetStatusIconOnly(statusIcon, "")
	}
	
	// Build the line with better spacing
	content := lipgloss.JoinHorizontal(
		lipgloss.Left,
		indent,
		expandIcon,
		" ",
		statusIcon,
		" ",
		node.Label,
	)
	
	// Add status if present with text labels for accessibility
	if node.Status != "" {
		var statusStyle lipgloss.Style
		var statusText string
		switch node.Status {
		case "pass", "ok":
			statusStyle = style.SuccessStyle
			statusText = style.GetStatusIcon(style.PassIcon, style.PassText) + " " + node.Status
		case "fail", "error":
			statusStyle = style.ErrorStyle
			statusText = style.GetStatusIcon(style.FailIcon, style.FailText) + " " + node.Status
		case "skip":
			statusStyle = style.WarningStyle
			statusText = style.GetStatusIcon(style.SkipIcon, style.SkipText) + " " + node.Status
		default:
			statusStyle = lipgloss.NewStyle().Foreground(style.SubtleColor)
			statusText = node.Status
		}
		
		content = lipgloss.JoinHorizontal(
			lipgloss.Left,
			content,
			" ",
			statusStyle.Render(statusText),
		)
	}
	
	// Apply selection styling with better visibility
	if selected {
		// Add selection indicator for better visibility
		selectedContent := "▶ " + content
		return style.SelectedListItemStyle.Width(tl.width).Render(selectedContent)
	}
	
	return style.ListItemStyle.Width(tl.width).Render(content)
}

// Clear removes all nodes from the tree
func (tl *TreeList) Clear() {
	tl.nodes = make([]*TreeNode, 0)
	tl.flattenedView = make([]*TreeNode, 0)
	tl.selected = 0
	tl.offset = 0
}
