package issuepane

import (
	lipgloss "charm.land/lipgloss/v2"
	"github.com/m-oehme/jiji/internal/ui/common"
	"github.com/m-oehme/jiji/internal/ui/components/jqlsearch"
	"github.com/m-oehme/jiji/internal/ui/pages/issuelist"
)

// Model represents the issue list page.
type Model struct {
	common *common.Common
	width  int
	height int

	JqlSearch jqlsearch.Model
	IssueList issuelist.Model
}

// New creates a new issue list page.
func New(c *common.Common) Model {
	issuelist := issuelist.New(c)
	jqlsearch := jqlsearch.New(c)
	return Model{
		common:    c,
		JqlSearch: jqlsearch,
		IssueList: issuelist,
	}
}

// SetSize updates the available dimensions.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height

	jqlBarHeight := m.JqlSearch.SetSize(m.width)
	m.IssueList.SetSize(m.width, m.height-jqlBarHeight)
}

func (m *Model) View() string {
	if m.width <= 0 || m.height <= 0 {
		return ""
	}

	jqlsearchView := m.JqlSearch.View()
	issuelistView := m.IssueList.View()

	return lipgloss.JoinVertical(lipgloss.Top, jqlsearchView, issuelistView)
}
