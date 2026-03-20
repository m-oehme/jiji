// Package common provides shared types for all UI components.
package common

import (
	"github.com/m-oehme/jiji/internal/config"
	"github.com/m-oehme/jiji/internal/ui/styles"
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
