package common

import "testing"

func TestNewFocus(t *testing.T) {
	f := NewFocus()
	if f.ActivePane() != PaneIssueList {
		t.Fatalf("expected PaneIssueList, got %d", f.ActivePane())
	}
	if f.HasOverlay() {
		t.Fatal("expected no overlay")
	}
}

func TestFocus_TogglePane(t *testing.T) {
	f := NewFocus()

	f.TogglePane()
	if f.ActivePane() != PaneDetail {
		t.Fatalf("expected PaneDetail after toggle, got %d", f.ActivePane())
	}

	f.TogglePane()
	if f.ActivePane() != PaneIssueList {
		t.Fatalf("expected PaneIssueList after second toggle, got %d", f.ActivePane())
	}
}

func TestFocus_SetPane(t *testing.T) {
	f := NewFocus()

	f.SetPane(PaneDetail)
	if f.ActivePane() != PaneDetail {
		t.Fatalf("expected PaneDetail, got %d", f.ActivePane())
	}

	f.SetPane(PaneIssueList)
	if f.ActivePane() != PaneIssueList {
		t.Fatalf("expected PaneIssueList, got %d", f.ActivePane())
	}
}

func TestFocus_OverlayPushPop(t *testing.T) {
	f := NewFocus()

	if f.HasOverlay() {
		t.Fatal("expected no overlay initially")
	}

	f.PushOverlay(OverlayHelp)
	if !f.HasOverlay() {
		t.Fatal("expected overlay after push")
	}
	if f.TopOverlay() != OverlayHelp {
		t.Fatalf("expected OverlayHelp, got %d", f.TopOverlay())
	}

	f.PushOverlay(OverlayTransitionPicker)
	if f.TopOverlay() != OverlayTransitionPicker {
		t.Fatalf("expected OverlayTransitionPicker, got %d", f.TopOverlay())
	}

	popped := f.PopOverlay()
	if popped != OverlayTransitionPicker {
		t.Fatalf("expected popped OverlayTransitionPicker, got %d", popped)
	}
	if f.TopOverlay() != OverlayHelp {
		t.Fatalf("expected OverlayHelp after pop, got %d", f.TopOverlay())
	}

	popped = f.PopOverlay()
	if popped != OverlayHelp {
		t.Fatalf("expected popped OverlayHelp, got %d", popped)
	}
	if f.HasOverlay() {
		t.Fatal("expected no overlay after popping all")
	}
}

func TestFocus_PopEmptyOverlay(t *testing.T) {
	f := NewFocus()
	popped := f.PopOverlay()
	if popped != OverlayID(-1) {
		t.Fatalf("expected -1 for empty pop, got %d", popped)
	}
}

func TestFocus_TopEmptyOverlay(t *testing.T) {
	f := NewFocus()
	top := f.TopOverlay()
	if top != OverlayID(-1) {
		t.Fatalf("expected -1 for empty top, got %d", top)
	}
}

func TestFocus_PaneRing(t *testing.T) {
	f := NewFocus()

	if f.PaneRingIndex(PaneIssueList) != 0 {
		t.Fatalf("expected default ring index 0, got %d", f.PaneRingIndex(PaneIssueList))
	}

	f.SetPaneRingIndex(PaneIssueList, 3)
	if f.PaneRingIndex(PaneIssueList) != 3 {
		t.Fatalf("expected ring index 3, got %d", f.PaneRingIndex(PaneIssueList))
	}

	// Different pane has independent index
	if f.PaneRingIndex(PaneDetail) != 0 {
		t.Fatalf("expected detail ring index 0, got %d", f.PaneRingIndex(PaneDetail))
	}

	f.SetPaneRingIndex(PaneDetail, 1)
	if f.PaneRingIndex(PaneDetail) != 1 {
		t.Fatalf("expected detail ring index 1, got %d", f.PaneRingIndex(PaneDetail))
	}
}

func TestFocus_OverlayDoesNotAffectPane(t *testing.T) {
	f := NewFocus()
	f.SetPane(PaneDetail)

	f.PushOverlay(OverlayHelp)
	if f.ActivePane() != PaneDetail {
		t.Fatal("overlay push should not change active pane")
	}

	f.PopOverlay()
	if f.ActivePane() != PaneDetail {
		t.Fatal("overlay pop should not change active pane")
	}
}
