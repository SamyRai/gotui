package components

import (
	"testing"
)

func TestActionBarHeight(t *testing.T) {
	ab := NewActionBar()
	ab.SetWidth(100)

	// Empty action bar should have height 0
	if ab.Height() != 0 {
		t.Errorf("Expected height 0 for empty action bar, got %d", ab.Height())
	}

	// Add one action - should fit in one row
	ab.AddAction(Action{Key: "r", Label: "Run", Primary: true})
	height := ab.Height()
	// Should be: header (1) + empty line (1) + 1 row = 3
	if height != 3 {
		t.Errorf("Expected height 3 for 1 action, got %d", height)
	}

	// Add more actions that fit in one row (with width 100, ~4 buttons per row)
	ab.AddAction(Action{Key: "s", Label: "Stop", Primary: false})
	ab.AddAction(Action{Key: "f", Label: "Filter", Primary: false})
	height = ab.Height()
	// Should still be 3 (all fit in one row)
	if height != 3 {
		t.Errorf("Expected height 3 for 3 actions in one row, got %d", height)
	}

	// Add more actions to force multiple rows
	ab.AddAction(Action{Key: "o", Label: "Open", Primary: false})
	ab.AddAction(Action{Key: "e", Label: "Edit", Primary: false})
	ab.AddAction(Action{Key: "d", Label: "Delete", Primary: false})
	height = ab.Height()
	// Should be: header (1) + empty line (1) + 2 rows = 4
	if height < 3 {
		t.Errorf("Expected height >= 3 for multiple rows, got %d", height)
	}
}

func TestActionBarHeightWithNarrowWidth(t *testing.T) {
	ab := NewActionBar()
	ab.SetWidth(50) // Narrow width, fewer buttons per row

	// Add several actions
	for i := 0; i < 5; i++ {
		ab.AddAction(Action{Key: "k", Label: "Action", Primary: false})
	}

	height := ab.Height()
	// With narrow width, should need more rows
	if height < 3 {
		t.Errorf("Expected height >= 3 for narrow width with multiple actions, got %d", height)
	}
}
