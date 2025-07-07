package runner

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// GoTestEvent represents a single event from `go test -json`
type GoTestEvent struct {
	Time    time.Time `json:"Time"`
	Action  string    `json:"Action"`
	Package string    `json:"Package"`
	Test    string    `json:"Test"`
	Output  string    `json:"Output"`
	Elapsed float64   `json:"Elapsed"`
}

// TestPackage represents the state of a package being tested
type TestPackage struct {
	Name     string
	Status   string // "running", "pass", "fail", "skip"
	Tests    map[string]*TestCase
	Output   []string
	Elapsed  time.Duration
	Started  time.Time
	Finished time.Time
}

// TestCase represents an individual test case
type TestCase struct {
	Name     string
	Package  string
	Status   string // "running", "pass", "fail", "skip"
	Output   []string
	Elapsed  time.Duration
	Started  time.Time
	Finished time.Time
}

// TestSummary represents overall test results
type TestSummary struct {
	Packages map[string]*TestPackage
	Total    TestCounts
	Running  bool
	Started  time.Time
	Finished time.Time
	Duration time.Duration
}

// TestCounts represents test result counters
type TestCounts struct {
	Pass int
	Fail int
	Skip int
	Run  int
}

// GoTestParser parses Go test JSON output
type GoTestParser struct {
	summary *TestSummary
}

// NewGoTestParser creates a new test parser
func NewGoTestParser() *GoTestParser {
	return &GoTestParser{
		summary: &TestSummary{
			Packages: make(map[string]*TestPackage),
			Started:  time.Now(),
			Running:  true,
		},
	}
}

// ParseLine parses a single line of JSON output from go test
func (p *GoTestParser) ParseLine(line string) (*GoTestEvent, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, nil
	}
	
	var event GoTestEvent
	if err := json.Unmarshal([]byte(line), &event); err != nil {
		// Not a JSON line, might be direct output
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	
	p.processEvent(&event)
	return &event, nil
}

// processEvent updates the test summary based on the event
func (p *GoTestParser) processEvent(event *GoTestEvent) {
	// Ensure package exists
	if event.Package != "" {
		if _, exists := p.summary.Packages[event.Package]; !exists {
			p.summary.Packages[event.Package] = &TestPackage{
				Name:    event.Package,
				Status:  "running",
				Tests:   make(map[string]*TestCase),
				Output:  make([]string, 0),
				Started: event.Time,
			}
		}
	}
	
	pkg := p.summary.Packages[event.Package]
	if pkg == nil {
		return
	}
	
	switch event.Action {
	case "run":
		p.handleRunEvent(event, pkg)
	case "pass", "fail", "skip":
		p.handleResultEvent(event, pkg)
	case "output":
		p.handleOutputEvent(event, pkg)
	case "start":
		p.handleStartEvent(event, pkg)
	case "cont":
		p.handleContinueEvent(event, pkg)
	case "pause":
		p.handlePauseEvent(event, pkg)
	}
}

// handleRunEvent handles test/package start events
func (p *GoTestParser) handleRunEvent(event *GoTestEvent, pkg *TestPackage) {
	if event.Test != "" {
		// Individual test started
		pkg.Tests[event.Test] = &TestCase{
			Name:    event.Test,
			Package: event.Package,
			Status:  "running",
			Output:  make([]string, 0),
			Started: event.Time,
		}
	} else {
		// Package started
		pkg.Status = "running"
		pkg.Started = event.Time
	}
}

// handleResultEvent handles pass/fail/skip events
func (p *GoTestParser) handleResultEvent(event *GoTestEvent, pkg *TestPackage) {
	if event.Test != "" {
		// Individual test result
		if test, exists := pkg.Tests[event.Test]; exists {
			test.Status = event.Action
			test.Finished = event.Time
			test.Elapsed = time.Duration(event.Elapsed * float64(time.Second))
		}
	} else {
		// Package result
		pkg.Status = event.Action
		pkg.Finished = event.Time
		pkg.Elapsed = time.Duration(event.Elapsed * float64(time.Second))
	}
	
	p.updateCounts()
}

// handleOutputEvent handles output events
func (p *GoTestParser) handleOutputEvent(event *GoTestEvent, pkg *TestPackage) {
	output := strings.TrimSuffix(event.Output, "\n")
	if output == "" {
		return
	}
	
	if event.Test != "" {
		// Test-specific output
		if test, exists := pkg.Tests[event.Test]; exists {
			test.Output = append(test.Output, output)
		}
	} else {
		// Package-level output
		pkg.Output = append(pkg.Output, output)
	}
}

// handleStartEvent handles package start events
func (p *GoTestParser) handleStartEvent(event *GoTestEvent, pkg *TestPackage) {
	pkg.Status = "running"
	pkg.Started = event.Time
}

// handleContinueEvent handles package continue events
func (p *GoTestParser) handleContinueEvent(event *GoTestEvent, pkg *TestPackage) {
	pkg.Status = "running"
}

// handlePauseEvent handles package pause events
func (p *GoTestParser) handlePauseEvent(event *GoTestEvent, pkg *TestPackage) {
	// Package paused (not commonly used)
}

// updateCounts recalculates test counts
func (p *GoTestParser) updateCounts() {
	counts := TestCounts{}
	
	for _, pkg := range p.summary.Packages {
		for _, test := range pkg.Tests {
			switch test.Status {
			case "pass":
				counts.Pass++
			case "fail":
				counts.Fail++
			case "skip":
				counts.Skip++
			}
			if test.Status != "" {
				counts.Run++
			}
		}
	}
	
	p.summary.Total = counts
}

// GetSummary returns the current test summary
func (p *GoTestParser) GetSummary() *TestSummary {
	return p.summary
}

// IsRunning returns true if tests are still running
func (p *GoTestParser) IsRunning() bool {
	return p.summary.Running
}

// Finish marks the test run as complete
func (p *GoTestParser) Finish() {
	p.summary.Running = false
	p.summary.Finished = time.Now()
	p.summary.Duration = p.summary.Finished.Sub(p.summary.Started)
}

// GetFailedTests returns all failed test cases
func (p *GoTestParser) GetFailedTests() []*TestCase {
	var failed []*TestCase
	
	for _, pkg := range p.summary.Packages {
		for _, test := range pkg.Tests {
			if test.Status == "fail" {
				failed = append(failed, test)
			}
		}
	}
	
	return failed
}

// GetPassedTests returns all passed test cases
func (p *GoTestParser) GetPassedTests() []*TestCase {
	var passed []*TestCase
	
	for _, pkg := range p.summary.Packages {
		for _, test := range pkg.Tests {
			if test.Status == "pass" {
				passed = append(passed, test)
			}
		}
	}
	
	return passed
}

// GetTestsByPackage returns tests organized by package
func (p *GoTestParser) GetTestsByPackage() map[string][]*TestCase {
	result := make(map[string][]*TestCase)
	
	for pkgName, pkg := range p.summary.Packages {
		tests := make([]*TestCase, 0, len(pkg.Tests))
		for _, test := range pkg.Tests {
			tests = append(tests, test)
		}
		result[pkgName] = tests
	}
	
	return result
}
