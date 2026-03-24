// Package common provides shared types for all UI components.
package common

import (
	"github.com/m-oehme/jiji/internal/ui/styles"
)

// Component extends tea.Model with dimension management.
// Every UI component in jiji implements this interface (ADR-002).
type Component interface {
	SetSize(width, height int)
}

// Common holds mutable per-pane UI state (dimensions, focus, styles).
type Common struct {
	Width, Height int
	Styles        *styles.Styles
	Focused       bool
}
