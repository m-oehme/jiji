// Package common provides shared types for all UI components.
package common

import (
	"unicode/utf8"

	"github.com/m-oehme/jiji/internal/config"
	"github.com/m-oehme/jiji/internal/ui/styles"
	"github.com/mattn/go-runewidth"
)

// Component extends tea.Model with dimension management.
// Every UI component in jiji implements this interface (ADR-002).
type Component interface {
	SetSize(width, height int)
}

// Common holds shared state passed to every component.
// Components embed a pointer to Common for shared access to styles and keymap.
type Common struct {
	Width, Height int
	Styles        *styles.Styles
	Keys          *config.KeyConfig
	Focused       bool
}

// truncate cuts a string to maxLen, appending "…" if truncated.
func Truncate(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	if utf8.RuneCountInString(s) <= maxLen {
		return s
	}
	if maxLen <= 1 {
		return "…"
	}

	return runewidth.Truncate(s, maxLen, "…")
}
