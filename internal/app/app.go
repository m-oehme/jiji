// Package app contains the root Bubbletea model that orchestrates all UI components.
package app

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"

	"github.com/m-oehme/jiji/internal/config"
	"github.com/m-oehme/jiji/internal/jira"
	"github.com/m-oehme/jiji/internal/ui/common"
	"github.com/m-oehme/jiji/internal/ui/components/help"
	"github.com/m-oehme/jiji/internal/ui/components/statusbar"
	"github.com/m-oehme/jiji/internal/ui/components/tabs"
	"github.com/m-oehme/jiji/internal/ui/pages/detail"
	"github.com/m-oehme/jiji/internal/ui/pages/issuelist/entry"
	"github.com/m-oehme/jiji/internal/ui/pages/issuepane"
	"github.com/m-oehme/jiji/internal/ui/styles"
)

// errorDisplayDuration is how long error messages show before auto-clearing.
const errorDisplayDuration = 5 * time.Second

// Model is the root tea.Model composing all UI elements (ADR-002, ADR-009, ADR-010).
type Model struct {
	cfg    *config.Config
	client jira.Client
	log    *slog.Logger

	focus  *common.Focus
	styles *styles.Styles

	// Shared common state for left and right panes
	listCommon   *common.Common
	detailCommon *common.Common

	// Components
	tabs      tabs.Model
	statusBar statusbar.Model
	help      help.Model

	// Pages
	issuepane issuepane.Model
	// issueList issuelist.Model
	detail detail.Model

	// Dimensions
	width, height int
}

// New creates the root app model.
func New(cfg *config.Config, client jira.Client, log *slog.Logger) Model {
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
		log:          log,
		focus:        f,
		styles:       s,
		listCommon:   listCommon,
		detailCommon: detailCommon,
		tabs:         tabs.New(cfg.Tabs, s),
		statusBar:    statusbar.New(s),
		help:         help.New(&cfg.Keys, s),
		issuepane:    issuepane.New(listCommon, entry.ColumnsFromConfig(cfg.UI.Fields.List)),
		detail:       detail.New(detailCommon),
	}

	m.issuepane.JqlSearch.SetJQL(cfg.Tabs[0].JQL)

	return m
}

// Init returns initial commands — request terminal size and fire first search.
func (m Model) Init() tea.Cmd {
	m.log.Info("starting jiji", "tabs", len(m.cfg.Tabs))
	m.statusBar.SetLoading(true)
	cmds := []tea.Cmd{
		func() tea.Msg { return tea.RequestWindowSize() },
	}
	if m.client != nil {
		cmds = append(cmds, m.searchIssues(m.cfg.Tabs[0].JQL, 0))
	}
	return tea.Batch(cmds...)
}

// Update handles messages per the routing priority from ADR-009 and async results (ADR-010).
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.log.Debug("msg received", "type", fmt.Sprintf("%T", msg))

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
		m.log.Debug("window resize", "width", msg.Width, "height", msg.Height)
		m.width = msg.Width
		m.height = msg.Height
		m.recalcLayout()
		return m, nil

	// Async API results (ADR-010)
	case SearchResultMsg:
		m.log.Info("search results", "count", len(msg.Issues), "tab", msg.TabIndex)
		m.issuepane.IssueList.SetItems(msg.Issues)
		m.statusBar.SetLoading(false)
		m.statusBar.SetIssueCount(len(msg.Issues))
		// Auto-select first issue and load detail
		if len(msg.Issues) > 0 {
			m.issuepane.IssueList.JumpToTop()
			issue := m.issuepane.IssueList.SelectedIssue()
			m.statusBar.SetCurrentIssue(issue.Key)
			return m, tea.Batch(
				m.loadIssueDetail(issue.Key),
				m.loadComments(issue.Key),
			)
		}
		return m, nil

	case IssueDetailMsg:
		m.log.Info("issue detail loaded", "key", msg.Issue.Key)
		m.detail.SetIssue(msg.Issue)
		m.statusBar.SetLoading(false)
		return m, nil

	case CommentsMsg:
		m.log.Info("comments loaded", "key", msg.IssueKey, "count", len(msg.Comments))
		m.detail.SetComments(msg.Comments)
		return m, nil

	case ErrorMsg:
		m.log.Error("api error", "context", msg.Context, "err", msg.Err)
		m.statusBar.SetError(fmt.Errorf("%s: %w", msg.Context, msg.Err))
		m.statusBar.SetLoading(false)
		return m, clearErrorAfter(errorDisplayDuration)

	case clearErrorMsg:
		m.statusBar.ClearError()
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
	leftPane := m.issuepane.View()
	rightPane := m.detail.View()
	body := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)
	statusBar := m.statusBar.View()

	// Clip body to its allocated height. Panes use MaxHeight internally,
	// but JoinHorizontal can pad to the tallest pane.
	bodyH := m.height - tabBarHeight - statusBarHeight
	if bodyLines := strings.Split(body, "\n"); len(bodyLines) > bodyH {
		body = strings.Join(bodyLines[:bodyH], "\n")
	}

	content := lipgloss.JoinVertical(lipgloss.Left, tabBar, body, statusBar)

	// Overlay on top if active
	if m.focus.HasOverlay() && m.focus.TopOverlay() == common.OverlayHelp {
		content = m.help.View()
	}

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

// Fixed heights for tab bar and status bar.
const (
	tabBarHeight    = 1
	statusBarHeight = 1
)

// recalcLayout distributes available space to all components.
func (m *Model) recalcLayout() {
	// Body gets everything except the fixed tab bar and status bar.
	bodyH := m.height - tabBarHeight - statusBarHeight
	if bodyH < 1 {
		bodyH = 1
	}

	m.tabs.SetSize(m.width, tabBarHeight)
	m.statusBar.SetSize(m.width, statusBarHeight)
	m.help.SetSize(m.width, m.height)

	// Split body into left (issue list) and right (detail)
	leftW, rightW := common.SplitHorizontal(m.width, m.cfg.UI.ListRatio)
	// m.issueList.SetSize(leftW, bodyH)
	m.issuepane.SetSize(leftW, bodyH)
	m.detail.SetSize(rightW, bodyH)
}

// syncFocus updates the Focused field on both pane commons.
func (m *Model) syncFocus() {
	m.listCommon.Focused = m.focus.ActivePane() == common.PaneIssueList
	m.detailCommon.Focused = m.focus.ActivePane() == common.PaneDetail
}
