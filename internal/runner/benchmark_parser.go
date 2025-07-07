package runner

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// BenchmarkResult represents a single benchmark result
type BenchmarkResult struct {
	Name         string
	Package      string
	Operations   int64
	NsPerOp      float64
	BytesPerOp   int64
	AllocsPerOp  int64
	MBPerSec     float64
	Iterations   int
	Parallelism  int
	Raw          string
	Output       []string
	Started      time.Time
	Finished     time.Time
	Duration     time.Duration
}

// BenchmarkSummary represents all benchmark results
type BenchmarkSummary struct {
	Results  map[string]*BenchmarkResult
	Packages map[string]*BenchmarkPackage
	Started  time.Time
	Finished time.Time
	Duration time.Duration
	Running  bool
	Total    BenchmarkCounts
}

// BenchmarkPackage represents benchmark results for a package
type BenchmarkPackage struct {
	Name      string
	Status    string // "running", "pass", "fail"
	Results   map[string]*BenchmarkResult
	Output    []string
	Started   time.Time
	Finished  time.Time
	Duration  time.Duration
}

// BenchmarkCounts represents benchmark result counters
type BenchmarkCounts struct {
	Total     int
	Completed int
	Failed    int
	Running   int
}

// BenchmarkParser parses Go benchmark output
type BenchmarkParser struct {
	summary *BenchmarkSummary
}

// NewBenchmarkParser creates a new benchmark parser
func NewBenchmarkParser() *BenchmarkParser {
	return &BenchmarkParser{
		summary: &BenchmarkSummary{
			Results:  make(map[string]*BenchmarkResult),
			Packages: make(map[string]*BenchmarkPackage),
			Started:  time.Now(),
			Running:  true,
		},
	}
}

// ParseLine parses a line of benchmark output
func (p *BenchmarkParser) ParseLine(line string) (*GoTestEvent, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, nil
	}

	// Try to parse as JSON first (from go test -bench -json)
	var event GoTestEvent
	if err := json.Unmarshal([]byte(line), &event); err == nil {
		p.processEvent(&event)
		return &event, nil
	}

	// Fall back to parsing traditional benchmark output
	if strings.HasPrefix(line, "Benchmark") {
		if err := p.parseBenchmarkLine(line, ""); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

// processEvent processes a JSON test event that might contain benchmark data
func (p *BenchmarkParser) processEvent(event *GoTestEvent) {
	// Ensure package exists
	if event.Package != "" {
		if _, exists := p.summary.Packages[event.Package]; !exists {
			p.summary.Packages[event.Package] = &BenchmarkPackage{
				Name:    event.Package,
				Status:  "running",
				Results: make(map[string]*BenchmarkResult),
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
	case "output":
		p.handleOutputEvent(event, pkg)
	case "pass", "fail":
		if event.Test == "" {
			// Package-level result
			pkg.Status = event.Action
			pkg.Finished = event.Time
			pkg.Duration = time.Duration(event.Elapsed * float64(time.Second))
		}
	case "run":
		if event.Test != "" && strings.HasPrefix(event.Test, "Benchmark") {
			// Benchmark started
			result := &BenchmarkResult{
				Name:    event.Test,
				Package: event.Package,
				Output:  make([]string, 0),
				Started: event.Time,
			}
			pkg.Results[event.Test] = result
			p.summary.Results[event.Package+"/"+event.Test] = result
		}
	}

	p.updateCounts()
}

// handleOutputEvent processes output from benchmark execution
func (p *BenchmarkParser) handleOutputEvent(event *GoTestEvent, pkg *BenchmarkPackage) {
	output := strings.TrimSuffix(event.Output, "\n")
	if output == "" {
		return
	}

	// Check if this is benchmark result output
	if strings.HasPrefix(output, "Benchmark") && strings.Contains(output, "ns/op") {
		if err := p.parseBenchmarkLine(output, event.Package); err == nil {
			return
		}
	}

	// Add to package output
	pkg.Output = append(pkg.Output, output)

	// If we have a current test, add to its output too
	if event.Test != "" {
		if result, exists := pkg.Results[event.Test]; exists {
			result.Output = append(result.Output, output)
		}
	}
}

// parseBenchmarkLine parses a benchmark result line
// Example: BenchmarkFoo-8    	    1000	   2000 ns/op	   1000 B/op	      2 allocs/op
func (p *BenchmarkParser) parseBenchmarkLine(line, pkg string) error {
	parts := strings.Fields(line)
	if len(parts) < 4 {
		return fmt.Errorf("invalid benchmark line format")
	}

	// Parse benchmark name and parallelism
	namePart := parts[0]
	var name string
	var parallelism int = 1

	if idx := strings.LastIndex(namePart, "-"); idx != -1 {
		name = namePart[:idx]
		if p, err := strconv.Atoi(namePart[idx+1:]); err == nil {
			parallelism = p
		}
	} else {
		name = namePart
	}

	// Parse iterations
	iterations, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("failed to parse iterations: %w", err)
	}

	// Parse ns/op
	nsOpStr := strings.TrimSuffix(parts[2], "ns/op")
	nsPerOp, err := strconv.ParseFloat(nsOpStr, 64)
	if err != nil {
		return fmt.Errorf("failed to parse ns/op: %w", err)
	}

	result := &BenchmarkResult{
		Name:        name,
		Package:     pkg,
		Operations:  int64(iterations),
		NsPerOp:     nsPerOp,
		Iterations:  iterations,
		Parallelism: parallelism,
		Raw:         line,
		Finished:    time.Now(),
	}

	// Parse additional metrics if present
	for i := 3; i < len(parts); i++ {
		part := parts[i]
		if strings.HasSuffix(part, "B/op") {
			if bytesStr := strings.TrimSuffix(part, "B/op"); bytesStr != "" {
				if bytes, err := strconv.ParseInt(bytesStr, 10, 64); err == nil {
					result.BytesPerOp = bytes
				}
			}
		} else if strings.HasSuffix(part, "allocs/op") {
			if allocsStr := strings.TrimSuffix(part, "allocs/op"); allocsStr != "" {
				if allocs, err := strconv.ParseInt(allocsStr, 10, 64); err == nil {
					result.AllocsPerOp = allocs
				}
			}
		} else if strings.HasSuffix(part, "MB/s") {
			if mbStr := strings.TrimSuffix(part, "MB/s"); mbStr != "" {
				if mb, err := strconv.ParseFloat(mbStr, 64); err == nil {
					result.MBPerSec = mb
				}
			}
		}
	}

	// Store the result
	key := result.Package + "/" + result.Name
	p.summary.Results[key] = result

	// Also store in package if we have one
	if pkg != "" {
		if pkgData, exists := p.summary.Packages[pkg]; exists {
			pkgData.Results[result.Name] = result
		}
	}

	return nil
}

// updateCounts recalculates benchmark counts
func (p *BenchmarkParser) updateCounts() {
	counts := BenchmarkCounts{}

	for _, pkg := range p.summary.Packages {
		for _, result := range pkg.Results {
			counts.Total++
			if result.Finished.IsZero() {
				counts.Running++
			} else {
				counts.Completed++
			}
		}
	}

	p.summary.Total = counts
}

// GetSummary returns the current benchmark summary
func (p *BenchmarkParser) GetSummary() *BenchmarkSummary {
	return p.summary
}

// IsRunning returns true if benchmarks are still running
func (p *BenchmarkParser) IsRunning() bool {
	return p.summary.Running
}

// Finish marks the benchmark run as complete
func (p *BenchmarkParser) Finish() {
	p.summary.Running = false
	p.summary.Finished = time.Now()
	p.summary.Duration = p.summary.Finished.Sub(p.summary.Started)
}

// GetResultsByPackage returns benchmark results organized by package
func (p *BenchmarkParser) GetResultsByPackage() map[string][]*BenchmarkResult {
	result := make(map[string][]*BenchmarkResult)

	for pkgName, pkg := range p.summary.Packages {
		results := make([]*BenchmarkResult, 0, len(pkg.Results))
		for _, benchResult := range pkg.Results {
			results = append(results, benchResult)
		}
		result[pkgName] = results
	}

	return result
}

// GetTopResults returns the top N benchmark results by performance
func (p *BenchmarkParser) GetTopResults(n int) []*BenchmarkResult {
	var results []*BenchmarkResult
	for _, result := range p.summary.Results {
		if !result.Finished.IsZero() {
			results = append(results, result)
		}
	}

	// Sort by ns/op ascending (fastest first)
	sort.Slice(results, func(i, j int) bool {
		return results[i].NsPerOp < results[j].NsPerOp
	})

	if len(results) > n {
		return results[:n]
	}
	return results
}

// CompareTo compares this benchmark run with another (for future benchstat integration)
func (p *BenchmarkParser) CompareTo(other *BenchmarkParser) map[string]*BenchmarkComparison {
	comparisons := make(map[string]*BenchmarkComparison)
	for name, oldResult := range p.summary.Results {
		if newResult, ok := other.summary.Results[name]; ok {
			change := float64(newResult.NsPerOp-oldResult.NsPerOp) / float64(oldResult.NsPerOp) * 100
			comparisons[name] = &BenchmarkComparison{
				Name:   name,
				Old:    oldResult,
				New:    newResult,
				Change: change,
				Better: newResult.NsPerOp < oldResult.NsPerOp,
			}
		}
	}
	return comparisons
}

// BenchmarkComparison represents a comparison between two benchmark results
type BenchmarkComparison struct {
	Name     string
	Old      *BenchmarkResult
	New      *BenchmarkResult
	Change   float64 // Percentage change in ns/op
	Better   bool    // True if new is better than old
}
