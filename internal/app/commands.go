package app

import (
	"context"

	tea "charm.land/bubbletea/v2"
)

// tea.Cmd functions for async Jira API calls (ADR-010).

func (m *Model) searchIssues(jql string, tabIdx int) tea.Cmd {
	m.log.Info("searching issues", "jql", jql, "tab", tabIdx)
	return func() tea.Msg {
		result, err := m.client.SearchIssues(
			context.Background(), jql, m.cfg.UI.Fields.List, "",
		)
		if err != nil {
			m.log.Error("search failed", "jql", jql, "err", err)
			return ErrorMsg{Err: err, Context: "search"}
		}
		return SearchResultMsg{
			Issues:   result.Issues,
			NextPage: result.NextPage,
			TabIndex: tabIdx,
		}
	}
}

func (m *Model) loadIssueDetail(key string) tea.Cmd {
	m.log.Info("loading issue detail", "key", key)
	return func() tea.Msg {
		issue, err := m.client.GetIssue(context.Background(), key)
		if err != nil {
			m.log.Error("detail load failed", "key", key, "err", err)
			return ErrorMsg{Err: err, Context: "detail"}
		}
		return IssueDetailMsg{Issue: issue}
	}
}

func (m *Model) loadComments(key string) tea.Cmd {
	m.log.Info("loading comments", "key", key)
	return func() tea.Msg {
		comments, err := m.client.GetComments(context.Background(), key)
		if err != nil {
			m.log.Error("comments load failed", "key", key, "err", err)
			return ErrorMsg{Err: err, Context: "comments"}
		}
		return CommentsMsg{Comments: comments, IssueKey: key}
	}
}
