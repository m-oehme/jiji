package borderbox

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/m-oehme/jiji/internal/ui/common"
)

type Model struct {
	ctx    *common.Context
	common *common.Common

	style lipgloss.Style

	width  int
	height int
}

func New(ctx *common.Context, c *common.Common, focused bool) Model {
	borderStyle := c.Styles.Border
	if focused {
		borderStyle = c.Styles.BorderFocused
	}
	return Model{
		ctx:    ctx,
		common: c,
		style:  borderStyle,
	}
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *Model) GetContentSize() (int, int) {
	frameW, frameH := m.style.GetFrameSize()
	contentW := m.width - frameW
	contentH := m.height - frameH
	return contentW, contentH
}

func (m *Model) Render(content, title string) string {
	contentW, contentH := m.GetContentSize()
	if contentW <= 0 || contentH <= 0 {
		return m.style.Width(m.width).Height(m.height).Render("")
	}

	border := m.style.
		Width(m.width).
		Height(m.height).
		MaxWidth(m.width).
		MaxHeight(m.height).
		Render(content)

	if title == "" {
		return border
	}

	borderColor := m.style.GetBorderTopForeground()
	styledTitle := lipgloss.NewStyle().Foreground(borderColor).Render(" " + title + " ")

	lines := strings.SplitN(border, "\n", 2)
	lines[0] = common.ReplaceAt(lines[0], 2, styledTitle)
	return strings.Join(lines, "\n")
}
