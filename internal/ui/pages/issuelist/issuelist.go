// Package issuelist renders the left pane: JQL input + issue table (ADR-005).
package issuelist

import (
	"strings"

	"github.com/m-oehme/jiji/internal/config"
	"github.com/m-oehme/jiji/internal/jira"
	"github.com/m-oehme/jiji/internal/ui/common"
	"github.com/m-oehme/jiji/internal/ui/components/borderbox"
	"github.com/m-oehme/jiji/internal/ui/pages/issuelist/entry"
	"github.com/mattn/go-runewidth"
)

// Model represents the issue list page.
type Model struct {
	ctx    *common.Context
	common *common.Common
	issues []jira.Issue
	rows   []Row
	cursor int
	offset int // first visible row index for scrolling
	width  int
	height int
}

// New creates a new issue list page.
func New(ctx *common.Context, c *common.Common) Model {
	return Model{
		ctx:    ctx,
		common: c,
	}
}

// SetItems replaces the issue list.
func (m *Model) SetItems(issues []jira.Issue, sections config.SectionsConfig) {
	m.issues = issues
	m.rows = BuildRows(issues, sections)
	if m.cursor >= len(issues) {
		m.cursor = max(0, len(issues)-1)
	}
	m.offset = 0
}

func (m *Model) Items() []jira.Issue {
	return m.issues
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

func (m *Model) Offset() int {
	return m.offset
}

// Restore the cursor and offset
func (m *Model) Restore(cursor, offset int) {
	if len(m.issues) == 0 {
		m.cursor = 0
		m.offset = 0
		return
	}
	m.cursor = min(cursor, len(m.issues)-1)
	m.offset = min(offset, max(0, len(m.issues)-1))
}

// MoveUp moves the cursor up by one.
func (m *Model) MoveUp() {
	for prev := m.cursor - 1; prev >= 0; prev-- {
		if !m.isHeader(prev) {
			m.cursor = prev
			return
		}
	}
}

// MoveDown moves the cursor down by one.
func (m *Model) MoveDown() {
	for next := m.cursor + 1; next < len(m.rows); next++ {
		if !m.isHeader(next) {
			m.cursor = next
			return
		}
	}
}

// JumpToTop moves the cursor to the first issue.
func (m *Model) JumpToTop() {
	m.cursor = m.firstIssueIndex()
	m.offset = 0
}

// JumpToBottom moves the cursor to the last issue.
func (m *Model) JumpToBottom() {
	m.cursor = m.lastIssueIndex()
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

func (m *Model) isHeader(i int) bool {
	return i >= 0 && i < len(m.rows) && m.rows[i].Header != ""
}

func (m *Model) firstIssueIndex() int {
	for i, r := range m.rows {
		if r.Issue != nil {
			return i
		}
	}
	return 0
}

func (m *Model) lastIssueIndex() int {
	for i := len(m.rows) - 1; i >= 0; i-- {
		if m.rows[i].Issue != nil {
			return i
		}
	}
	return 0
}

// View renders the issue list pane.
func (m Model) View() string {
	if m.width <= 0 || m.height <= 0 {
		return ""
	}

	border := borderbox.New(m.ctx, m.common, m.common.Focused)
	border.SetSize(m.width, m.height)
	contentW, contentH := border.GetContentSize()

	rowSpace := contentH
	if rowSpace < 0 {
		rowSpace = 0
	}

	// Reusable entry model for rendering rows.
	e := entry.New(m.ctx, m.common)
	e.SetSize(contentW)

	// Compute visible window
	start, end := m.visibleRange(rowSpace)
	var rows []string
	for i := start; i < end; i++ {
		row := m.rows[i]
		if row.Issue == nil {
			rows = append(rows, m.common.Styles.Dimmed.Render(
				renderSectionHeader(row.Header, contentW),
			))
		} else {
			e.SetIssue(*row.Issue)
			e.SetSelected(i == m.cursor)
			rows = append(rows, e.View())
		}
	}

	// Pad remaining lines
	for len(rows) < rowSpace {
		rows = append(rows, strings.Repeat(" ", contentW))
	}

	content := strings.Join(rows, "\n")

	return border.Render(content, "Issues")
}

func renderSectionHeader(title string, width int) string {
	prefix := "─── " + title + " "
	remaining := width - runewidth.StringWidth(prefix)
	if remaining > 0 {
		return prefix + strings.Repeat("─", remaining)
	}
	return common.Truncate(prefix, width)
}

// visibleRange returns the slice of issues to render, keeping the cursor visible.
func (m Model) visibleRange(viewportH int) (start, end int) {
	total := len(m.rows)
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
