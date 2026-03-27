package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/m-oehme/jiji/internal/config"
	"github.com/m-oehme/jiji/internal/jira"
	"github.com/m-oehme/jiji/internal/ui/common"
)

// mockClient implements jira.Client for testing.
type mockClient struct {
	issues   []jira.Issue
	comments []jira.Comment
	searchFn func(jql string) (*jira.SearchResult, error)
}

func (c *mockClient) SearchIssues(_ context.Context, jql string, _ []string, _ jira.PageToken) (*jira.SearchResult, error) {
	if c.searchFn != nil {
		return c.searchFn(jql)
	}
	return &jira.SearchResult{Issues: c.issues, Total: len(c.issues)}, nil
}

func (c *mockClient) GetIssue(_ context.Context, key string) (*jira.Issue, error) {
	for i := range c.issues {
		if c.issues[i].Key == key {
			return &c.issues[i], nil
		}
	}
	return nil, fmt.Errorf("issue %s not found", key)
}

func (c *mockClient) GetComments(_ context.Context, _ string) ([]jira.Comment, error) {
	return c.comments, nil
}

func (c *mockClient) GetFieldMetadata(_ context.Context) ([]jira.FieldMetadata, error) {
	return nil, nil
}

func (c *mockClient) GetAutocompleteSuggestions(_ context.Context, _, _ string) ([]jira.Suggestion, error) {
	return nil, nil
}

func (c *mockClient) GetTransitions(_ context.Context, _ string) ([]jira.Transition, error) {
	return nil, nil
}

func (c *mockClient) TransitionIssue(_ context.Context, _, _ string) error { return nil }

func (c *mockClient) UpdateIssue(_ context.Context, _ string, _ map[string]any) error { return nil }

func (c *mockClient) AddComment(_ context.Context, _ string, _ json.RawMessage) error { return nil }

func (c *mockClient) UpdateComment(_ context.Context, _, _ string, _ json.RawMessage) error {
	return nil
}

func (c *mockClient) GetLabels(_ context.Context) ([]string, error) { return nil, nil }

func testConfig() *config.Config {
	return &config.Config{
		UI: config.UIConfig{
			ListRatio:    30,
			DetailLayout: "stacked",
		},
		Tabs: []config.TabConfig{
			{Name: "My Issues", JQL: "assignee = currentUser() AND resolution = Unresolved ORDER BY updated DESC"},
		},
		Keys: config.Keybindings{
			Builtin: config.BuiltinKeybindings{
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

func testIssues() []jira.Issue {
	now := time.Now()
	return []jira.Issue{
		{Key: "TEST-1", Summary: "First issue", Status: "Open", Priority: "High", Assignee: "Alice", Created: now, Updated: now},
		{Key: "TEST-2", Summary: "Second issue", Status: "In Progress", Priority: "Medium", Assignee: "Bob", Created: now, Updated: now},
		{Key: "TEST-3", Summary: "Third issue", Status: "Done", Priority: "Low", Assignee: "Charlie", Created: now, Updated: now},
	}
}

func testClient() *mockClient {
	return &mockClient{
		issues: testIssues(),
		comments: []jira.Comment{
			{ID: "1", Author: "Alice", Body: json.RawMessage(`{"version":1,"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"A comment"}]}]}`), Created: time.Now(), Updated: time.Now()},
		},
	}
}

// setupModel creates a model with mock client and simulates initial window resize.
func setupModel(t *testing.T) Model {
	t.Helper()
	cfg := testConfig()
	client := testClient()
	m := New(cfg, client, slog.New(slog.DiscardHandler))
	// Simulate window resize
	sized, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = sized.(Model)
	// Simulate search results arriving
	m.issuepane.IssueList.SetItems(client.issues)
	m.statusBar.SetIssueCount(len(client.issues))
	if sel := m.issuepane.IssueList.SelectedIssue(); sel != nil {
		m.statusBar.SetCurrentIssue(sel.Key)
	}
	return m
}

// setupModelNoClient creates a model without a client (for tests that don't need API calls).
func setupModelNoClient(t *testing.T) Model {
	t.Helper()
	cfg := testConfig()
	m := New(cfg, nil, slog.New(slog.DiscardHandler))
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
	if m.issuepane.IssueList.SelectedIssue() == nil {
		t.Fatal("expected issues loaded")
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
	startIdx := m.issuepane.IssueList.SelectedIndex()

	// j moves down
	updated, _ := m.Update(tea.KeyPressMsg{Text: "j", Code: 'j'})
	m = updated.(Model)
	if m.issuepane.IssueList.SelectedIndex() != startIdx+1 {
		t.Fatalf("expected index %d after j, got %d", startIdx+1, m.issuepane.IssueList.SelectedIndex())
	}

	// k moves up
	updated, _ = m.Update(tea.KeyPressMsg{Text: "k", Code: 'k'})
	m = updated.(Model)
	if m.issuepane.IssueList.SelectedIndex() != startIdx {
		t.Fatalf("expected index %d after k, got %d", startIdx, m.issuepane.IssueList.SelectedIndex())
	}

	// G jumps to bottom
	updated, _ = m.Update(tea.KeyPressMsg{Text: "G", Code: 'G'})
	m = updated.(Model)
	if m.issuepane.IssueList.SelectedIndex() == 0 {
		t.Fatal("expected non-zero index after G")
	}

	// g jumps to top
	updated, _ = m.Update(tea.KeyPressMsg{Text: "g", Code: 'g'})
	m = updated.(Model)
	if m.issuepane.IssueList.SelectedIndex() != 0 {
		t.Fatalf("expected index 0 after g, got %d", m.issuepane.IssueList.SelectedIndex())
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
	startIdx := m.issuepane.IssueList.SelectedIndex()

	// Open help overlay
	updated, _ := m.Update(tea.KeyPressMsg{Text: "?", Code: '?'})
	m = updated.(Model)

	// j should NOT navigate issue list when overlay is active
	updated, _ = m.Update(tea.KeyPressMsg{Text: "j", Code: 'j'})
	m = updated.(Model)

	if m.issuepane.IssueList.SelectedIndex() != startIdx {
		t.Fatal("navigation should be blocked when overlay is active")
	}
}

func TestApp_DetailScrolling(t *testing.T) {
	m := setupModel(t)

	// Switch to detail pane
	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = updated.(Model)

	// j/k should scroll detail, not navigate issue list
	startIdx := m.issuepane.IssueList.SelectedIndex()
	updated, _ = m.Update(tea.KeyPressMsg{Text: "j", Code: 'j'})
	m = updated.(Model)

	if m.issuepane.IssueList.SelectedIndex() != startIdx {
		t.Fatal("j in detail pane should not move issue list cursor")
	}
}

func TestApp_WindowResize(t *testing.T) {
	cfg := testConfig()
	m := New(cfg, nil, slog.New(slog.DiscardHandler))

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

// --- Phase 3 tests ---

func TestApp_SearchResultMsg(t *testing.T) {
	m := setupModelNoClient(t)

	issues := testIssues()
	updated, cmd := m.Update(SearchResultMsg{Issues: issues, TabIndex: 0})
	m = updated.(Model)

	if m.issuepane.IssueList.SelectedIssue() == nil {
		t.Fatal("expected issues after SearchResultMsg")
	}
	if m.issuepane.IssueList.SelectedIssue().Key != "TEST-1" {
		t.Fatalf("expected first issue TEST-1, got %s", m.issuepane.IssueList.SelectedIssue().Key)
	}
	// Should return a batch cmd to load detail + comments
	if cmd == nil {
		t.Fatal("expected cmd to load detail after search results")
	}
}

func TestApp_IssueDetailMsg(t *testing.T) {
	m := setupModel(t)

	issue := &jira.Issue{Key: "DETAIL-1", Summary: "Detail test", Status: "Open"}
	updated, _ := m.Update(IssueDetailMsg{Issue: issue})
	m = updated.(Model)

	if m.detail.CurrentIssue() == nil {
		t.Fatal("expected issue set on detail")
	}
	if m.detail.CurrentIssue().Key != "DETAIL-1" {
		t.Fatalf("expected DETAIL-1, got %s", m.detail.CurrentIssue().Key)
	}
}

func TestApp_CommentsMsg(t *testing.T) {
	m := setupModel(t)

	comments := []jira.Comment{
		{ID: "1", Author: "Tester", Created: time.Now(), Updated: time.Now()},
	}
	updated, _ := m.Update(CommentsMsg{Comments: comments, IssueKey: "TEST-1"})
	_ = updated.(Model)
	// No panic = success; comments are set internally
}

func TestApp_ErrorMsg(t *testing.T) {
	m := setupModel(t)

	updated, cmd := m.Update(ErrorMsg{Err: fmt.Errorf("connection refused"), Context: "search"})
	_ = updated.(Model)

	// Should return a cmd for auto-clear
	if cmd == nil {
		t.Fatal("expected auto-clear cmd after error")
	}
}

func TestApp_ClearErrorMsg(t *testing.T) {
	m := setupModel(t)

	// Set an error
	m.statusBar.SetError(fmt.Errorf("test error"))

	// Clear it
	updated, _ := m.Update(clearErrorMsg{})
	_ = updated.(Model)
	// No panic = success
}

func TestApp_NavigationTriggersDetailLoad(t *testing.T) {
	m := setupModel(t)

	// j moves down and should return a cmd (when client != nil)
	updated, cmd := m.Update(tea.KeyPressMsg{Text: "j", Code: 'j'})
	_ = updated.(Model)

	if cmd == nil {
		t.Fatal("expected cmd to load detail on navigation")
	}
}

func TestApp_JQLFocus(t *testing.T) {
	m := setupModel(t)

	// / focuses JQL
	updated, _ := m.Update(tea.KeyPressMsg{Text: "/", Code: '/'})
	m = updated.(Model)

	if !m.issuepane.JqlSearch.IsJQLFocused() {
		t.Fatal("expected JQL focused after /")
	}

	// Esc cancels JQL
	updated, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	m = updated.(Model)

	if m.issuepane.JqlSearch.IsJQLFocused() {
		t.Fatal("expected JQL unfocused after Esc")
	}
}

func TestApp_JQLSubmit(t *testing.T) {
	m := setupModel(t)

	// Focus JQL
	updated, _ := m.Update(tea.KeyPressMsg{Text: "/", Code: '/'})
	m = updated.(Model)

	// Submit with Enter
	updated, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = updated.(Model)

	if m.issuepane.JqlSearch.IsJQLFocused() {
		t.Fatal("expected JQL unfocused after Enter")
	}
	if cmd == nil {
		t.Fatal("expected search cmd after JQL submit")
	}
}

func TestApp_JQLBlocksOtherKeys(t *testing.T) {
	m := setupModel(t)
	startIdx := m.issuepane.IssueList.SelectedIndex()

	// Focus JQL
	updated, _ := m.Update(tea.KeyPressMsg{Text: "/", Code: '/'})
	m = updated.(Model)

	// j should go to textinput, not navigate
	updated, _ = m.Update(tea.KeyPressMsg{Text: "j", Code: 'j'})
	m = updated.(Model)

	if m.issuepane.IssueList.SelectedIndex() != startIdx {
		t.Fatal("j should not navigate issue list when JQL is focused")
	}

	// q should go to textinput, not quit
	_, cmd := m.Update(tea.KeyPressMsg{Text: "q", Code: 'q'})
	if cmd != nil {
		// Check it's not a quit command
		msg := cmd()
		if _, ok := msg.(tea.QuitMsg); ok {
			t.Fatal("q should not quit when JQL is focused")
		}
	}
}
