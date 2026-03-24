package app

import (
	"context"

	tea "charm.land/bubbletea/v2"
)

// tea.Cmd functions for async Jira API calls (ADR-010).

func (m *Model) searchIssues(jql string, tabIdx int) tea.Cmd {
	m.ctx.Logger.Info("searching issues", "jql", jql, "tab", tabIdx)
	return func() tea.Msg {
		result, err := m.client.SearchIssues(
			context.Background(), jql, m.ctx.Config.UI.Fields.List, "",
		)
		if err != nil {
			m.ctx.Logger.Error("search failed", "jql", jql, "err", err)
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
	m.ctx.Logger.Info("loading issue detail", "key", key)
	return func() tea.Msg {
		issue, err := m.client.GetIssue(context.Background(), key)
		if err != nil {
			m.ctx.Logger.Error("detail load failed", "key", key, "err", err)
			return ErrorMsg{Err: err, Context: "detail"}
		}
		return IssueDetailMsg{Issue: issue}
	}
}

func (m *Model) loadComments(key string) tea.Cmd {
	m.ctx.Logger.Info("loading comments", "key", key)
	return func() tea.Msg {
		comments, err := m.client.GetComments(context.Background(), key)
		if err != nil {
			m.ctx.Logger.Error("comments load failed", "key", key, "err", err)
			return ErrorMsg{Err: err, Context: "comments"}
		}
		return CommentsMsg{Comments: comments, IssueKey: key}
	}
}
