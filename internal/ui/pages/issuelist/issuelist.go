// Package issuelist renders the left pane: JQL input + issue table (ADR-005).
package issuelist

import (
	"fmt"
	"strings"

	"github.com/m-oehme/jiji/internal/jira"
	"github.com/m-oehme/jiji/internal/ui/common"
	lipgloss "charm.land/lipgloss/v2"
)

// Column widths for the issue table.
const (
	colKeyWidth      = 12
	colPriorityWidth = 8
	colAssigneeWidth = 15
)

// Model represents the issue list page.
type Model struct {
	common   *common.Common
	issues   []jira.Issue
	cursor   int
	jql      string
	width    int
	height   int
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
}

// JumpToBottom moves the cursor to the last issue.
func (m *Model) JumpToBottom() {
	if len(m.issues) > 0 {
		m.cursor = len(m.issues) - 1
	}
}

// SetJQL sets the JQL query string displayed above the table.
func (m *Model) SetJQL(jql string) {
	m.jql = jql
}

// JQLValue returns the current JQL string.
func (m *Model) JQLValue() string {
	return m.jql
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

	innerW, innerH := common.InnerSize(m.width, m.height, true)
	if innerW <= 0 || innerH <= 0 {
		return borderStyle.Width(m.width - 2).Height(m.height - 2).Render("")
	}

	// JQL line (static for Phase 2)
	jqlLine := m.common.Styles.Dimmed.Width(innerW).Render(
		truncate("JQL: "+m.jql, innerW),
	)

	// Header
	header := m.renderHeader(innerW)

	// Remaining space for rows
	rowSpace := innerH - 2 // jql + header
	if rowSpace < 0 {
		rowSpace = 0
	}

	// Compute visible window
	start, end := m.visibleRange(rowSpace)
	var rows []string
	for i := start; i < end; i++ {
		rows = append(rows, m.renderRow(i, innerW))
	}

	// Pad remaining lines
	for len(rows) < rowSpace {
		rows = append(rows, strings.Repeat(" ", innerW))
	}

	content := jqlLine + "\n" + header + "\n" + strings.Join(rows, "\n")
	return borderStyle.Width(innerW).Height(innerH).Render(content)
}

// visibleRange calculates which issues to show given the viewport size.
func (m Model) visibleRange(viewportH int) (start, end int) {
	total := len(m.issues)
	if total == 0 || viewportH <= 0 {
		return 0, 0
	}

	start = 0
	if m.cursor >= viewportH {
		start = m.cursor - viewportH + 1
	}
	end = start + viewportH
	if end > total {
		end = total
	}
	return start, end
}

// renderHeader renders the column header row.
func (m Model) renderHeader(width int) string {
	summaryW := width - colKeyWidth - colPriorityWidth - colAssigneeWidth
	if summaryW < 4 {
		summaryW = 4
	}
	header := fmt.Sprintf("%-*s%-*s%-*s%-*s",
		colKeyWidth, "KEY",
		colPriorityWidth, "PRI",
		colAssigneeWidth, "ASSIGNEE",
		summaryW, "SUMMARY",
	)
	return m.common.Styles.Dimmed.Width(width).Render(truncate(header, width))
}

// renderRow renders a single issue row.
func (m Model) renderRow(idx, width int) string {
	issue := m.issues[idx]
	summaryW := width - colKeyWidth - colPriorityWidth - colAssigneeWidth
	if summaryW < 4 {
		summaryW = 4
	}

	row := fmt.Sprintf("%-*s%-*s%-*s%-*s",
		colKeyWidth, truncate(issue.Key, colKeyWidth),
		colPriorityWidth, truncate(issue.Priority, colPriorityWidth),
		colAssigneeWidth, truncate(issue.Assignee, colAssigneeWidth),
		summaryW, truncate(issue.Summary, summaryW),
	)
	row = truncate(row, width)

	if idx == m.cursor {
		return m.common.Styles.IssueSelected.Width(width).Render(row)
	}
	return lipgloss.NewStyle().Width(width).Render(row)
}

// truncate cuts a string to maxLen, appending "…" if truncated.
func truncate(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 1 {
		return "…"
	}
	return s[:maxLen-1] + "…"
}
