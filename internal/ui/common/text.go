package common

import (
	"strings"
	"unicode"

	"github.com/charmbracelet/x/ansi"
	"github.com/mattn/go-runewidth"
)

// Truncate cuts a string to maxLen, appending "…" if truncated.
func Truncate(s string, maxLen int) string {
	trimed := strings.TrimRightFunc(s, unicode.IsSpace)
	if maxLen <= 0 {
		return ""
	}
	if runewidth.StringWidth(trimed) <= maxLen {
		return trimed
	}
	if maxLen <= 1 {
		return "…"
	}

	return runewidth.Truncate(trimed, maxLen, "…")
}

// ReplaceAt replaces visible characters in s at visual position start with replacement.
// ANSI escape sequences are preserved.
func ReplaceAt(s string, start int, replacement string) string {
	end := start + ansi.StringWidth(replacement)
	return ansi.Truncate(s, start, "") + replacement + ansi.TruncateLeft(s, end, "")
}
