package entry

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/m-oehme/jiji/internal/jira"
	"github.com/m-oehme/jiji/internal/ui/common"
)

const (
	lineHeight       = 1
	colKeyWidth      = 12
	colPriorityWidth = 3
	colAssigneeWidth = 12
)

func RenderListEntry(c *common.Common, issue jira.Issue, width int, selected bool) string {
	keyView := common.Truncate(issue.Key, colKeyWidth)
	priorityView := prioritySymbol(issue.Priority)
	assigneeView := common.Truncate(issue.Assignee, colAssigneeWidth)

	summaryW := width - colKeyWidth - colPriorityWidth - colAssigneeWidth
	// summeryView := common.Truncate(issue.Summary, summaryW)

	row := fmt.Sprintf("%-*s%-*s%-*s%-*s",
		colKeyWidth, keyView,
		colPriorityWidth, priorityView,
		colAssigneeWidth, assigneeView,
		summaryW, issue.Summary,
	)
	row = common.Truncate(row, width)

	if selected {
		return c.Styles.IssueSelected.Width(width).Render(row)
	}
	return lipgloss.NewStyle().Width(width).Render(row)
}

// prioritySymbol maps Jira priority names to compact single-width symbols.
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

// renderHeader renders the column header row.
func RenderHeader(c *common.Common, width int) string {
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
	return c.Styles.Dimmed.Width(width).Render(common.Truncate(header, width))
}
