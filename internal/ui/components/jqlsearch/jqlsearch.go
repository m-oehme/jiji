package jqlsearch

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/m-oehme/jiji/internal/ui/common"
	"github.com/m-oehme/jiji/internal/ui/components/borderbox"
)

type Model struct {
	common   *common.Common
	jql      string // last submitted JQL
	jqlInput textinput.Model
	jqlFocus bool // true when JQL input is being edited

	width  int
	height int
}

const (
	componentHeight = 3
	lineHeight      = 1
)

func New(c *common.Common) Model {
	ti := textinput.New()
	ti.Prompt = "JQL: "
	ti.SetVirtualCursor(true)
	return Model{
		common:   c,
		jqlInput: ti,
	}
}

// SetSize updates the available dimensions.
func (m *Model) SetSize(width int) int {
	m.width = width
	m.height = componentHeight

	frameW, _ := m.common.Styles.Border.GetFrameSize()
	contentW := width - frameW

	// Account for the "JQL: " prompt width
	inputW := contentW - lipgloss.Width(m.jqlInput.Prompt)
	if inputW < 1 {
		inputW = 1
	}
	m.jqlInput.SetWidth(inputW)

	return m.height
}

func (m *Model) View() string {
	if m.width <= 0 || m.height <= 0 {
		return ""
	}

	border := borderbox.New(m.common, m.common.Focused)
	border.SetSize(m.width, m.height)
	contentW, _ := border.GetContentSize()

	// JQL line: interactive textinput when focused, static text when unfocused.
	var jqlLine string
	if m.jqlFocus {
		jqlLine = m.jqlInput.View()
	} else {
		jqlLine = m.common.Styles.Dimmed.
			Height(1).
			Width(contentW).
			MaxHeight(1).
			MaxWidth(contentW).
			Render(common.Truncate(m.jql, contentW))
	}

	return border.Render(jqlLine, "JQL")
}

// SetJQL sets the JQL query string displayed above the table.
func (m *Model) SetJQL(jql string) {
	m.jql = jql
	m.jqlInput.SetValue(jql)
}

// JQLValue returns the current JQL string (from the text input when focused, or last submitted).
func (m *Model) JQLValue() string {
	if m.jqlFocus {
		return m.jqlInput.Value()
	}
	return m.jql
}

// FocusJQL activates the JQL text input for editing.
func (m *Model) FocusJQL() tea.Cmd {
	m.jqlFocus = true
	m.jqlInput.SetValue(m.jql)
	m.jqlInput.CursorEnd()
	return m.jqlInput.Focus()
}

// UnfocusJQL deactivates the JQL text input, reverting to the last submitted value.
func (m *Model) UnfocusJQL() {
	m.jqlFocus = false
	m.jqlInput.Blur()
	// On cancel, revert to last submitted JQL
	m.jqlInput.SetValue(m.jql)
}

// IsJQLFocused returns whether the JQL input is being edited.
func (m *Model) IsJQLFocused() bool {
	return m.jqlFocus
}

// UpdateJQL delegates a key message to the textinput and returns any command.
func (m *Model) UpdateJQL(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.jqlInput, cmd = m.jqlInput.Update(msg)
	return cmd
}
