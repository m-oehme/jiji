// Package app contains the root Bubbletea model that orchestrates all UI components.
package app

import (
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"

	"github.com/m-oehme/jiji/internal/config"
	"github.com/m-oehme/jiji/internal/jira"
	"github.com/m-oehme/jiji/internal/ui/common"
	"github.com/m-oehme/jiji/internal/ui/components/help"
	"github.com/m-oehme/jiji/internal/ui/components/statusbar"
	"github.com/m-oehme/jiji/internal/ui/components/tabs"
	"github.com/m-oehme/jiji/internal/ui/pages/detail"
	"github.com/m-oehme/jiji/internal/ui/pages/issuelist"
	"github.com/m-oehme/jiji/internal/ui/styles"
)

// Model is the root tea.Model composing all UI elements (ADR-002, ADR-009, ADR-010).
type Model struct {
	cfg    *config.Config
	client jira.Client // not used in Phase 2

	focus  *common.Focus
	styles *styles.Styles

	// Shared common state for left and right panes
	listCommon   *common.Common
	detailCommon *common.Common

	// Components
	tabs      tabs.Model
	statusBar statusbar.Model
	help      help.Model

	// Pages (Phase 2: single tab)
	issueList issuelist.Model
	detail    detail.Model

	// Dimensions
	width, height int
}

// New creates the root app model.
func New(cfg *config.Config, client jira.Client) Model {
	s := styles.NewStyles(cfg.Theme)
	f := common.NewFocus()

	listCommon := &common.Common{
		Styles:  s,
		Keys:    &cfg.Keys,
		Focused: true, // issue list starts focused
	}
	detailCommon := &common.Common{
		Styles:  s,
		Keys:    &cfg.Keys,
		Focused: false,
	}

	m := Model{
		cfg:          cfg,
		client:       client,
		focus:        f,
		styles:       s,
		listCommon:   listCommon,
		detailCommon: detailCommon,
		tabs:         tabs.New(cfg.Tabs, s),
		statusBar:    statusbar.New(s),
		help:         help.New(&cfg.Keys, s),
		issueList:    issuelist.New(listCommon),
		detail:       detail.New(detailCommon),
	}

	// Load mock data
	issues := mockIssues()
	m.issueList.SetItems(issues)
	m.issueList.SetJQL(cfg.Tabs[0].JQL)
	m.statusBar.SetIssueCount(len(issues))

	// Show first issue in detail
	if sel := m.issueList.SelectedIssue(); sel != nil {
		m.detail.SetIssue(sel)
		m.detail.SetComments(mockComments())
		m.statusBar.SetCurrentIssue(sel.Key)
	}

	return m
}

// Init returns the initial command — request terminal size.
func (m Model) Init() tea.Cmd {
	return func() tea.Msg { return tea.RequestWindowSize() }
}

// Update handles messages per the routing priority from ADR-009.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// 1. Ctrl+C — always quit
	case tea.KeyPressMsg:
		if msg.Code == 'c' && msg.Mod == tea.ModCtrl {
			return m, tea.Quit
		}
	}

	switch msg := msg.(type) {
	// 2. Window resize
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.recalcLayout()
		return m, nil

	// 3-5. Key routing
	case tea.KeyPressMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

// View composes the full layout.
func (m Model) View() tea.View {
	if m.width == 0 || m.height == 0 {
		v := tea.NewView("Loading…")
		v.AltScreen = true
		return v
	}

	tabBar := m.tabs.View()
	leftPane := m.issueList.View()
	rightPane := m.detail.View()
	body := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)
	statusBar := m.statusBar.View()

	content := lipgloss.JoinVertical(lipgloss.Left, tabBar, body, statusBar)

	// Overlay on top if active
	if m.focus.HasOverlay() && m.focus.TopOverlay() == common.OverlayHelp {
		content = m.help.View()
	}

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

// recalcLayout distributes available space to all components.
func (m *Model) recalcLayout() {
	// Tab bar: 1 line at top
	tabH := 1
	// Status bar: 1 line at bottom
	statusH := 1
	// Body gets the rest
	bodyH := m.height - tabH - statusH
	if bodyH < 1 {
		bodyH = 1
	}

	m.tabs.SetSize(m.width, tabH)
	m.statusBar.SetSize(m.width, statusH)
	m.help.SetSize(m.width, m.height)

	// Split body into left (issue list) and right (detail)
	leftW, rightW := common.SplitHorizontal(m.width, m.cfg.UI.ListRatio)
	m.issueList.SetSize(leftW, bodyH)
	m.detail.SetSize(rightW, bodyH)
}

// syncFocus updates the Focused field on both pane commons.
func (m *Model) syncFocus() {
	m.listCommon.Focused = m.focus.ActivePane() == common.PaneIssueList
	m.detailCommon.Focused = m.focus.ActivePane() == common.PaneDetail
}
