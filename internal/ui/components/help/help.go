// Package help renders the help overlay listing all keybindings.
package help

import (
	"fmt"
	"strings"

	"github.com/m-oehme/jiji/internal/ui/common"
	"github.com/m-oehme/jiji/internal/ui/styles"
	lipgloss "charm.land/lipgloss/v2"
)

// binding pairs a description with its key strings.
type binding struct {
	desc string
	keys []string
}

// Model represents the help overlay.
type Model struct {
	ctx    *common.Context
	styles *styles.Styles
	width  int
	height int
	scroll int
	lines  []string // precomputed content lines
	built  bool
}

// New creates a help overlay.
func New(ctx *common.Context, s *styles.Styles) Model {
	return Model{
		ctx:    ctx,
		styles: s,
	}
}

// SetSize updates the help overlay dimensions.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.built = false // rebuild on resize
}

// ScrollDown scrolls the help text down.
func (m *Model) ScrollDown() {
	maxScroll := len(m.lines) - m.innerHeight()
	if maxScroll < 0 {
		maxScroll = 0
	}
	if m.scroll < maxScroll {
		m.scroll++
	}
}

// ScrollUp scrolls the help text up.
func (m *Model) ScrollUp() {
	if m.scroll > 0 {
		m.scroll--
	}
}

func (m *Model) innerHeight() int {
	// Account for border + title + padding
	return m.height - 6
}

func (m *Model) build() {
	keys := &m.ctx.Config.Keys
	bindings := []binding{
		{"Move up", keys.Up},
		{"Move down", keys.Down},
		{"Jump to top", keys.Top},
		{"Jump to bottom", keys.Bottom},
		{"Next tab", keys.TabNext},
		{"Previous tab", keys.TabPrev},
		{"Switch pane", keys.PaneSwitch},
		{"Focus JQL", keys.FocusJQL},
		{"Confirm", keys.Confirm},
		{"Cancel / Back", keys.Cancel},
		{"Transition", keys.Transition},
		{"Comment", keys.Comment},
		{"Labels", keys.Labels},
		{"Edit summary", keys.Summary},
		{"Edit description", keys.Edit},
		{"Refresh", keys.Refresh},
		{"Help", keys.Help},
		{"Quit", keys.Quit},
	}

	keyStyle := m.styles.Heading
	descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#cccccc"))

	m.lines = nil
	m.lines = append(m.lines, "")
	for _, b := range bindings {
		if len(b.keys) == 0 {
			continue
		}
		keyStr := strings.Join(b.keys, ", ")
		line := fmt.Sprintf("  %s  %s",
			keyStyle.Width(16).Render(keyStr),
			descStyle.Render(b.desc),
		)
		m.lines = append(m.lines, line)
	}
	m.lines = append(m.lines, "")
	m.lines = append(m.lines, m.styles.Dimmed.Render("  Press ? or Esc to close"))
	m.built = true
}

// View renders the help overlay as a centered box.
func (m Model) View() string {
	if !m.built {
		m.build()
	}

	ih := m.innerHeight()
	if ih <= 0 {
		return ""
	}

	// Visible slice
	start := m.scroll
	end := start + ih
	if end > len(m.lines) {
		end = len(m.lines)
	}
	if start > end {
		start = end
	}

	visible := strings.Join(m.lines[start:end], "\n")

	boxW := m.width * 2 / 3
	if boxW < 40 {
		boxW = min(40, m.width-4)
	}

	title := m.styles.Heading.Render(" Keybindings ")

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7b2fbe")).
		Width(boxW).
		Height(ih + 2).
		Padding(1, 2).
		Render(title + "\n" + visible)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
}
