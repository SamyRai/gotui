package runner

import (
	"bufio"
	"context"
	"io"
	"os/exec"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// NotificationType defines the nature of a notification.
type NotificationType int

const (
	InfoNotification NotificationType = iota
	ErrorNotification
	SuccessNotification
)

// NotificationMsg is a message used to send system-wide alerts to the TUI.
// It's for events that aren't tied to a specific command's output,
// like a command failing to start.
type NotificationMsg struct {
	Type    NotificationType
	Message string
}

// OutputSource indicates whether a line of output is from stdout or stderr.
// This allows the TUI to render them differently.
type OutputSource int

const (
	Stdout OutputSource = iota
	Stderr
)

// CommandOutput now includes the source of the output line.
type CommandOutput struct {
	Line   string
	Source OutputSource
}

// CommandResult represents the result of a command execution
type CommandResult struct {
	Command   string
	Args      []string
	ExitCode  int
	Error     error
	Duration  time.Duration
	StartTime time.Time
	EndTime   time.Time
}

// CommandStarted is sent when a command starts
type CommandStarted struct {
	Command string
	Args    []string
	PID     int
}

// CommandFinished is sent when a command completes
type CommandFinished struct {
	Result CommandResult
}

// CommandRunner manages the lifecycle of a single external command.
type CommandRunner struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup // Ensures graceful shutdown
}

// NewCommandRunner creates a new command runner.
func NewCommandRunner(ctx context.Context) *CommandRunner {
	childCtx, cancel := context.WithCancel(ctx)

	return &CommandRunner{
		ctx:    childCtx,
		cancel: cancel,
	}
}

// Run executes a command and streams its output (stdout and stderr).
func (cr *CommandRunner) Run(name string, args ...string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.CommandContext(cr.ctx, name, args...)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return NotificationMsg{Type: ErrorNotification, Message: "Failed to get stdout pipe: " + err.Error()}
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			return NotificationMsg{Type: ErrorNotification, Message: "Failed to get stderr pipe: " + err.Error()}
		}

		outputChan := make(chan tea.Msg)

		cr.wg.Add(2) // One for stdout, one for stderr

		// Goroutine to stream stdout
		go func() {
			defer cr.wg.Done()
			streamOutput(stdout, outputChan, Stdout)
		}()

		// Goroutine to stream stderr
		go func() {
			defer cr.wg.Done()
			streamOutput(stderr, outputChan, Stderr)
		}()

		if err := cmd.Start(); err != nil {
			return NotificationMsg{Type: ErrorNotification, Message: "Failed to start command: " + err.Error()}
		}

		// Goroutine to wait for command completion and send final messages
		go func() {
			cr.wg.Wait() // Wait for stdout and stderr to be fully read
			exitCode := 0
			if err := cmd.Wait(); err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					exitCode = exitErr.ExitCode()
				} else {
					exitCode = -1 // Unknown error
				}
			}
			outputChan <- CommandFinished{Result: CommandResult{ExitCode: exitCode}}
			close(outputChan)
		}()

		return <-outputChan // Return the first message from the channel
	}
}

// streamOutput reads from an io.Reader and sends lines to a channel.
func streamOutput(r io.Reader, ch chan<- tea.Msg, source OutputSource) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		ch <- CommandOutput{
			Line:   scanner.Text(),
			Source: source,
		}
	}
}

// Stop terminates the running command.
func (cr *CommandRunner) Stop() {
	cr.cancel()
}

// Wait blocks until the command and all its output streams are closed.
func (cr *CommandRunner) Wait() {
	cr.wg.Wait()
}

// Common errors
var (
	ErrCommandAlreadyRunning = &CommandError{"command already running"}
	ErrCommandNotFound       = &CommandError{"command not found"}
)

// CommandError represents a command execution error
type CommandError struct {
	message string
}

func (e *CommandError) Error() string {
	return e.message
}
