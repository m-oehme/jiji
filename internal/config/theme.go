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
