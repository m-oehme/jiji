package config

// DefaultKeyConfig returns the default keybinding configuration per ADR-005.
func DefaultKeyConfig() KeyConfig {
	return KeyConfig{
		Up:         []string{"k", "up"},
		Down:       []string{"j", "down"},
		TabNext:    []string{"l", "right"},
		TabPrev:    []string{"h", "left"},
		PaneSwitch: []string{"tab"},
		Top:        []string{"g"},
		Bottom:     []string{"G"},
		Confirm:    []string{"enter"},
		FocusJQL:   []string{"/"},
		Cancel:     []string{"esc"},
		Quit:       []string{"q"},
		Help:       []string{"?"},
		Transition: []string{"t"},
		Comment:    []string{"c"},
		Labels:     []string{"L"},
		Summary:    []string{"s"},
		Edit:       []string{"e"},
		Refresh:    []string{"r"},
	}
}
