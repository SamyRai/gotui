package tabs

import (
	"context"
	"fmt"
	"goutui/internal/runner"
	"goutui/internal/style"
	"goutui/internal/tui/components"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// BenchmarkRunner manages benchmark execution
type BenchmarkRunner struct {
	ctx       context.Context
	width     int
	height    int
	runner    *runner.CommandRunner
	parser    *runner.BenchmarkParser
	table     table.Model
	actionBar components.ActionBar
	running   bool
	status    string
}

// NewBenchmarkRunner creates a new benchmark runner
func NewBenchmarkRunner(ctx context.Context) *BenchmarkRunner {
	// Create table columns
	columns := []table.Column{
		{Title: "Name", Width: 30},
		{Title: "Iterations", Width: 12},
		{Title: "ns/op", Width: 12},
		{Title: "B/op", Width: 10},
		{Title: "allocs/op", Width: 12},
		{Title: "MB/s", Width: 10},
		{Title: "Package", Width: 20},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(15),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(style.BorderColor).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(style.ActiveTabColor).
		Background(style.SelectedColor).
		Bold(false)
	t.SetStyles(s)

	br := &BenchmarkRunner{
		ctx:       ctx,
		runner:    runner.NewCommandRunner(ctx),
		parser:    runner.NewBenchmarkParser(),
		table:     t,
		actionBar: components.NewActionBar(),
		status:    "Ready to run benchmarks",
	}
	br.updateActionBar()
	return br
}

// updateActionBar updates the action bar based on current state
func (br *BenchmarkRunner) updateActionBar() {
	br.actionBar.Clear()
	
	if br.running {
		br.actionBar.AddAction(components.Action{
			Key:         "s",
			Label:       "Stop",
			Description: "Stop running benchmarks",
			Primary:     true,
		})
	} else {
		br.actionBar.AddAction(components.Action{
			Key:         "r",
			Label:       "Run Benchmarks",
			Description: "Run all benchmarks",
			Primary:     true,
		})
	}
}

// Init initializes the benchmark runner
func (br BenchmarkRunner) Init() tea.Cmd {
	return nil
}

// SetSize sets the dimensions
func (br *BenchmarkRunner) SetSize(width, height int) {
	br.width = width
	br.height = height
	br.actionBar.SetWidth(width)
	// Account for header, status, and action bar (calculate dynamically)
	actionBarHeight := br.actionBar.Height()
	br.table.SetWidth(width - 4)
	br.table.SetHeight(height - 8 - actionBarHeight)
}

// runBenchmarks executes go test -bench
func (br BenchmarkRunner) runBenchmarks() tea.Cmd {
	return br.runner.Run("go", "test", "-bench=.", "-json", "./...")
}

// handleBenchmarkOutput processes output from go test -bench
func (br *BenchmarkRunner) handleBenchmarkOutput(output runner.CommandOutput) tea.Cmd {
	event, err := br.parser.ParseLine(output.Line)
	if err != nil {
		// Ignore parse errors for now
		return nil
	}

	if event != nil {
		// Update the table with new results
		br.updateTable()
	}

	return nil
}

// updateTable updates the table with current benchmark results
func (br *BenchmarkRunner) updateTable() {
	summary := br.parser.GetSummary()
	var rows []table.Row

	for _, result := range summary.Results {
		if result.Finished.IsZero() {
			continue // Skip running benchmarks
		}

		row := table.Row{
			result.Name,
			fmt.Sprintf("%d", result.Operations),
			fmt.Sprintf("%.2f", result.NsPerOp),
			formatBytes(result.BytesPerOp),
			fmt.Sprintf("%d", result.AllocsPerOp),
			formatMBPerSec(result.MBPerSec),
			result.Package,
		}
		rows = append(rows, row)
	}

	br.table.SetRows(rows)

	// Update status
	counts := summary.Total
	if summary.Running {
		br.status = fmt.Sprintf("Running: %d completed, %d running", counts.Completed, counts.Running)
	} else {
		br.status = fmt.Sprintf("Completed: %d benchmarks", counts.Completed)
	}
}

// formatBytes formats byte count for display
func formatBytes(bytes int64) string {
	if bytes == 0 {
		return "-"
	}
	if bytes < 1024 {
		return fmt.Sprintf("%d", bytes)
	}
	if bytes < 1024*1024 {
		return fmt.Sprintf("%.1fK", float64(bytes)/1024)
	}
	return fmt.Sprintf("%.1fM", float64(bytes)/(1024*1024))
}

// formatMBPerSec formats MB/s for display
func formatMBPerSec(mbps float64) string {
	if mbps == 0 {
		return "-"
	}
	return fmt.Sprintf("%.1f", mbps)
}

// Update handles messages
func (br *BenchmarkRunner) Update(msg tea.Msg) (TabInterface, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			if !br.running {
				br.running = true
				br.parser = runner.NewBenchmarkParser()
				br.status = "Starting benchmarks..."
				br.updateActionBar()
				return br, br.runBenchmarks()
			}
		case "s":
			if br.running {
				br.running = false
				br.parser.Finish()
				br.status = "Stopped"
				br.runner.Stop()
				br.updateActionBar()
				return br, nil
			}
		}

	case runner.CommandOutput:
		if br.running {
			cmd := br.handleBenchmarkOutput(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}

	case runner.CommandFinished:
		if br.running {
			br.running = false
			br.parser.Finish()
			br.updateTable()
			br.updateActionBar()
		}
	}

	// Update table
	var cmd tea.Cmd
	br.table, cmd = br.table.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return br, tea.Batch(cmds...)
}

// View renders the benchmark runner
func (br BenchmarkRunner) View() string {
	var parts []string

	// Action bar at the top
	actionBarView := br.actionBar.View()
	if actionBarView != "" {
		parts = append(parts, actionBarView, "")
	}

	// Header
	header := style.HeaderStyle.Render("Benchmarks")
	parts = append(parts, header, "")

	// Status
	statusBar := style.StatusStyle.Render(br.status)
	parts = append(parts, statusBar, "")

	// Table
	tableView := br.table.View()
	if tableView == "" && actionBarView != "" {
		parts = append(parts, style.SubtleStyle.Render("Press a key above to run benchmarks"))
	} else {
		parts = append(parts, tableView)
	}

	return strings.Join(parts, "\n")
}

// Refresh triggers a refresh
func (br BenchmarkRunner) Refresh() tea.Cmd {
	return nil
}

// GetStatus returns the current status
func (br BenchmarkRunner) GetStatus() string {
	return "Benchmarks ready"
}

// GetKeyHints returns key binding hints
func (br BenchmarkRunner) GetKeyHints() []string {
	return []string{"r: run benchmarks"}
}

// Cleanup performs cleanup
func (br BenchmarkRunner) Cleanup() {
}
