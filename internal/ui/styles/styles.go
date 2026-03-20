// Package styles converts ThemeConfig into lipgloss styles.
package styles

import (
	"github.com/m-oehme/jiji/internal/config"
	lipgloss "charm.land/lipgloss/v2"
)

// Styles holds all precomputed lipgloss styles for the application.
type Styles struct {
	Border        lipgloss.Style
	BorderFocused lipgloss.Style
	TabActive     lipgloss.Style
	TabInactive   lipgloss.Style
	StatusBar     lipgloss.Style
	IssueSelected lipgloss.Style
	IssueCurrent  lipgloss.Style
	MetadataKey   lipgloss.Style
	MetadataValue lipgloss.Style
	Heading       lipgloss.Style
	Dimmed        lipgloss.Style
	Error         lipgloss.Style
	Success       lipgloss.Style
	Warning       lipgloss.Style
}

// NewStyles creates Styles from a ThemeConfig, converting hex colors to lipgloss styles.
func NewStyles(theme config.ThemeConfig) *Styles {
	primary := lipgloss.Color(theme.Primary)
	secondary := lipgloss.Color(theme.Secondary)
	border := lipgloss.Color(theme.Border)
	borderFocused := lipgloss.Color(theme.BorderFocused)
	text := lipgloss.Color(theme.Text)
	textDim := lipgloss.Color(theme.TextDim)
	errColor := lipgloss.Color(theme.Error)
	successColor := lipgloss.Color(theme.Success)
	warningColor := lipgloss.Color(theme.Warning)

	return &Styles{
		Border: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(border),

		BorderFocused: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderFocused),

		TabActive: lipgloss.NewStyle().
			Bold(true).
			Foreground(primary).
			PaddingRight(2),

		TabInactive: lipgloss.NewStyle().
			Foreground(textDim).
			PaddingRight(2),

		StatusBar: lipgloss.NewStyle().
			Foreground(text).
			Background(lipgloss.Color("#1a1a2e")).
			PaddingLeft(1).
			PaddingRight(1),

		IssueSelected: lipgloss.NewStyle().
			Background(primary).
			Foreground(lipgloss.Color("#000000")).
			Bold(true),

		IssueCurrent: lipgloss.NewStyle().
			Foreground(secondary),

		MetadataKey: lipgloss.NewStyle().
			Foreground(textDim).
			Bold(true),

		MetadataValue: lipgloss.NewStyle().
			Foreground(text),

		Heading: lipgloss.NewStyle().
			Bold(true).
			Foreground(primary),

		Dimmed: lipgloss.NewStyle().
			Foreground(textDim),

		Error: lipgloss.NewStyle().
			Foreground(errColor).
			Bold(true),

		Success: lipgloss.NewStyle().
			Foreground(successColor),

		Warning: lipgloss.NewStyle().
			Foreground(warningColor),
	}
}
