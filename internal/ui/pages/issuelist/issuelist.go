// Package issuelist renders the left pane: JQL input + issue table (ADR-005).
package issuelist

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/textinput"
	lipgloss "charm.land/lipgloss/v2"

	"github.com/m-oehme/jiji/internal/jira"
	"github.com/m-oehme/jiji/internal/ui/common"
)

// Column widths for the issue table.
const (
	colKeyWidth      = 12
	colPriorityWidth = 3
	colAssigneeWidth = 15
)

// Model represents the issue list page.
type Model struct {
	common   *common.Common
	issues   []jira.Issue
	cursor   int
	offset   int // first visible row index for scrolling
	jql      string // last submitted JQL
	jqlInput textinput.Model
	jqlFocus bool // true when JQL input is being edited
	width    int
	height   int
}

// New creates a new issue list page.
func New(c *common.Common) Model {
	ti := textinput.New()
	ti.Prompt = "JQL: "
	ti.SetVirtualCursor(true)
	return Model{
		common:   c,
		jqlInput: ti,
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

// SetJQL sets the JQL query string displayed above the table.
func (m *Model) SetJQL(jql string) {
	m.jql = jql
	m.jqlInput.SetValue(jql)
}

// JQLValue returns the current JQL string (from the text input when focused, or last submitted).
func (m *Model) JQLValue() string {
	if m.jqlFocus {
		return m.jqlInput.Value()
	}
	return m.jql
}

// FocusJQL activates the JQL text input for editing.
func (m *Model) FocusJQL() tea.Cmd {
	m.jqlFocus = true
	m.jqlInput.SetValue(m.jql)
	m.jqlInput.CursorEnd()
	return m.jqlInput.Focus()
}

// UnfocusJQL deactivates the JQL text input, reverting to the last submitted value.
func (m *Model) UnfocusJQL() {
	m.jqlFocus = false
	m.jqlInput.Blur()
	// On cancel, revert to last submitted JQL
	m.jqlInput.SetValue(m.jql)
}

// IsJQLFocused returns whether the JQL input is being edited.
func (m *Model) IsJQLFocused() bool {
	return m.jqlFocus
}

// UpdateJQL delegates a key message to the textinput and returns any command.
func (m *Model) UpdateJQL(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.jqlInput, cmd = m.jqlInput.Update(msg)
	return cmd
}

// SetFocused updates the focused state.
func (m *Model) SetFocused(focused bool) {
	m.common.Focused = focused
}

// SetSize updates the available dimensions.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	frameW, _ := m.common.Styles.Border.GetFrameSize()
	contentW := width - frameW
	// Account for the "JQL: " prompt width
	inputW := contentW - lipgloss.Width(m.jqlInput.Prompt)
	if inputW < 1 {
		inputW = 1
	}
	m.jqlInput.SetWidth(inputW)
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

	// JQL line: interactive textinput when focused, static text when unfocused.
	var jqlLine string
	if m.jqlFocus {
		jqlLine = m.jqlInput.View()
	} else {
		jqlLine = m.common.Styles.Dimmed.Render(truncate("JQL: "+m.jql, contentW))
	}

	// Header
	header := m.renderHeader(contentW)

	// Row space = content height minus JQL and header lines
	rowSpace := contentH - lipgloss.Height(jqlLine) - 1
	if rowSpace < 0 {
		rowSpace = 0
	}

	// Compute visible window
	start, end := m.visibleRange(rowSpace)
	var rows []string
	for i := start; i < end; i++ {
		rows = append(rows, m.renderRow(i, contentW))
	}

	// Pad remaining lines
	for len(rows) < rowSpace {
		rows = append(rows, strings.Repeat(" ", contentW))
	}

	content := jqlLine + "\n" + header + "\n" + strings.Join(rows, "\n")

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

// prioritySymbol maps Jira priority names to compact single-width symbols.
func prioritySymbol(name string) string {
	switch strings.ToLower(name) {
	case "highest":
		return "↑↑"
	case "high":
		return "↑"
	case "medium":
		return "●"
	case "low":
		return "↓"
	case "lowest":
		return "↓↓"
	default:
		return "·"
	}
}

// renderHeader renders the column header row.
func (m Model) renderHeader(width int) string {
	summaryW := width - colKeyWidth - colPriorityWidth - colAssigneeWidth
	if summaryW < 4 {
		summaryW = 4
	}
	header := fmt.Sprintf("%-*s%-*s%-*s%-*s",
		colKeyWidth, "KEY",
		colPriorityWidth, "P",
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

	pri := prioritySymbol(issue.Priority)
	row := fmt.Sprintf("%-*s%-*s%-*s%-*s",
		colKeyWidth, truncate(issue.Key, colKeyWidth),
		colPriorityWidth, pri,
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
