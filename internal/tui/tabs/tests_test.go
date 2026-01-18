package tabs

import (
	"context"
	"fmt"
	"testing"
	"time"

	"goutui/internal/discovery"
)

func TestTestRunnerDiscovery(t *testing.T) {
	ctx := context.Background()
	tr := NewTestRunner(ctx)

	// Initially, discovery should not be done and no tests should be found
	if tr.discoveryDone {
		t.Error("Expected discoveryDone to be false initially")
	}
	if len(tr.availableTests) != 0 {
		t.Error("Expected no available tests initially")
	}

	// Simulate discovery completion
	tr.availableTests = []discovery.FileMatch{
		{Path: "/test/file1_test.go", Package: "test", RelPath: "file1_test.go"},
		{Path: "/test/file2_test.go", Package: "test", RelPath: "file2_test.go"},
	}
	tr.discoveryDone = true

	// Now it should show tests are available
	if !tr.discoveryDone {
		t.Error("Expected discoveryDone to be true after simulation")
	}
	if len(tr.availableTests) != 2 {
		t.Errorf("Expected 2 available tests, got %d", len(tr.availableTests))
	}
}

func TestTestRunnerStatus(t *testing.T) {
	ctx := context.Background()
	tr := NewTestRunner(ctx)

	// Initially should show discovering
	status := tr.GetStatus()
	if status == "" {
		t.Error("Expected non-empty status")
	}

	// Simulate discovery completion with no tests
	tr.discoveryDone = true
	status = tr.GetStatus()
	// Should show "No tests run" or similar

	// Simulate finding tests
	tr.availableTests = []discovery.FileMatch{
		{Path: "/test/file1_test.go", Package: "test", RelPath: "file1_test.go"},
	}
	status = tr.GetStatus()
	if status == "" {
		t.Error("Expected non-empty status after finding tests")
	}
}

func TestTestRunnerInit(t *testing.T) {
	ctx := context.Background()
	tr := NewTestRunner(ctx)

	// Init should return a command (not nil)
	cmd := tr.Init()
	if cmd == nil {
		t.Error("Expected Init to return a non-nil command")
	}
}

func TestTestRunnerRunningState(t *testing.T) {
	ctx := context.Background()
	tr := NewTestRunner(ctx)

	// Initially not running
	if tr.running {
		t.Error("Expected not running initially")
	}

	// Simulate starting a test run
	tr.running = true
	tr.lastRun = time.Now()

	if !tr.running {
		t.Error("Expected running after setting")
	}

	status := tr.GetStatus()
	if status == "" || status == "No tests run" {
		t.Error("Expected running status")
	}
}

func TestDiscoveryCompletedMsg(t *testing.T) {
	ctx := context.Background()
	tr := NewTestRunner(ctx)

	// Initially discovery not done
	if tr.discoveryDone {
		t.Error("Expected discovery not done initially")
	}

	// Simulate successful discovery
	msg := DiscoveryCompletedMsg{
		TestFiles: []discovery.FileMatch{
			{Path: "/test/file1_test.go", Package: "test", RelPath: "file1_test.go"},
			{Path: "/test/file2_test.go", Package: "test", RelPath: "file2_test.go"},
		},
		Error: nil,
	}

	_, cmd := tr.Update(msg)
	if cmd != nil {
		t.Error("Expected no command returned from discovery completion")
	}

	// Should now be marked as discovery done
	if !tr.discoveryDone {
		t.Error("Expected discovery done after processing message")
	}

	if len(tr.availableTests) != 2 {
		t.Errorf("Expected 2 test files, got %d", len(tr.availableTests))
	}

	status := tr.GetStatus()
	if status == "" {
		t.Error("Expected non-empty status after discovery")
	}
}

func TestDiscoveryCompletedMsgWithError(t *testing.T) {
	ctx := context.Background()
	tr := NewTestRunner(ctx)

	// Simulate failed discovery
	msg := DiscoveryCompletedMsg{
		TestFiles: nil,
		Error:     fmt.Errorf("discovery failed"),
	}

	_, cmd := tr.Update(msg)
	if cmd != nil {
		t.Error("Expected no command returned from discovery completion")
	}

	// Should still be marked as discovery done
	if !tr.discoveryDone {
		t.Error("Expected discovery done even with error")
	}

	if len(tr.availableTests) != 0 {
		t.Errorf("Expected 0 test files on error, got %d", len(tr.availableTests))
	}
}

func TestSetSizeWithActionBar(t *testing.T) {
	ctx := context.Background()
	tr := NewTestRunner(ctx)

	// Test that SetSize properly accounts for action bar height
	tr.SetSize(100, 20)

	// The action bar should reduce available height for components
	// We can't easily test the internal component sizes, but we can verify
	// the method doesn't panic and sets the expected width/height
	if tr.width != 100 {
		t.Errorf("Expected width 100, got %d", tr.width)
	}
	if tr.height != 20 {
		t.Errorf("Expected height 20, got %d", tr.height)
	}
}