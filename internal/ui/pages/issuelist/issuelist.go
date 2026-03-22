// Package issuelist renders the left pane: JQL input + issue table (ADR-005).
package issuelist

import (
	"strings"

	"github.com/m-oehme/jiji/internal/jira"
	"github.com/m-oehme/jiji/internal/ui/common"
	"github.com/m-oehme/jiji/internal/ui/pages/issuelist/entry"
)

// Column widths for the issue table.
const (
	colKeyWidth      = 12
	colPriorityWidth = 3
	colAssigneeWidth = 15
)

// Model represents the issue list page.
type Model struct {
	common *common.Common
	issues []jira.Issue
	cursor int
	offset int // first visible row index for scrolling
	width  int
	height int
}

// New creates a new issue list page.
func New(c *common.Common) Model {
	return Model{
		common: c,
	}
}

// SetItems replaces the issue list.
func (m *Model) SetItems(issues []jira.Issue) {
	m.issues = issues
	if m.cursor >= len(issues) {
		m.cursor = max(0, len(issues)-1)
	}
	m.offset = 0
}

// SelectedIssue returns the issue at the cursor, or nil if empty.
func (m *Model) SelectedIssue() *jira.Issue {
	if len(m.issues) == 0 || m.cursor < 0 || m.cursor >= len(m.issues) {
		return nil
	}
	return &m.issues[m.cursor]
}

// SelectedIndex returns the cursor position.
func (m *Model) SelectedIndex() int {
	return m.cursor
}

// MoveUp moves the cursor up by one.
func (m *Model) MoveUp() {
	if m.cursor > 0 {
		m.cursor--
		if m.cursor < m.offset {
			m.offset = m.cursor
		}
	}
}

// MoveDown moves the cursor down by one.
func (m *Model) MoveDown() {
	if m.cursor < len(m.issues)-1 {
		m.cursor++
	}
}

// JumpToTop moves the cursor to the first issue.
func (m *Model) JumpToTop() {
	m.cursor = 0
	m.offset = 0
}

// JumpToBottom moves the cursor to the last issue.
func (m *Model) JumpToBottom() {
	if len(m.issues) > 0 {
		m.cursor = len(m.issues) - 1
	}
}

// SetFocused updates the focused state.
func (m *Model) SetFocused(focused bool) {
	m.common.Focused = focused
}

// SetSize updates the available dimensions.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// View renders the issue list pane.
func (m Model) View() string {
	if m.width <= 0 || m.height <= 0 {
		return ""
	}

	borderStyle := m.common.Styles.Border
	if m.common.Focused {
		borderStyle = m.common.Styles.BorderFocused
	}

	// In lipgloss v2, Width/Height set the TOTAL rendered size including
	// borders/padding. Content area = total - frame.
	frameW, frameH := borderStyle.GetFrameSize()
	contentW := m.width - frameW
	contentH := m.height - frameH
	if contentW <= 0 || contentH <= 0 {
		return borderStyle.Width(m.width).Height(m.height).Render("")
	}

	// Header
	header := entry.RenderHeader(m.common, contentW)

	// Row space = content height minus JQL and header lines
	rowSpace := contentH - 1
	if rowSpace < 0 {
		rowSpace = 0
	}

	// Compute visible window
	start, end := m.visibleRange(rowSpace)
	var rows []string
	for i := start; i < end; i++ {
		rows = append(rows, entry.RenderListEntry(m.common, m.issues[i], contentW, i == m.cursor))
	}

	// Pad remaining lines
	for len(rows) < rowSpace {
		rows = append(rows, strings.Repeat(" ", contentW))
	}

	content := header + "\n" + strings.Join(rows, "\n")

	// Width/Height = outer size; lipgloss subtracts the frame for content area.
	// MaxWidth/MaxHeight = hard clip safety net.
	return borderStyle.
		Width(m.width).
		Height(m.height).
		MaxWidth(m.width).
		MaxHeight(m.height).
		Render(content)
}

// visibleRange returns the slice of issues to render, keeping the cursor visible.
func (m Model) visibleRange(viewportH int) (start, end int) {
	total := len(m.issues)
	if total == 0 || viewportH <= 0 {
		return 0, 0
	}

	offset := m.offset

	// Scroll down: cursor moved below visible area
	if m.cursor >= offset+viewportH {
		offset = m.cursor - viewportH + 1
	}
	// Scroll up: cursor moved above visible area
	if m.cursor < offset {
		offset = m.cursor
	}
	// Clamp offset
	if maxOffset := total - viewportH; offset > maxOffset {
		offset = max(0, maxOffset)
	}

	start = offset
	end = start + viewportH
	if end > total {
		end = total
	}
	return start, end
}
