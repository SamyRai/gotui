package components

import (
	"goutui/internal/style"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// DiffViewer displays file diffs with syntax highlighting
type DiffViewer struct {
	viewport    viewport.Model
	width       int
	height      int
	title       string
	borderStyle lipgloss.Style
	diffLines   []DiffLine
	showLineNumbers bool
}

// DiffLine represents a single line in a diff
type DiffLine struct {
	Content    string
	Type       DiffLineType
	LineNum    int
	OldLineNum int
	NewLineNum int
}

// DiffLineType represents the type of a diff line
type DiffLineType int

const (
	DiffLineContext DiffLineType = iota
	DiffLineAddition
	DiffLineRemoval
	DiffLineHeader
	DiffLineHunk
	DiffLineNoNewline
)

// NewDiffViewer creates a new diff viewer
func NewDiffViewer(title string) DiffViewer {
	vp := viewport.New(80, 20)
	vp.Style = style.LogStyle
	
	return DiffViewer{
		viewport:        vp,
		title:           title,
		borderStyle:     style.BorderStyle,
		diffLines:       make([]DiffLine, 0),
		showLineNumbers: true,
	}
}

// SetSize sets the dimensions of the diff viewer
func (dv *DiffViewer) SetSize(width, height int) {
	dv.width = width
	dv.height = height
	dv.viewport.Width = width - 2   // Account for border
	dv.viewport.Height = height - 3 // Account for border and title
}

// SetDiff sets the diff content
func (dv *DiffViewer) SetDiff(diff string) {
	dv.diffLines = dv.parseDiff(diff)
	dv.updateContent()
}

// SetShowLineNumbers toggles line number display
func (dv *DiffViewer) SetShowLineNumbers(show bool) {
	dv.showLineNumbers = show
	dv.updateContent()
}

// parseDiff parses a unified diff
func (dv *DiffViewer) parseDiff(diff string) []DiffLine {
	lines := strings.Split(diff, "\n")
	var diffLines []DiffLine
	var oldLineNum, newLineNum int

	for i, line := range lines {
		diffLine := DiffLine{
			Content:    line,
			LineNum:    i + 1,
			OldLineNum: oldLineNum,
			NewLineNum: newLineNum,
		}

		if len(line) == 0 {
			diffLine.Type = DiffLineContext
			oldLineNum++
			newLineNum++
		} else {
			switch line[0] {
			case '+':
				if strings.HasPrefix(line, "+++") {
					diffLine.Type = DiffLineHeader
				} else {
					diffLine.Type = DiffLineAddition
					newLineNum++
				}
			case '-':
				if strings.HasPrefix(line, "---") {
					diffLine.Type = DiffLineHeader
				} else {
					diffLine.Type = DiffLineRemoval
					oldLineNum++
				}
			case '@':
				if strings.HasPrefix(line, "@@") {
					diffLine.Type = DiffLineHunk
					// Parse hunk header to extract line numbers
					if parts := strings.Fields(line); len(parts) >= 3 {
						// Example: @@ -1,4 +1,6 @@
						oldPart := strings.TrimPrefix(parts[1], "-")
						newPart := strings.TrimPrefix(parts[2], "+")
						
						if oldComma := strings.Split(oldPart, ","); len(oldComma) > 0 {
							if num := parseInt(oldComma[0]); num > 0 {
								oldLineNum = num - 1
							}
						}
						if newComma := strings.Split(newPart, ","); len(newComma) > 0 {
							if num := parseInt(newComma[0]); num > 0 {
								newLineNum = num - 1
							}
						}
					}
				} else {
					diffLine.Type = DiffLineContext
					oldLineNum++
					newLineNum++
				}
			case '\\':
				if strings.HasPrefix(line, "\\ No newline") {
					diffLine.Type = DiffLineNoNewline
				} else {
					diffLine.Type = DiffLineContext
					oldLineNum++
					newLineNum++
				}
			default:
				if strings.HasPrefix(line, "diff --git") || 
				   strings.HasPrefix(line, "index ") ||
				   strings.HasPrefix(line, "new file") ||
				   strings.HasPrefix(line, "deleted file") {
					diffLine.Type = DiffLineHeader
				} else {
					diffLine.Type = DiffLineContext
					oldLineNum++
					newLineNum++
				}
			}
		}

		diffLine.OldLineNum = oldLineNum
		diffLine.NewLineNum = newLineNum
		diffLines = append(diffLines, diffLine)
	}

	return diffLines
}

// parseInt safely parses an integer string
func parseInt(s string) int {
	var result int
	for _, r := range s {
		if r >= '0' && r <= '9' {
			result = result*10 + int(r-'0')
		} else {
			break
		}
	}
	return result
}

// updateContent updates the viewport content
func (dv *DiffViewer) updateContent() {
	var content strings.Builder
	
	for _, line := range dv.diffLines {
		lineContent := dv.formatDiffLine(line)
		content.WriteString(lineContent + "\n")
	}
	
	dv.viewport.SetContent(content.String())
}

// formatDiffLine formats a single diff line with appropriate styling
func (dv *DiffViewer) formatDiffLine(line DiffLine) string {
	var prefix string
	var styledContent string

	// Add line numbers if enabled
	if dv.showLineNumbers && line.Type != DiffLineHeader && line.Type != DiffLineHunk {
		oldNum := " "
		newNum := " "
		
		if line.Type == DiffLineContext || line.Type == DiffLineRemoval {
			if line.OldLineNum > 0 {
				oldNum = lipgloss.NewStyle().
					Width(4).
					Align(lipgloss.Right).
					Foreground(style.SubtleColor).
					Render(string(rune(line.OldLineNum)))
			}
		}
		
		if line.Type == DiffLineContext || line.Type == DiffLineAddition {
			if line.NewLineNum > 0 {
				newNum = lipgloss.NewStyle().
					Width(4).
					Align(lipgloss.Right).
					Foreground(style.SubtleColor).
					Render(string(rune(line.NewLineNum)))
			}
		}
		
		prefix = lipgloss.NewStyle().
			Foreground(style.SubtleColor).
			Render(oldNum + " " + newNum + " ")
	}

	// Apply styling based on line type
	switch line.Type {
	case DiffLineAddition:
		styledContent = style.DiffAddStyle.Render(line.Content)
	case DiffLineRemoval:
		styledContent = style.DiffRemoveStyle.Render(line.Content)
	case DiffLineHeader:
		styledContent = lipgloss.NewStyle().
			Bold(true).
			Foreground(style.InfoColor).
			Render(line.Content)
	case DiffLineHunk:
		styledContent = lipgloss.NewStyle().
			Bold(true).
			Foreground(style.AccentColor).
			Background(lipgloss.Color("#2D2A1F")).
			Render(line.Content)
	case DiffLineNoNewline:
		styledContent = lipgloss.NewStyle().
			Italic(true).
			Foreground(style.SubtleColor).
			Render(line.Content)
	default:
		styledContent = style.DiffContextStyle.Render(line.Content)
	}

	return prefix + styledContent
}

// Update handles messages
func (dv DiffViewer) Update(msg tea.Msg) (DiffViewer, tea.Cmd) {
	var cmd tea.Cmd
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "n":
			// Toggle line numbers
			dv.showLineNumbers = !dv.showLineNumbers
			dv.updateContent()
			return dv, nil
		}
	}
	
	dv.viewport, cmd = dv.viewport.Update(msg)
	return dv, cmd
}

// View renders the diff viewer
func (dv DiffViewer) View() string {
	titleBar := style.HeaderStyle.Render(dv.title)
	
	// Add key hints
	keyHints := lipgloss.NewStyle().
		Foreground(style.SubtleColor).
		Render("n: toggle line numbers • ↑↓: scroll • q: back")
	
	content := dv.borderStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			titleBar,
			dv.viewport.View(),
			keyHints,
		),
	)
	
	return content
}

// GetContent returns just the viewport content without borders or styling
func (dv DiffViewer) GetContent() string {
	return dv.viewport.View()
}

// GetTitle returns the title of the diff viewer
func (dv DiffViewer) GetTitle() string {
	return dv.title
}

// GetStats returns statistics about the diff
func (dv DiffViewer) GetStats() DiffStats {
	var stats DiffStats
	
	for _, line := range dv.diffLines {
		switch line.Type {
		case DiffLineAddition:
			stats.Additions++
		case DiffLineRemoval:
			stats.Removals++
		case DiffLineContext:
			stats.Context++
		}
	}
	
	return stats
}

// DiffStats represents statistics about a diff
type DiffStats struct {
	Additions int
	Removals  int
	Context   int
}

// String returns a string representation of the diff stats
func (ds DiffStats) String() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		style.SuccessStyle.Render("+"+string(rune(ds.Additions))),
		" ",
		style.ErrorStyle.Render("-"+string(rune(ds.Removals))),
	)
}

// IsEmpty returns true if there are no changes in the diff
func (ds DiffStats) IsEmpty() bool {
	return ds.Additions == 0 && ds.Removals == 0
}

// GetRawDiff returns the raw diff as a string (for accumulation)
func (dv *DiffViewer) GetRawDiff() string {
	var lines []string
	for _, l := range dv.diffLines {
		lines = append(lines, l.Content)
	}
	return strings.Join(lines, "\n")
}

// GetFirstFilePathFromDiff returns the first file path found in the diff header
func (dv *DiffViewer) GetFirstFilePathFromDiff() string {
	for _, l := range dv.diffLines {
		if l.Type == DiffLineHeader && strings.HasPrefix(l.Content, "+++ ") {
			// Example: +++ b/internal/foo.go
			parts := strings.Fields(l.Content)
			if len(parts) > 1 {
				path := parts[1]
				// Remove leading b/ or a/ if present
				if len(path) > 2 && (path[:2] == "b/" || path[:2] == "a/") {
					return path[2:]
				}
				return path
			}
		}
	}
	return ""
}
