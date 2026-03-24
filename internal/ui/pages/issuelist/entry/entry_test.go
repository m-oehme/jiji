package entry

import (
	"log/slog"
	"testing"

	"github.com/m-oehme/jiji/internal/config"
	"github.com/m-oehme/jiji/internal/jira"
	"github.com/m-oehme/jiji/internal/ui/common"
	"github.com/m-oehme/jiji/internal/ui/styles"
)

func testContext(cols []string) *common.Context {
	return &common.Context{
		Config: &config.Config{
			UI: config.UIConfig{
				Fields: config.FieldsConfig{
					List: cols,
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

func TestColumnsFromConfig_Order(t *testing.T) {
	cols := ColumnsFromConfig([]string{"summary", "key", "priority"})
	if len(cols) != 3 {
		t.Fatalf("expected 3 columns, got %d", len(cols))
	}
	if cols[0].ID != "summary" {
		t.Errorf("expected first column 'summary', got %q", cols[0].ID)
	}
	if cols[1].ID != "key" {
		t.Errorf("expected second column 'key', got %q", cols[1].ID)
	}
	if cols[2].ID != "priority" {
		t.Errorf("expected third column 'priority', got %q", cols[2].ID)
	}
}

func TestColumnsFromConfig_SkipsUnknown(t *testing.T) {
	cols := ColumnsFromConfig([]string{"key", "bogus", "summary"})
	if len(cols) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(cols))
	}
	if cols[0].ID != "key" {
		t.Errorf("expected 'key', got %q", cols[0].ID)
	}
	if cols[1].ID != "summary" {
		t.Errorf("expected 'summary', got %q", cols[1].ID)
	}
}

func TestColumnsFromConfig_EnsuresFlex(t *testing.T) {
	cols := ColumnsFromConfig([]string{"key", "priority", "assignee"})
	if len(cols) != 3 {
		t.Fatalf("expected 3 columns, got %d", len(cols))
	}
	// Last column should become flex (width 0)
	if cols[2].Width != 0 {
		t.Errorf("expected last column to be flex (width 0), got %d", cols[2].Width)
	}
}

func TestColumnsFromConfig_Empty(t *testing.T) {
	cols := ColumnsFromConfig(nil)
	if cols != nil {
		t.Fatalf("expected nil, got %v", cols)
	}
	cols = ColumnsFromConfig([]string{"bogus"})
	if cols != nil {
		t.Fatalf("expected nil for all unknown IDs, got %v", cols)
	}
}

func TestModel_View_NonEmpty(t *testing.T) {
	c := testCommon()
	cols := []string{"key", "priority", "assignee", "summary"}
	m := New(testContext(cols), c)
	m.SetSize(80)
	m.SetIssue(jira.Issue{Key: "TEST-1", Summary: "A summary", Priority: "High", Assignee: "Alice"})
	m.SetSelected(false)

	view := m.View()
	if view == "" {
		t.Fatal("expected non-empty view")
	}
}

func TestModel_View_ZeroWidth(t *testing.T) {
	c := testCommon()
	cols := []string{"key", "summary"}
	m := New(testContext(cols), c)
	m.SetSize(0)

	if m.View() != "" {
		t.Fatal("expected empty view for zero width")
	}
}
