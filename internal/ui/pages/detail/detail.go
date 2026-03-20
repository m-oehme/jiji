package detail

import (
	"strings"

	"github.com/m-oehme/jiji/internal/jira"
	"github.com/m-oehme/jiji/internal/ui/common"
	"charm.land/bubbles/v2/viewport"
)

// Model is the detail pane orchestrator: metadata banner + scrollable content.
type Model struct {
	common   *common.Common
	issue    *jira.Issue
	comments []jira.Comment
	viewport viewport.Model
	width    int
	height   int
	ready    bool
}

// New creates a new detail pane.
func New(c *common.Common) Model {
	return Model{
		common: c,
	}
}

// SetIssue updates the displayed issue and rebuilds content.
func (m *Model) SetIssue(issue *jira.Issue) {
	m.issue = issue
	m.rebuildContent()
}

// SetComments updates the displayed comments and rebuilds content.
func (m *Model) SetComments(comments []jira.Comment) {
	m.comments = comments
	m.rebuildContent()
}

// CurrentIssue returns the currently displayed issue.
func (m *Model) CurrentIssue() *jira.Issue {
	return m.issue
}

// ScrollUp scrolls the content viewport up.
func (m *Model) ScrollUp() {
	m.viewport.ScrollUp(1)
}

// ScrollDown scrolls the content viewport down.
func (m *Model) ScrollDown() {
	m.viewport.ScrollDown(1)
}

// ScrollToTop jumps to the top of the content.
func (m *Model) ScrollToTop() {
	m.viewport.GotoTop()
}

// ScrollToBottom jumps to the bottom of the content.
func (m *Model) ScrollToBottom() {
	m.viewport.GotoBottom()
}

// SetFocused updates the focused state.
func (m *Model) SetFocused(focused bool) {
	m.common.Focused = focused
}

// SetSize updates the available dimensions and reconfigures the viewport.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height

	frameW, frameH := m.common.Styles.Border.GetFrameSize()
	contentW := width - frameW
	contentH := height - frameH
	vpH := contentH - metadataBannerHeight
	if vpH < 1 {
		vpH = 1
	}

	if !m.ready {
		m.viewport = viewport.New(
			viewport.WithWidth(contentW),
			viewport.WithHeight(vpH),
		)
		m.ready = true
	} else {
		m.viewport.SetWidth(contentW)
		m.viewport.SetHeight(vpH)
	}

	m.rebuildContent()
}

// rebuildContent regenerates the viewport content from the current issue/comments.
func (m *Model) rebuildContent() {
	if !m.ready {
		return
	}
	content := buildContent(m.issue, m.comments, m.common.Styles)
	m.viewport.SetContent(content)
}

// View renders the detail pane: metadata + scrollable content inside a border.
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

	meta := renderMetadata(m.issue, contentW, m.common.Styles)
	vpContent := m.viewport.View()

	// Stack metadata + viewport content
	body := meta + "\n" + vpContent

	// Pad/truncate to fit content dimensions
	lines := strings.Split(body, "\n")
	for len(lines) < contentH {
		lines = append(lines, strings.Repeat(" ", contentW))
	}
	if len(lines) > contentH {
		lines = lines[:contentH]
	}

	return borderStyle.
		Width(m.width).
		Height(m.height).
		MaxWidth(m.width).
		MaxHeight(m.height).
		Render(strings.Join(lines, "\n"))
}
