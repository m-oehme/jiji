package common

// PaneID identifies a pane in the split layout (ADR-009 Layer 2).
type PaneID int

const (
	PaneIssueList PaneID = iota
	PaneDetail
	paneCount // sentinel for toggling
)

// OverlayID identifies a modal overlay (ADR-009 Layer 3).
type OverlayID int

const (
	OverlayHelp OverlayID = iota
	OverlayTransitionPicker
	OverlayLabelPicker
	OverlayAutocompleteDrop
)

// Focus implements the three-layer focus system from ADR-009:
//
//	Layer 3: Overlay stack (modals capture all input)
//	Layer 2: Pane switching (Tab toggles between issue list and detail)
//	Layer 1: Intra-pane ring (element focus within a pane)
type Focus struct {
	overlays   []OverlayID
	activePane PaneID
	paneRings  map[PaneID]int // focused element index within pane
}

// NewFocus returns a Focus with PaneIssueList active and no overlays.
func NewFocus() *Focus {
	return &Focus{
		activePane: PaneIssueList,
		paneRings:  make(map[PaneID]int),
	}
}

// --- Layer 3: Overlay management ---

// PushOverlay adds an overlay to the top of the stack.
func (f *Focus) PushOverlay(id OverlayID) {
	f.overlays = append(f.overlays, id)
}

// PopOverlay removes and returns the topmost overlay.
// Returns -1 if the stack is empty.
func (f *Focus) PopOverlay() OverlayID {
	if len(f.overlays) == 0 {
		return OverlayID(-1)
	}
	top := f.overlays[len(f.overlays)-1]
	f.overlays = f.overlays[:len(f.overlays)-1]
	return top
}

// HasOverlay returns true if any overlay is active.
func (f *Focus) HasOverlay() bool {
	return len(f.overlays) > 0
}

// TopOverlay returns the topmost overlay without removing it.
// Returns -1 if the stack is empty.
func (f *Focus) TopOverlay() OverlayID {
	if len(f.overlays) == 0 {
		return OverlayID(-1)
	}
	return f.overlays[len(f.overlays)-1]
}

// --- Layer 2: Pane switching ---

// SetPane sets the active pane.
func (f *Focus) SetPane(id PaneID) {
	f.activePane = id
}

// ActivePane returns the currently active pane.
func (f *Focus) ActivePane() PaneID {
	return f.activePane
}

// TogglePane switches between IssueList and Detail panes.
func (f *Focus) TogglePane() {
	f.activePane = (f.activePane + 1) % paneCount
}

// --- Layer 1: Intra-pane ring ---

// PaneRingIndex returns the focused element index within the given pane.
func (f *Focus) PaneRingIndex(pane PaneID) int {
	return f.paneRings[pane]
}

// SetPaneRingIndex sets the focused element index within the given pane.
func (f *Focus) SetPaneRingIndex(pane PaneID, idx int) {
	f.paneRings[pane] = idx
}
