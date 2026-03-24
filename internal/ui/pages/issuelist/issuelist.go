// Package issuelist renders the left pane: JQL input + issue table (ADR-005).
package issuelist

import (
	"strings"

	"github.com/m-oehme/jiji/internal/jira"
	"github.com/m-oehme/jiji/internal/ui/common"
	"github.com/m-oehme/jiji/internal/ui/components/borderbox"
	"github.com/m-oehme/jiji/internal/ui/pages/issuelist/entry"
)

// Model represents the issue list page.
type Model struct {
	common  *common.Common
	columns []entry.Column
	issues  []jira.Issue
	cursor  int
	offset  int // first visible row index for scrolling
	width   int
	height  int
}

// New creates a new issue list page.
func New(c *common.Common, columns []entry.Column) Model {
	return Model{
		common:  c,
		columns: columns,
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

	border := borderbox.New(m.common, m.common.Focused)
	border.SetSize(m.width, m.height)
	contentW, contentH := border.GetContentSize()

	rowSpace := contentH
	if rowSpace < 0 {
		rowSpace = 0
	}

	// Reusable entry model for rendering rows.
	e := entry.New(m.common, m.columns)
	e.SetSize(contentW)

	// Compute visible window
	start, end := m.visibleRange(rowSpace)
	var rows []string
	for i := start; i < end; i++ {
		e.SetIssue(m.issues[i])
		e.SetSelected(i == m.cursor)
		rows = append(rows, e.View())
	}

	// Pad remaining lines
	for len(rows) < rowSpace {
		rows = append(rows, strings.Repeat(" ", contentW))
	}

	content := strings.Join(rows, "\n")

	return border.Render(content, "Issues")
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
