// Package config handles loading and validating application configuration.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// CLI holds Jira connection settings parsed from CLI flags / env vars.
// Parsed by kong in main.go.
type CLI struct {
	Host    string `kong:"help='Jira instance URL',env='JIJI_JIRA_HOST',short='H'"`
	Email   string `kong:"help='Jira user email',env='JIJI_JIRA_EMAIL',short='e'"`
	Token   string `kong:"help='Jira API token',env='JIJI_JIRA_TOKEN',short='t'"`
	Version bool   `kong:"help='Print version and exit',short='v'"`
}

// ValidateConnection checks that all required Jira connection fields are set.
func (c *CLI) ValidateConnection() error {
	if c.Host == "" {
		return fmt.Errorf("--host is required (or set JIJI_JIRA_HOST)")
	}
	if c.Email == "" {
		return fmt.Errorf("--email is required (or set JIJI_JIRA_EMAIL)")
	}
	if c.Token == "" {
		return fmt.Errorf("--token is required (or set JIJI_JIRA_TOKEN)")
	}
	return nil
}

// JiraConnection holds the resolved Jira connection info.
type JiraConnection struct {
	Host  string
	Email string
	Token string
}

// Config is the top-level application configuration.
type Config struct {
	Jira  JiraConnection // from CLI/env, not from file
	UI    UIConfig       `koanf:"ui"`
	Tabs  []TabConfig    `koanf:"tabs"`
	Keys  KeyConfig      `koanf:"keybindings"`
	Theme ThemeConfig    `koanf:"theme"`
	Cache CacheConfig    `koanf:"cache"`
}

// UIConfig holds user-interface preferences.
type UIConfig struct {
	Theme        string       `koanf:"theme"`
	ListRatio    int          `koanf:"list_ratio"`
	DetailLayout string       `koanf:"detail_layout"`
	Editor       string       `koanf:"editor"`
	Fields       FieldsConfig `koanf:"fields"`
}

// FieldsConfig controls which fields appear in the issue list.
type FieldsConfig struct {
	List []string `koanf:"list"`
}

// TabConfig defines a named tab with a JQL query.
type TabConfig struct {
	Name string `koanf:"name"`
	JQL  string `koanf:"jql"`
}

// KeyConfig holds user-provided keybinding overrides.
// Each key is a list of bindings (e.g. ["k", "up"]).
type KeyConfig struct {
	Up         []string `koanf:"up"`
	Down       []string `koanf:"down"`
	TabNext    []string `koanf:"tab_next"`
	TabPrev    []string `koanf:"tab_prev"`
	PaneSwitch []string `koanf:"pane_switch"`
	Top        []string `koanf:"top"`
	Bottom     []string `koanf:"bottom"`
	Confirm    []string `koanf:"confirm"`
	FocusJQL   []string `koanf:"focus_jql"`
	Cancel     []string `koanf:"cancel"`
	Quit       []string `koanf:"quit"`
	Help       []string `koanf:"help"`
	Transition []string `koanf:"transition"`
	Comment    []string `koanf:"comment"`
	Labels     []string `koanf:"labels"`
	Summary    []string `koanf:"summary"`
	Edit       []string `koanf:"edit"`
	Refresh    []string `koanf:"refresh"`
}

// CacheConfig controls in-memory caching behavior.
type CacheConfig struct {
	IssueCapacity   int `koanf:"issue_capacity"`
	CommentCapacity int `koanf:"comment_capacity"`
	PrefetchCount   int `koanf:"prefetch_count"`
}

// xdgDir returns an XDG base directory, reading env vars at call time.
func xdgDir(envVar, defaultSuffix string) string {
	if v := os.Getenv(envVar); v != "" {
		return v
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, defaultSuffix)
}

// configDir returns the jiji config directory path.
func configDir() string {
	return filepath.Join(xdgDir("XDG_CONFIG_HOME", ".config"), "jiji")
}

// configFilePath returns the path to config.toml.
func configFilePath() string {
	return filepath.Join(configDir(), "config.toml")
}

// Load reads the config file and returns a Config with defaults applied.
// If no config file exists, a default one is created.
func Load() (*Config, error) {
	cfg := defaults()

	if err := ensureXDGDirs(); err != nil {
		return nil, fmt.Errorf("creating XDG directories: %w", err)
	}

	path := configFilePath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := writeDefaultConfig(path); err != nil {
			return nil, fmt.Errorf("writing default config: %w", err)
		}
		return cfg, nil
	}

	k := koanf.New(".")
	if err := k.Load(file.Provider(path), toml.Parser()); err != nil {
		return nil, fmt.Errorf("loading config file %s: %w", path, err)
	}
	if err := k.Unmarshal("", cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if err := validate(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// defaults returns a Config populated with default values.
func defaults() *Config {
	return &Config{
		UI: UIConfig{
			Theme:        "default",
			ListRatio:    30,
			DetailLayout: "stacked",
			Editor:       "",
			Fields: FieldsConfig{
				List: []string{"key", "summary", "assignee", "priority"},
			},
		},
		Tabs: []TabConfig{
			{Name: "All Issues", JQL: "ORDER BY updated DESC"},
		},
		Theme: DefaultTheme(),
		Cache: CacheConfig{
			IssueCapacity:   50,
			CommentCapacity: 20,
			PrefetchCount:   10,
		},
	}
}

// Editor returns the configured editor, falling back to $EDITOR, then "vi".
func (c *Config) Editor() string {
	if c.UI.Editor != "" {
		return c.UI.Editor
	}
	if e := os.Getenv("EDITOR"); e != "" {
		return e
	}
	return "vi"
}

var hexColorRegex = regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)

// validate checks config values and returns a descriptive error on failure.
func validate(cfg *Config) error {
	if cfg.UI.ListRatio < 1 || cfg.UI.ListRatio > 99 {
		return fmt.Errorf("config: ui.list_ratio must be between 1 and 99, got %d", cfg.UI.ListRatio)
	}

	validLayouts := map[string]bool{"stacked": true, "side-by-side": true}
	if !validLayouts[cfg.UI.DetailLayout] {
		return fmt.Errorf("config: ui.detail_layout must be \"stacked\" or \"side-by-side\", got %q", cfg.UI.DetailLayout)
	}

	if len(cfg.Tabs) == 0 {
		cfg.Tabs = []TabConfig{{Name: "All Issues", JQL: "ORDER BY updated DESC"}}
	}
	if len(cfg.Tabs) > 9 {
		return fmt.Errorf("config: at most 9 tabs allowed, got %d", len(cfg.Tabs))
	}

	themeColors := map[string]string{
		"theme.primary":        cfg.Theme.Primary,
		"theme.secondary":      cfg.Theme.Secondary,
		"theme.border":         cfg.Theme.Border,
		"theme.border_focused": cfg.Theme.BorderFocused,
		"theme.text":           cfg.Theme.Text,
		"theme.text_dim":       cfg.Theme.TextDim,
		"theme.success":        cfg.Theme.Success,
		"theme.warning":        cfg.Theme.Warning,
		"theme.error":          cfg.Theme.Error,
	}
	for key, val := range themeColors {
		if val != "" && !hexColorRegex.MatchString(val) {
			return fmt.Errorf("config: %s must be a hex color (#RRGGBB), got %q", key, val)
		}
	}

	return nil
}

// ensureXDGDirs creates the XDG directories for jiji.
func ensureXDGDirs() error {
	dirs := []string{
		filepath.Join(xdgDir("XDG_CONFIG_HOME", ".config"), "jiji"),
		filepath.Join(xdgDir("XDG_DATA_HOME", filepath.Join(".local", "share")), "jiji"),
		filepath.Join(xdgDir("XDG_CACHE_HOME", ".cache"), "jiji"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0o755); err != nil {
			return err
		}
	}
	return nil
}

// writeDefaultConfig writes a minimal config file with commented-out defaults.
func writeDefaultConfig(path string) error {
	content := strings.TrimSpace(`
# Jiji configuration file
# Uncomment and modify settings as needed.

# [ui]
# theme = "default"
# list_ratio = 30
# detail_layout = "stacked"   # "stacked" or "side-by-side"
# editor = ""                  # Falls back to $EDITOR, then "vi"

# [ui.fields]
# list = ["key", "summary", "assignee", "priority"]

# [[tabs]]
# name = "All Issues"
# jql = "ORDER BY updated DESC"

# [keybindings]
# up = ["k", "up"]
# down = ["j", "down"]
# tab_next = ["l", "right"]
# tab_prev = ["h", "left"]
# pane_switch = ["tab"]

# [theme]
# primary = "#7C3AED"
# secondary = "#06B6D4"
# border = "#404040"
# border_focused = "#7C3AED"
# text = "#E4E4E7"
# text_dim = "#71717A"
# success = "#22C55E"
# warning = "#F59E0B"
# error = "#EF4444"

# [cache]
# issue_capacity = 50
# comment_capacity = 20
# prefetch_count = 10
`) + "\n"
	return os.WriteFile(path, []byte(content), 0o644)
}
