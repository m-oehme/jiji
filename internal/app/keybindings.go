package app

import (
	"os"
	"os/exec"
	"slices"

	tea "charm.land/bubbletea/v2"

	"github.com/m-oehme/jiji/internal/ui/common"
)

// handleKey routes key presses per ADR-009 priority.
func (m Model) handleKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	key := msg.String()
	m.ctx.Logger.Debug("key press", "key", key, "pane", m.focus.ActivePane(), "overlay", m.focus.HasOverlay(), "jql_focused", m.issuepane.JqlSearch.IsJQLFocused())

	if model, cmd, handled := m.handleUserKeys(key); handled {
		return model, cmd
	}

	// JQL input focused — route all keys to text input except Enter/Esc
	if m.issuepane.JqlSearch.IsJQLFocused() {
		return m.handleJQLKey(msg)
	}

	// 3. Overlay active — route to overlay
	if m.focus.HasOverlay() {
		return m.handleOverlayKey(key)
	}

	// 4. Global keys
	if model, cmd, handled := m.handleGlobalKey(key); handled {
		return model, cmd
	}

	// 5. Delegate to active pane
	switch m.focus.ActivePane() {
	case common.PaneIssueList:
		return m.handleIssueListKey(key, msg)
	case common.PaneDetail:
		return m.handleDetailKey(key)
	}

	return m, nil
}

// handleOverlayKey handles keys when an overlay is active.
func (m Model) handleOverlayKey(key string) (tea.Model, tea.Cmd) {
	switch {
	case matchKey(key, m.ctx.Config.Keys.Builtin.Cancel) || key == "esc":
		m.focus.PopOverlay()
	case matchKey(key, m.ctx.Config.Keys.Builtin.Help):
		if m.focus.TopOverlay() == common.OverlayHelp {
			m.focus.PopOverlay()
		}
	case matchKey(key, m.ctx.Config.Keys.Builtin.Down):
		m.help.ScrollDown()
	case matchKey(key, m.ctx.Config.Keys.Builtin.Up):
		m.help.ScrollUp()
	}
	return m, nil
}

// handleGlobalKey handles keys that work regardless of active pane.
// Returns handled=true if the key was consumed.
func (m Model) handleGlobalKey(key string) (tea.Model, tea.Cmd, bool) {
	switch {
	case matchKey(key, m.ctx.Config.Keys.Builtin.Help):
		m.focus.PushOverlay(common.OverlayHelp)
		return m, nil, true

	case matchKey(key, m.ctx.Config.Keys.Builtin.Quit):
		return m, tea.Quit, true

	case matchKey(key, m.ctx.Config.Keys.Builtin.TabNext):
		idx := m.tabs.Active() + 1
		if idx >= m.tabs.Count() {
			idx = 0
		}
		m.tabs.SetActive(idx)
		return m, nil, true

	case matchKey(key, m.ctx.Config.Keys.Builtin.TabPrev):
		idx := m.tabs.Active() - 1
		if idx < 0 {
			idx = m.tabs.Count() - 1
		}
		m.tabs.SetActive(idx)
		return m, nil, true

	case matchKey(key, m.ctx.Config.Keys.Builtin.PaneSwitch):
		m.focus.TogglePane()
		m.syncFocus()
		return m, nil, true

	// Number keys 1-9 jump to tab
	case key >= "1" && key <= "9":
		idx := int(key[0]-'0') - 1
		if idx < m.tabs.Count() {
			m.tabs.SetActive(idx)
		}
		return m, nil, true
	}

	return m, nil, false
}

// handleJQLKey handles keys when the JQL input is focused.
// Enter submits, Esc cancels, everything else routes to the text input.
func (m Model) handleJQLKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	switch {
	case matchKey(key, m.ctx.Config.Keys.Builtin.Confirm) || key == "enter":
		jql := m.issuepane.JqlSearch.JQLValue()
		m.issuepane.JqlSearch.UnfocusJQL()
		m.statusBar.SetLoading(true)
		return m, m.searchIssues(jql, m.tabs.Active())

	case matchKey(key, m.ctx.Config.Keys.Builtin.Cancel) || key == "esc":
		m.issuepane.JqlSearch.UnfocusJQL()
		return m, nil
	}

	// Delegate to textinput
	cmd := m.issuepane.JqlSearch.UpdateJQL(msg)
	return m, cmd
}

// handleIssueListKey handles keys when the issue list pane is active.
func (m Model) handleIssueListKey(key string, msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	// Focus JQL input
	if matchKey(key, m.ctx.Config.Keys.Builtin.FocusJQL) {
		cmd := m.issuepane.JqlSearch.FocusJQL()
		return m, cmd
	}

	prevIdx := m.issuepane.IssueList.SelectedIndex()

	switch {
	case matchKey(key, m.ctx.Config.Keys.Builtin.Down):
		m.issuepane.IssueList.MoveDown()
	case matchKey(key, m.ctx.Config.Keys.Builtin.Up):
		m.issuepane.IssueList.MoveUp()
	case matchKey(key, m.ctx.Config.Keys.Builtin.Top):
		m.issuepane.IssueList.JumpToTop()
	case matchKey(key, m.ctx.Config.Keys.Builtin.Bottom):
		m.issuepane.IssueList.JumpToBottom()
	}

	// Auto-load detail if selection changed (ADR-005)
	if m.issuepane.IssueList.SelectedIndex() != prevIdx {
		if sel := m.issuepane.IssueList.SelectedIssue(); sel != nil {
			m.statusBar.SetCurrentIssue(sel.Key)
			if m.client != nil {
				return m, tea.Batch(
					m.loadIssueDetail(sel.Key),
					m.loadComments(sel.Key),
				)
			}
		}
	}

	return m, nil
}

// handleDetailKey handles keys when the detail pane is active.
func (m Model) handleDetailKey(key string) (tea.Model, tea.Cmd) {
	switch {
	case matchKey(key, m.ctx.Config.Keys.Builtin.Down):
		m.detail.ScrollDown()
	case matchKey(key, m.ctx.Config.Keys.Builtin.Up):
		m.detail.ScrollUp()
	case matchKey(key, m.ctx.Config.Keys.Builtin.Top):
		m.detail.ScrollToTop()
	case matchKey(key, m.ctx.Config.Keys.Builtin.Bottom):
		m.detail.ScrollToBottom()
	}
	return m, nil
}

// matchKey checks if a key string matches any of the configured bindings.
func matchKey(key string, bindings []string) bool {
	return slices.Contains(bindings, key)
}

func (m *Model) handleUserKeys(key string) (tea.Model, tea.Cmd, bool) {
	issueKeyConfig := m.ctx.Config.Keys.User.Issues
	if len(issueKeyConfig) == 0 {
		return m, nil, false
	}
	for _, issue := range issueKeyConfig {
		if key == issue.Key {
			command, err := m.issuepane.IssueList.SelectedIssue().Format(issue.Command)
			if err != nil {
				return m, func() tea.Msg {
					return ErrorMsg{
						Err:     err,
						Context: "error parsing command",
					}
				}, true
			}
			cmd := m.executeShellCommand(command)
			return m, cmd, true
		}
	}

	globalConfig := m.ctx.Config.Keys.User.Global
	if len(globalConfig) == 0 {
		return m, nil, false
	}

	for _, gKey := range globalConfig {
		if key == gKey.Key {
			cmd := m.executeShellCommand(gKey.Command)
			return m, cmd, true
		}
	}

	return m, nil, false
}

func (m *Model) executeShellCommand(cmd string) tea.Cmd {
	m.ctx.Logger.Debug("execute shell command", "command", cmd)
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "sh"
	}
	c := exec.Command(shell, "-c", cmd)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		if err != nil {
			return ErrorMsg{
				Err:     err,
				Context: "Error executing command",
			}
		}
		return nil
	})
}
