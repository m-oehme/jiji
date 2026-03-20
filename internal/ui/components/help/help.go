// Package help renders the help overlay listing all keybindings.
package help

import (
	"fmt"
	"strings"

	"github.com/m-oehme/jiji/internal/config"
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
	styles   *styles.Styles
	keys     *config.KeyConfig
	width    int
	height   int
	scroll   int
	lines    []string // precomputed content lines
	built    bool
}

// New creates a help overlay.
func New(keys *config.KeyConfig, s *styles.Styles) Model {
	return Model{
		keys:   keys,
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
	bindings := []binding{
		{"Move up", m.keys.Up},
		{"Move down", m.keys.Down},
		{"Jump to top", m.keys.Top},
		{"Jump to bottom", m.keys.Bottom},
		{"Next tab", m.keys.TabNext},
		{"Previous tab", m.keys.TabPrev},
		{"Switch pane", m.keys.PaneSwitch},
		{"Focus JQL", m.keys.FocusJQL},
		{"Confirm", m.keys.Confirm},
		{"Cancel / Back", m.keys.Cancel},
		{"Transition", m.keys.Transition},
		{"Comment", m.keys.Comment},
		{"Labels", m.keys.Labels},
		{"Edit summary", m.keys.Summary},
		{"Edit description", m.keys.Edit},
		{"Refresh", m.keys.Refresh},
		{"Help", m.keys.Help},
		{"Quit", m.keys.Quit},
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
