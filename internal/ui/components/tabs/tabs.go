// Package tabs renders the tab bar at the top of the screen (ADR-007).
package tabs

import (
	"fmt"
	"strings"

	"github.com/m-oehme/jiji/internal/config"
	"github.com/m-oehme/jiji/internal/ui/common"
	"github.com/m-oehme/jiji/internal/ui/styles"
	lipgloss "charm.land/lipgloss/v2"
)

// Model represents the tab bar component.
type Model struct {
	ctx    *common.Context
	tabs   []config.TabConfig
	active int
	width  int
	height int
	styles *styles.Styles
}

// New creates a tab bar from the configured tabs.
func New(ctx *common.Context, tabs []config.TabConfig, s *styles.Styles) Model {
	return Model{
		ctx:    ctx,
		tabs:   tabs,
		active: 0,
		styles: s,
		height: 1,
	}
}

// SetActive sets the active tab index.
func (m *Model) SetActive(idx int) {
	if idx >= 0 && idx < len(m.tabs) {
		m.active = idx
	}
}

// Active returns the index of the active tab.
func (m *Model) Active() int {
	return m.active
}

// Count returns the number of tabs.
func (m *Model) Count() int {
	return len(m.tabs)
}

// ActiveTab returns the active TabConfig.
func (m *Model) ActiveTab() config.TabConfig {
	if m.active < len(m.tabs) {
		return m.tabs[m.active]
	}
	return config.TabConfig{}
}

// SetSize updates the available width for the tab bar.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// View renders the tab bar.
func (m Model) View() string {
	if len(m.tabs) == 0 || m.width <= 0 {
		return ""
	}

	var parts []string
	remaining := m.width

	for i, tab := range m.tabs {
		label := fmt.Sprintf("%d %s", i+1, tab.Name)

		// Truncate if we're running out of space
		maxLen := remaining - 2 // leave room for padding
		if maxLen < 4 {
			break
		}
		if len(label) > maxLen {
			label = label[:maxLen-1] + "…"
		}

		var rendered string
		if i == m.active {
			rendered = m.styles.TabActive.Render(label)
		} else {
			rendered = m.styles.TabInactive.Render(label)
		}
		parts = append(parts, rendered)
		remaining -= lipgloss.Width(rendered)
	}

	row := strings.Join(parts, "")

	// Pad to full width
	if lipgloss.Width(row) < m.width {
		row += strings.Repeat(" ", m.width-lipgloss.Width(row))
	}

	return row
}
