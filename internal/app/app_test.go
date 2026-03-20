package app

import (
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/m-oehme/jiji/internal/config"
	"github.com/m-oehme/jiji/internal/ui/common"
)

func testConfig() *config.Config {
	return &config.Config{
		UI: config.UIConfig{
			ListRatio:    30,
			DetailLayout: "stacked",
		},
		Tabs: []config.TabConfig{
			{Name: "All Issues", JQL: "ORDER BY updated DESC"},
		},
		Keys: config.KeyConfig{
			Up:         []string{"k", "up"},
			Down:       []string{"j", "down"},
			Top:        []string{"g"},
			Bottom:     []string{"G"},
			TabNext:    []string{"l", "right"},
			TabPrev:    []string{"h", "left"},
			PaneSwitch: []string{"tab"},
			Help:       []string{"?"},
			Quit:       []string{"q"},
			Cancel:     []string{"esc"},
			Confirm:    []string{"enter"},
			FocusJQL:   []string{"/"},
			Transition: []string{"t"},
			Comment:    []string{"c"},
			Labels:     []string{"L"},
			Summary:    []string{"s"},
			Edit:       []string{"e"},
			Refresh:    []string{"r"},
		},
		Theme: config.ThemeConfig{
			Primary:       "#7b2fbe",
			Secondary:     "#00d4aa",
			Border:        "#555555",
			BorderFocused: "#7b2fbe",
			Text:          "#e0e0e0",
			TextDim:       "#888888",
			Success:       "#00d4aa",
			Warning:       "#f0c040",
			Error:         "#ff4444",
		},
	}
}

func setupModel(t *testing.T) Model {
	t.Helper()
	cfg := testConfig()
	m := New(cfg, nil)
	// Simulate window resize
	sized, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	return sized.(Model)
}

func TestApp_InitialState(t *testing.T) {
	m := setupModel(t)

	if m.focus.ActivePane() != common.PaneIssueList {
		t.Fatal("expected IssueList pane focused initially")
	}
	if m.focus.HasOverlay() {
		t.Fatal("expected no overlay initially")
	}
	if m.issueList.SelectedIssue() == nil {
		t.Fatal("expected mock issues loaded")
	}
}

func TestApp_PaneSwitch(t *testing.T) {
	m := setupModel(t)

	// Tab key switches pane
	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = updated.(Model)

	if m.focus.ActivePane() != common.PaneDetail {
		t.Fatal("expected Detail pane after Tab")
	}

	updated, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = updated.(Model)

	if m.focus.ActivePane() != common.PaneIssueList {
		t.Fatal("expected IssueList pane after second Tab")
	}
}

func TestApp_IssueListNavigation(t *testing.T) {
	m := setupModel(t)
	startIdx := m.issueList.SelectedIndex()

	// j moves down
	updated, _ := m.Update(tea.KeyPressMsg{Text: "j", Code: 'j'})
	m = updated.(Model)
	if m.issueList.SelectedIndex() != startIdx+1 {
		t.Fatalf("expected index %d after j, got %d", startIdx+1, m.issueList.SelectedIndex())
	}

	// k moves up
	updated, _ = m.Update(tea.KeyPressMsg{Text: "k", Code: 'k'})
	m = updated.(Model)
	if m.issueList.SelectedIndex() != startIdx {
		t.Fatalf("expected index %d after k, got %d", startIdx, m.issueList.SelectedIndex())
	}

	// G jumps to bottom
	updated, _ = m.Update(tea.KeyPressMsg{Text: "G", Code: 'G'})
	m = updated.(Model)
	if m.issueList.SelectedIndex() == 0 {
		t.Fatal("expected non-zero index after G")
	}

	// g jumps to top
	updated, _ = m.Update(tea.KeyPressMsg{Text: "g", Code: 'g'})
	m = updated.(Model)
	if m.issueList.SelectedIndex() != 0 {
		t.Fatalf("expected index 0 after g, got %d", m.issueList.SelectedIndex())
	}
}

func TestApp_HelpOverlay(t *testing.T) {
	m := setupModel(t)

	// ? opens help
	updated, _ := m.Update(tea.KeyPressMsg{Text: "?", Code: '?'})
	m = updated.(Model)

	if !m.focus.HasOverlay() {
		t.Fatal("expected help overlay after ?")
	}
	if m.focus.TopOverlay() != common.OverlayHelp {
		t.Fatalf("expected OverlayHelp, got %d", m.focus.TopOverlay())
	}

	// Esc closes help
	updated, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	m = updated.(Model)

	if m.focus.HasOverlay() {
		t.Fatal("expected no overlay after Esc")
	}
}

func TestApp_HelpToggle(t *testing.T) {
	m := setupModel(t)

	// ? opens help
	updated, _ := m.Update(tea.KeyPressMsg{Text: "?", Code: '?'})
	m = updated.(Model)
	if !m.focus.HasOverlay() {
		t.Fatal("expected overlay open")
	}

	// ? again closes help
	updated, _ = m.Update(tea.KeyPressMsg{Text: "?", Code: '?'})
	m = updated.(Model)
	if m.focus.HasOverlay() {
		t.Fatal("expected overlay closed after second ?")
	}
}

func TestApp_QuitKey(t *testing.T) {
	m := setupModel(t)

	_, cmd := m.Update(tea.KeyPressMsg{Text: "q", Code: 'q'})
	if cmd == nil {
		t.Fatal("expected quit command from q key")
	}
	// Execute the cmd and check it returns QuitMsg
	msg := cmd()
	if _, ok := msg.(tea.QuitMsg); !ok {
		t.Fatalf("expected QuitMsg, got %T", msg)
	}
}

func TestApp_CtrlC_AlwaysQuits(t *testing.T) {
	m := setupModel(t)

	// Even with overlay open
	m.focus.PushOverlay(common.OverlayHelp)

	_, cmd := m.Update(tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl})
	if cmd == nil {
		t.Fatal("expected quit command from Ctrl+C")
	}
	msg := cmd()
	if _, ok := msg.(tea.QuitMsg); !ok {
		t.Fatalf("expected QuitMsg from Ctrl+C, got %T", msg)
	}
}

func TestApp_OverlayBlocksNavigation(t *testing.T) {
	m := setupModel(t)
	startIdx := m.issueList.SelectedIndex()

	// Open help overlay
	updated, _ := m.Update(tea.KeyPressMsg{Text: "?", Code: '?'})
	m = updated.(Model)

	// j should NOT navigate issue list when overlay is active
	updated, _ = m.Update(tea.KeyPressMsg{Text: "j", Code: 'j'})
	m = updated.(Model)

	if m.issueList.SelectedIndex() != startIdx {
		t.Fatal("navigation should be blocked when overlay is active")
	}
}

func TestApp_DetailScrolling(t *testing.T) {
	m := setupModel(t)

	// Switch to detail pane
	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = updated.(Model)

	// j/k should scroll detail, not navigate issue list
	startIdx := m.issueList.SelectedIndex()
	updated, _ = m.Update(tea.KeyPressMsg{Text: "j", Code: 'j'})
	m = updated.(Model)

	if m.issueList.SelectedIndex() != startIdx {
		t.Fatal("j in detail pane should not move issue list cursor")
	}
}

func TestApp_WindowResize(t *testing.T) {
	cfg := testConfig()
	m := New(cfg, nil)

	updated, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = updated.(Model)

	if m.width != 80 || m.height != 24 {
		t.Fatalf("expected 80x24, got %dx%d", m.width, m.height)
	}

	// View should not be empty after resize
	v := m.View()
	if v.Content == "" {
		t.Fatal("expected non-empty view after resize")
	}
}

func TestMatchKey(t *testing.T) {
	tests := []struct {
		key      string
		bindings []string
		want     bool
	}{
		{"j", []string{"j", "down"}, true},
		{"down", []string{"j", "down"}, true},
		{"k", []string{"j", "down"}, false},
		{"j", nil, false},
		{"j", []string{}, false},
	}
	for _, tt := range tests {
		got := matchKey(tt.key, tt.bindings)
		if got != tt.want {
			t.Errorf("matchKey(%q, %v) = %v, want %v", tt.key, tt.bindings, got, tt.want)
		}
	}
}
