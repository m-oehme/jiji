package app

import (
	"slices"

	tea "charm.land/bubbletea/v2"

	"github.com/m-oehme/jiji/internal/ui/common"
)

// handleKey routes key presses per ADR-009 priority.
func (m Model) handleKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

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
		return m.handleIssueListKey(key)
	case common.PaneDetail:
		return m.handleDetailKey(key)
	}

	return m, nil
}

// handleOverlayKey handles keys when an overlay is active.
func (m Model) handleOverlayKey(key string) (tea.Model, tea.Cmd) {
	switch {
	case matchKey(key, m.cfg.Keys.Cancel) || key == "esc":
		m.focus.PopOverlay()
	case matchKey(key, m.cfg.Keys.Help):
		if m.focus.TopOverlay() == common.OverlayHelp {
			m.focus.PopOverlay()
		}
	case matchKey(key, m.cfg.Keys.Down):
		m.help.ScrollDown()
	case matchKey(key, m.cfg.Keys.Up):
		m.help.ScrollUp()
	}
	return m, nil
}

// handleGlobalKey handles keys that work regardless of active pane.
// Returns handled=true if the key was consumed.
func (m Model) handleGlobalKey(key string) (tea.Model, tea.Cmd, bool) {
	switch {
	case matchKey(key, m.cfg.Keys.Help):
		m.focus.PushOverlay(common.OverlayHelp)
		return m, nil, true

	case matchKey(key, m.cfg.Keys.Quit):
		return m, tea.Quit, true

	case matchKey(key, m.cfg.Keys.TabNext):
		idx := m.tabs.Active() + 1
		if idx >= m.tabs.Count() {
			idx = 0
		}
		m.tabs.SetActive(idx)
		return m, nil, true

	case matchKey(key, m.cfg.Keys.TabPrev):
		idx := m.tabs.Active() - 1
		if idx < 0 {
			idx = m.tabs.Count() - 1
		}
		m.tabs.SetActive(idx)
		return m, nil, true

	case matchKey(key, m.cfg.Keys.PaneSwitch):
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

// handleIssueListKey handles keys when the issue list pane is active.
func (m Model) handleIssueListKey(key string) (tea.Model, tea.Cmd) {
	prevIdx := m.issueList.SelectedIndex()

	switch {
	case matchKey(key, m.cfg.Keys.Down):
		m.issueList.MoveDown()
	case matchKey(key, m.cfg.Keys.Up):
		m.issueList.MoveUp()
	case matchKey(key, m.cfg.Keys.Top):
		m.issueList.JumpToTop()
	case matchKey(key, m.cfg.Keys.Bottom):
		m.issueList.JumpToBottom()
	}

	// Update detail if selection changed
	if m.issueList.SelectedIndex() != prevIdx {
		if sel := m.issueList.SelectedIssue(); sel != nil {
			m.detail.SetIssue(sel)
			m.detail.SetComments(mockComments())
			m.statusBar.SetCurrentIssue(sel.Key)
		}
	}

	return m, nil
}

// handleDetailKey handles keys when the detail pane is active.
func (m Model) handleDetailKey(key string) (tea.Model, tea.Cmd) {
	switch {
	case matchKey(key, m.cfg.Keys.Down):
		m.detail.ScrollDown()
	case matchKey(key, m.cfg.Keys.Up):
		m.detail.ScrollUp()
	case matchKey(key, m.cfg.Keys.Top):
		m.detail.ScrollToTop()
	case matchKey(key, m.cfg.Keys.Bottom):
		m.detail.ScrollToBottom()
	}
	return m, nil
}

// matchKey checks if a key string matches any of the configured bindings.
func matchKey(key string, bindings []string) bool {
	return slices.Contains(bindings, key)
}
