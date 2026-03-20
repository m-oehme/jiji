package app

import (
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/m-oehme/jiji/internal/jira"
)

// Async message types for API responses (ADR-010).
// These are the only typed messages — synchronous updates use method calls.

// SearchResultMsg carries issues returned from a JQL search.
type SearchResultMsg struct {
	Issues   []jira.Issue
	NextPage jira.PageToken
	TabIndex int // which tab initiated the search
}

// IssueDetailMsg carries a fully loaded issue.
type IssueDetailMsg struct {
	Issue *jira.Issue
}

// CommentsMsg carries comments for an issue.
type CommentsMsg struct {
	Comments []jira.Comment
	IssueKey string
}

// TransitionsMsg carries available workflow transitions.
type TransitionsMsg struct {
	Transitions []jira.Transition
}

// LabelsListMsg carries all available labels.
type LabelsListMsg struct {
	Labels []string
}

// ErrorMsg carries an API error with context.
type ErrorMsg struct {
	Err     error
	Context string // "search", "detail", "comments", etc.
}

// clearErrorMsg is an internal message to auto-clear errors after a timeout.
type clearErrorMsg struct{}

// clearErrorAfter returns a tea.Cmd that sends clearErrorMsg after the given duration.
func clearErrorAfter(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}
