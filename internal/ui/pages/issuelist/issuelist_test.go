package issuelist

import (
	"log/slog"
	"testing"

	"github.com/m-oehme/jiji/internal/config"
	"github.com/m-oehme/jiji/internal/jira"
	"github.com/m-oehme/jiji/internal/ui/common"
	"github.com/m-oehme/jiji/internal/ui/styles"
)

func testContext() *common.Context {
	return &common.Context{
		Config: &config.Config{
			UI: config.UIConfig{
				Fields: config.FieldsConfig{
					List: []string{"key", "priority", "assignee", "summary"},
				},
			},
		},
		Logger: slog.New(slog.DiscardHandler),
	}
}

func testCommon() *common.Common {
	s := styles.NewStyles(config.ThemeConfig{
		Primary:       "#7b2fbe",
		Secondary:     "#00d4aa",
		Border:        "#555555",
		BorderFocused: "#7b2fbe",
		Text:          "#e0e0e0",
		TextDim:       "#888888",
		Success:       "#00d4aa",
		Warning:       "#f0c040",
		Error:         "#ff4444",
	})
	return &common.Common{
		Styles:  s,
		Focused: true,
	}
}

func testIssues() []jira.Issue {
	return []jira.Issue{
		{Key: "TEST-1", Summary: "First issue", Priority: "High", Assignee: "Alice"},
		{Key: "TEST-2", Summary: "Second issue", Priority: "Medium", Assignee: "Bob"},
		{Key: "TEST-3", Summary: "Third issue", Priority: "Low", Assignee: "Charlie"},
	}
}

func TestModel_Navigation(t *testing.T) {
	c := testCommon()
	m := New(testContext(), c)
	m.SetItems(testIssues())

	if m.SelectedIndex() != 0 {
		t.Fatalf("expected index 0, got %d", m.SelectedIndex())
	}

	m.MoveDown()
	if m.SelectedIndex() != 1 {
		t.Fatalf("expected index 1 after MoveDown, got %d", m.SelectedIndex())
	}

	m.MoveDown()
	if m.SelectedIndex() != 2 {
		t.Fatalf("expected index 2, got %d", m.SelectedIndex())
	}

	// Can't move past end
	m.MoveDown()
	if m.SelectedIndex() != 2 {
		t.Fatalf("expected index 2 (clamped), got %d", m.SelectedIndex())
	}

	m.MoveUp()
	if m.SelectedIndex() != 1 {
		t.Fatalf("expected index 1 after MoveUp, got %d", m.SelectedIndex())
	}

	// Can't move before start
	m.JumpToTop()
	m.MoveUp()
	if m.SelectedIndex() != 0 {
		t.Fatalf("expected index 0 (clamped), got %d", m.SelectedIndex())
	}
}

func TestModel_JumpToTopBottom(t *testing.T) {
	c := testCommon()
	m := New(testContext(), c)
	m.SetItems(testIssues())

	m.JumpToBottom()
	if m.SelectedIndex() != 2 {
		t.Fatalf("expected index 2 after JumpToBottom, got %d", m.SelectedIndex())
	}

	m.JumpToTop()
	if m.SelectedIndex() != 0 {
		t.Fatalf("expected index 0 after JumpToTop, got %d", m.SelectedIndex())
	}
}

func TestModel_SelectedIssue(t *testing.T) {
	c := testCommon()
	m := New(testContext(), c)

	// Empty list returns nil
	if m.SelectedIssue() != nil {
		t.Fatal("expected nil for empty list")
	}

	m.SetItems(testIssues())
	sel := m.SelectedIssue()
	if sel == nil {
		t.Fatal("expected non-nil selected issue")
	}
	if sel.Key != "TEST-1" {
		t.Fatalf("expected TEST-1, got %s", sel.Key)
	}

	m.MoveDown()
	sel = m.SelectedIssue()
	if sel.Key != "TEST-2" {
		t.Fatalf("expected TEST-2, got %s", sel.Key)
	}
}

func TestModel_SetItems_CursorClamp(t *testing.T) {
	c := testCommon()
	m := New(testContext(), c)
	m.SetItems(testIssues())
	m.JumpToBottom() // cursor = 2

	// Replace with shorter list
	m.SetItems([]jira.Issue{
		{Key: "TEST-1", Summary: "Only one"},
	})
	if m.SelectedIndex() != 0 {
		t.Fatalf("expected cursor clamped to 0, got %d", m.SelectedIndex())
	}
}

func TestModel_View_NonEmpty(t *testing.T) {
	c := testCommon()
	m := New(testContext(), c)
	m.SetItems(testIssues())
	m.SetSize(80, 20)

	view := m.View()
	if view == "" {
		t.Fatal("expected non-empty view")
	}
}

func TestModel_View_ZeroSize(t *testing.T) {
	c := testCommon()
	m := New(testContext(), c)
	m.SetSize(0, 0)

	if m.View() != "" {
		t.Fatal("expected empty view for zero size")
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		input  string
		maxLen int
		want   string
	}{
		{"hello", 10, "hello"},
		{"hello world", 5, "hell…"},
		{"hi", 2, "hi"},
		{"hi", 1, "…"},
		{"hi", 0, ""},
		{"", 5, ""},
	}
	for _, tt := range tests {
		got := common.Truncate(tt.input, tt.maxLen)
		if got != tt.want {
			t.Errorf("common.Truncate(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.want)
		}
	}
}
