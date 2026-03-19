package config

// ThemeConfig holds hex color values for the UI theme.
type ThemeConfig struct {
	Primary       string `koanf:"primary"`
	Secondary     string `koanf:"secondary"`
	Border        string `koanf:"border"`
	BorderFocused string `koanf:"border_focused"`
	Text          string `koanf:"text"`
	TextDim       string `koanf:"text_dim"`
	Success       string `koanf:"success"`
	Warning       string `koanf:"warning"`
	Error         string `koanf:"error"`
}

// DefaultTheme returns the default color theme from ADR-003.
func DefaultTheme() ThemeConfig {
	return ThemeConfig{
		Primary:       "#7C3AED",
		Secondary:     "#06B6D4",
		Border:        "#404040",
		BorderFocused: "#7C3AED",
		Text:          "#E4E4E7",
		TextDim:       "#71717A",
		Success:       "#22C55E",
		Warning:       "#F59E0B",
		Error:         "#EF4444",
	}
}
