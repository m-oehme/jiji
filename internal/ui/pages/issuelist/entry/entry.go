package entry

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/m-oehme/jiji/internal/jira"
	"github.com/m-oehme/jiji/internal/ui/common"
)

// Column defines a single column in the issue list.
type Column struct {
	ID      string                  // matches config field ID
	Header  string                  // display name (unused for now, kept for future use)
	Width   int                     // fixed width; 0 = flex (fills remaining space)
	Extract func(jira.Issue) string // pulls display value from issue
}

// DefaultColumns maps field IDs to their column definitions.
var DefaultColumns = map[string]Column{
	"key":      {ID: "key", Header: "KEY", Width: 12, Extract: func(i jira.Issue) string { return i.Key }},
	"priority": {ID: "priority", Header: "P", Width: 3, Extract: func(i jira.Issue) string { return prioritySymbol(i.Priority) }},
	"assignee": {ID: "assignee", Header: "ASSIGNEE", Width: 12, Extract: func(i jira.Issue) string { return i.Assignee }},
	"status":   {ID: "status", Header: "STATUS", Width: 12, Extract: func(i jira.Issue) string { return i.Status }},
	"type":     {ID: "type", Header: "TYPE", Width: 10, Extract: func(i jira.Issue) string { return i.Type }},
	"reporter": {ID: "reporter", Header: "REPORTER", Width: 12, Extract: func(i jira.Issue) string { return i.Reporter }},
	"sprint":   {ID: "sprint", Header: "SPRINT", Width: 15, Extract: func(i jira.Issue) string { return i.Sprint }},
	"summary":  {ID: "summary", Header: "SUMMARY", Width: 0, Extract: func(i jira.Issue) string { return i.Summary }},
}

// ColumnsFromConfig resolves config field IDs into an ordered slice of Columns.
// Unknown IDs are skipped. If no flex column exists, the last column becomes flex.
func ColumnsFromConfig(fieldIDs []string) []Column {
	var cols []Column
	for _, id := range fieldIDs {
		if col, ok := DefaultColumns[id]; ok {
			cols = append(cols, col)
		}
	}
	if len(cols) == 0 {
		return nil
	}
	// Ensure at least one flex column.
	hasFlex := false
	for _, c := range cols {
		if c.Width == 0 {
			hasFlex = true
			break
		}
	}
	if !hasFlex {
		cols[len(cols)-1].Width = 0
	}
	return cols
}

// Model renders a single issue row. Create one and reuse it across rows
// by calling SetIssue + SetSelected before each View().
type Model struct {
	common   *common.Common
	columns  []Column
	issue    jira.Issue
	selected bool
	width    int
}

// New creates a new entry model with the given columns.
func New(c *common.Common, columns []Column) Model {
	return Model{
		common:  c,
		columns: columns,
	}
}

// SetIssue sets the issue to render.
func (m *Model) SetIssue(issue jira.Issue) {
	m.issue = issue
}

// SetSelected sets the selected state.
func (m *Model) SetSelected(selected bool) {
	m.selected = selected
}

// SetSize sets the available width.
func (m *Model) SetSize(width int) {
	m.width = width
}

// View renders the issue row.
func (m Model) View() string {
	if m.width <= 0 || len(m.columns) == 0 {
		return ""
	}

	// Calculate flex width: total minus all fixed columns.
	fixedSum := 0
	flexCount := 0
	for _, c := range m.columns {
		if c.Width > 0 {
			fixedSum += c.Width
		} else {
			flexCount++
		}
	}
	flexW := 0
	if flexCount > 0 {
		flexW = (m.width - fixedSum) / flexCount
		if flexW < 1 {
			flexW = 1
		}
	}

	var row strings.Builder
	for _, c := range m.columns {
		w := c.Width
		if w == 0 {
			w = flexW
		}
		val := c.Extract(m.issue)
		// Priority column is already a symbol, don't truncate it.
		if c.ID != "priority" {
			val = common.Truncate(val, w)
		}
		fmt.Fprintf(&row, "%-*s", w, val)
	}

	rendered := common.Truncate(row.String(), m.width)

	if m.selected {
		return m.common.Styles.IssueSelected.Width(m.width).Render(rendered)
	}
	return lipgloss.NewStyle().Width(m.width).Render(rendered)
}

func prioritySymbol(name string) string {
	switch strings.ToLower(name) {
	case "highest":
		return "󰄿"
	case "high":
		return "󰅃"
	case "medium":
		return "-"
	case "low":
		return "󰅀"
	case "lowest":
		return "󰄼"
	default:
		return "·"
	}
}
