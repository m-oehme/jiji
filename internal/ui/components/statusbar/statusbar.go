// Package statusbar renders the bottom status bar (ADR-010).
package statusbar

import (
	"fmt"
	"strings"

	"github.com/m-oehme/jiji/internal/ui/styles"
	lipgloss "charm.land/lipgloss/v2"
)

// Model represents the status bar at the bottom of the screen.
type Model struct {
	styles       *styles.Styles
	width        int
	height       int
	loading      bool
	errMsg       string
	status       string
	issueCount   int
	currentIssue string
}

// New creates a new status bar.
func New(s *styles.Styles) Model {
	return Model{
		styles: s,
		height: 1,
	}
}

// SetLoading toggles the loading indicator.
func (m *Model) SetLoading(loading bool) {
	m.loading = loading
}

// SetError sets a temporary error message.
func (m *Model) SetError(err error) {
	if err != nil {
		m.errMsg = err.Error()
	} else {
		m.errMsg = ""
	}
}

// ClearError clears the current error message.
func (m *Model) ClearError() {
	m.errMsg = ""
}

// SetStatus sets a custom status message.
func (m *Model) SetStatus(text string) {
	m.status = text
}

// SetIssueCount sets the displayed issue count.
func (m *Model) SetIssueCount(n int) {
	m.issueCount = n
}

// SetCurrentIssue sets the current issue key displayed on the left.
func (m *Model) SetCurrentIssue(key string) {
	m.currentIssue = key
}

// SetSize updates the available width for the status bar.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// View renders the status bar.
func (m Model) View() string {
	if m.width <= 0 {
		return ""
	}

	// Error takes priority
	if m.errMsg != "" {
		errText := m.styles.Error.Render(" " + m.errMsg + " ")
		return m.pad(errText)
	}

	// Left: current issue key
	left := ""
	if m.currentIssue != "" {
		left = " " + m.currentIssue
	}

	// Center: issue count or loading or custom status
	center := ""
	if m.loading {
		center = "loading…"
	} else if m.status != "" {
		center = m.status
	} else if m.issueCount > 0 {
		center = fmt.Sprintf("%d issues", m.issueCount)
	}

	// Right: help hint
	right := "? help "

	return m.compose(left, center, right)
}

// compose lays out left/center/right across the full width.
func (m Model) compose(left, center, right string) string {
	lw := lipgloss.Width(left)
	cw := lipgloss.Width(center)
	rw := lipgloss.Width(right)

	// Available space between left and right
	totalContent := lw + cw + rw
	if totalContent >= m.width {
		// Truncate center if needed
		avail := m.width - lw - rw
		if avail > 0 && len(center) > avail {
			center = center[:avail-1] + "…"
			cw = lipgloss.Width(center)
		} else if avail <= 0 {
			center = ""
			cw = 0
		}
	}

	// Build with separator chars
	leftPad := ""
	if lw > 0 && cw > 0 {
		leftPad = " │ "
	}

	content := left + leftPad + center
	contentW := lipgloss.Width(content)
	gap := m.width - contentW - rw
	if gap < 0 {
		gap = 0
	}

	bar := content + strings.Repeat(" ", gap) + right
	return m.styles.StatusBar.Width(m.width).Render(bar)
}

// pad renders text centered in the full-width status bar.
func (m Model) pad(text string) string {
	return m.styles.StatusBar.Width(m.width).Render(text)
}
