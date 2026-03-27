package config

import (
	"os"
	"path/filepath"
	"testing"
)

func testDefaults(t *testing.T) *Config {
	t.Helper()
	cfg, err := parseDefaults()
	if err != nil {
		t.Fatalf("parseDefaults: %v", err)
	}
	return cfg
}

func TestDefaults(t *testing.T) {
	cfg := testDefaults(t)

	if cfg.UI.ListRatio != 40 {
		t.Errorf("expected list_ratio 40, got %d", cfg.UI.ListRatio)
	}
	if cfg.UI.DetailLayout != "stacked" {
		t.Errorf("expected detail_layout stacked, got %q", cfg.UI.DetailLayout)
	}
	if len(cfg.Tabs) != 1 {
		t.Fatalf("expected 1 default tab, got %d", len(cfg.Tabs))
	}
	if cfg.Tabs[0].Name != "My Issues" {
		t.Errorf("expected default tab name 'My Issues', got %q", cfg.Tabs[0].Name)
	}
	if cfg.Cache.IssueCapacity != 50 {
		t.Errorf("expected cache.issue_capacity 50, got %d", cfg.Cache.IssueCapacity)
	}
	if cfg.Theme.Primary != "#7C3AED" {
		t.Errorf("expected theme.primary #7C3AED, got %q", cfg.Theme.Primary)
	}
	// Keybinding defaults should be populated (was a bug before embed)
	if len(cfg.Keys.Builtin.Up) == 0 {
		t.Error("expected keybinding defaults for 'up' to be populated")
	}
	if len(cfg.Keys.Builtin.Quit) == 0 {
		t.Error("expected keybinding defaults for 'quit' to be populated")
	}
}

func TestValidate_ValidConfig(t *testing.T) {
	cfg := testDefaults(t)
	if err := validate(cfg); err != nil {
		t.Errorf("defaults should be valid: %v", err)
	}
}

func TestValidate_InvalidListRatio(t *testing.T) {
	tests := []struct {
		name  string
		ratio int
	}{
		{"zero", 0},
		{"negative", -1},
		{"100", 100},
		{"200", 200},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := testDefaults(t)
			cfg.UI.ListRatio = tt.ratio
			if err := validate(cfg); err == nil {
				t.Error("expected validation error for invalid list_ratio")
			}
		})
	}
}

func TestValidate_InvalidDetailLayout(t *testing.T) {
	cfg := testDefaults(t)
	cfg.UI.DetailLayout = "invalid"
	if err := validate(cfg); err == nil {
		t.Error("expected validation error for invalid detail_layout")
	}
}

func TestValidate_TooManyTabs(t *testing.T) {
	cfg := testDefaults(t)
	cfg.Tabs = make([]TabConfig, 10)
	for i := range cfg.Tabs {
		cfg.Tabs[i] = TabConfig{Name: "Tab", JQL: "ORDER BY updated DESC"}
	}
	if err := validate(cfg); err == nil {
		t.Error("expected validation error for too many tabs")
	}
}

func TestValidate_EmptyTabsFallback(t *testing.T) {
	cfg := testDefaults(t)
	cfg.Tabs = nil
	if err := validate(cfg); err != nil {
		t.Errorf("empty tabs should fallback to default: %v", err)
	}
	if len(cfg.Tabs) != 1 {
		t.Errorf("expected 1 fallback tab, got %d", len(cfg.Tabs))
	}
}

func TestValidate_InvalidThemeColor(t *testing.T) {
	cfg := testDefaults(t)
	cfg.Theme.Primary = "not-a-color"
	if err := validate(cfg); err == nil {
		t.Error("expected validation error for invalid hex color")
	}
}

func TestValidate_ValidThemeColors(t *testing.T) {
	cfg := testDefaults(t)
	cfg.Theme.Primary = "#AABBCC"
	cfg.Theme.Error = "#ff0000"
	if err := validate(cfg); err != nil {
		t.Errorf("valid hex colors should pass: %v", err)
	}
}

func TestEditor_Fallback(t *testing.T) {
	cfg := testDefaults(t)

	// Explicit editor in config
	cfg.UI.Editor = "nvim"
	if got := cfg.Editor(); got != "nvim" {
		t.Errorf("expected nvim, got %q", got)
	}

	// Falls back to $EDITOR
	cfg.UI.Editor = ""
	t.Setenv("EDITOR", "emacs")
	if got := cfg.Editor(); got != "emacs" {
		t.Errorf("expected emacs, got %q", got)
	}

	// Falls back to vi
	t.Setenv("EDITOR", "")
	if got := cfg.Editor(); got != "vi" {
		t.Errorf("expected vi, got %q", got)
	}
}

func TestLoad_FromTOML(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("XDG_DATA_HOME", filepath.Join(dir, "data"))
	t.Setenv("XDG_CACHE_HOME", filepath.Join(dir, "cache"))

	configDir := filepath.Join(dir, "jiji")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}

	toml := `
[ui]
list_ratio = 40
detail_layout = "side-by-side"
editor = "code"

[ui.fields]
list = ["key", "summary", "status"]

[[tabs]]
name = "My Work"
jql = "assignee = currentUser() ORDER BY updated DESC"

[[tabs]]
name = "Sprint"
jql = "sprint in openSprints() ORDER BY rank ASC"

[theme]
primary = "#FF0000"

[cache]
issue_capacity = 100
`
	if err := os.WriteFile(filepath.Join(configDir, "config.toml"), []byte(toml), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if cfg.UI.ListRatio != 40 {
		t.Errorf("expected list_ratio 40, got %d", cfg.UI.ListRatio)
	}
	if cfg.UI.DetailLayout != "side-by-side" {
		t.Errorf("expected side-by-side, got %q", cfg.UI.DetailLayout)
	}
	if cfg.UI.Editor != "code" {
		t.Errorf("expected editor 'code', got %q", cfg.UI.Editor)
	}
	if len(cfg.Tabs) != 2 {
		t.Fatalf("expected 2 tabs, got %d", len(cfg.Tabs))
	}
	if cfg.Tabs[0].Name != "My Work" {
		t.Errorf("expected tab name 'My Work', got %q", cfg.Tabs[0].Name)
	}
	if cfg.Theme.Primary != "#FF0000" {
		t.Errorf("expected theme.primary #FF0000, got %q", cfg.Theme.Primary)
	}
	if cfg.Cache.IssueCapacity != 100 {
		t.Errorf("expected cache.issue_capacity 100, got %d", cfg.Cache.IssueCapacity)
	}
	// Fields from TOML that weren't set should keep defaults
	if cfg.Theme.Secondary != "#06B6D4" {
		t.Errorf("expected default theme.secondary, got %q", cfg.Theme.Secondary)
	}
}

func TestLoad_NoUserConfig(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("XDG_DATA_HOME", filepath.Join(dir, "data"))
	t.Setenv("XDG_CACHE_HOME", filepath.Join(dir, "cache"))

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	// Should return defaults
	if cfg.UI.ListRatio != 40 {
		t.Errorf("expected default list_ratio 40, got %d", cfg.UI.ListRatio)
	}

	// Keybinding defaults should be populated
	if len(cfg.Keys.Builtin.Up) != 2 || cfg.Keys.Builtin.Up[0] != "k" {
		t.Errorf("expected keybinding up=[k, up], got %v", cfg.Keys.Builtin.Up)
	}
	if len(cfg.Keys.Builtin.Quit) != 1 || cfg.Keys.Builtin.Quit[0] != "q" {
		t.Errorf("expected keybinding quit=[q], got %v", cfg.Keys.Builtin.Quit)
	}
}
