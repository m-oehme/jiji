package common

import "testing"

func TestTruncate(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		maxLen int
		want   string
	}{
		// ASCII
		{"ascii fits", "hello", 10, "hello"},
		{"ascii exact", "hello", 5, "hello"},
		{"ascii truncated", "hello world", 5, "hell…"},
		{"ascii one char", "hi", 1, "…"},
		{"ascii zero", "hi", 0, ""},
		{"empty string", "", 5, ""},
		{"negative maxLen", "hello", -1, ""},
		{"trailing spaces trimmed", "hello   ", 5, "hello"},
		{"trailing spaces trimmed then truncated", "hello world   ", 5, "hell…"},

		// Multi-byte UTF-8 (single-width)
		// 2-byte: ñ (U+00F1), ö (U+00F6)
		{"2byte fits", "señor", 10, "señor"},
		{"2byte exact width", "señor", 5, "señor"},
		{"2byte truncated", "señor amigo", 5, "seño…"},

		// 3-byte: 日 (U+65E5) — rendered as double-width
		{"cjk fits", "日本語", 10, "日本語"},
		{"cjk exact width", "日本語", 6, "日本語"},
		{"cjk truncated", "日本語テスト", 7, "日本語…"},
		{"cjk truncated no room for wide char", "日本語", 5, "日本…"},
		{"cjk one col", "日", 1, "…"},

		// 4-byte: 𝕳 (U+1D573) — single-width mathematical letter
		{"4byte fits", "𝕳ello", 10, "𝕳ello"},
		{"4byte exact", "𝕳ello", 5, "𝕳ello"},
		{"4byte truncated", "𝕳ello world", 5, "𝕳ell…"},

		// Emoji — double-width
		{"emoji fits", "🎉done", 10, "🎉done"},
		{"emoji exact width", "🎉done", 6, "🎉done"},
		{"emoji truncated", "🎉hello world", 5, "🎉he…"},
		{"emoji only truncated", "🎉🎊🎈", 5, "🎉🎊…"},
		{"emoji only exact", "🎉🎊🎈", 6, "🎉🎊🎈"},
		{"emoji at boundary no room", "🎉🎊🎈", 3, "🎉…"},
		{"emoji maxLen 2 no room with ellipsis", "🎉🎊", 2, "…"},
		{"emoji maxLen 1", "🎉", 1, "…"},

		// ZWJ sequences (family emoji) — rendered double-width
		{"zwj emoji fits", "👨‍👩‍👧ok", 10, "👨‍👩‍👧ok"},

		// Mixed scripts
		{"mixed ascii cjk", "hi日本", 6, "hi日本"},
		{"mixed ascii cjk truncated", "hi日本語", 6, "hi日…"},
		{"mixed emoji ascii", "🚀launch", 8, "🚀launch"},
		{"mixed emoji ascii truncated", "🚀launch pad", 8, "🚀launc…"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Truncate(tt.input, tt.maxLen)
			if got != tt.want {
				t.Errorf("Truncate(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.want)
			}
		})
	}
}

func TestReplaceAt(t *testing.T) {
	tests := []struct {
		name        string
		s           string
		start       int
		replacement string
		want        string
	}{
		// ASCII
		{"ascii middle", "hello world", 5, "-", "hello-world"},
		{"ascii start", "hello", 0, "HE", "HEllo"},
		{"ascii end", "hello", 4, "O", "hellO"},

		// Multi-byte UTF-8 (single-width)
		{"2byte replace in middle", "señor", 2, "N", "seNor"},
		{"2byte replaced by 2byte", "señor", 2, "ñ", "señor"},

		// CJK double-width
		{"cjk replace at start", "日本語", 0, "月", "月本語"},
		{"cjk replace middle", "日本語", 2, "月", "日月語"},
		{"cjk ascii into wide slot", "日本語", 0, "AB", "AB本語"},

		// Emoji double-width
		{"emoji replace at start", "🎉🎊🎈", 0, "🚀", "🚀🎊🎈"},
		{"emoji replace middle", "🎉🎊🎈", 2, "🚀", "🎉🚀🎈"},
		{"emoji replace with ascii", "🎉🎊🎈", 0, "AB", "AB🎊🎈"},

		// ANSI escape sequences preserved (library re-emits sequences at cut boundaries)
		{"ansi preserved", "\x1b[31mhello\x1b[0m world", 5, "-", "\x1b[31mhello\x1b[0m-\x1b[31m\x1b[0mworld"},
		{"ansi at start", "\x1b[1mbold\x1b[0m text", 0, "BO", "\x1b[1m\x1b[0mBO\x1b[1mld\x1b[0m text"},

		// Mixed
		{"mixed emoji and ascii", "hi🎉bye", 2, "🚀", "hi🚀bye"},
		{"mixed cjk and ascii", "hi日bye", 2, "本", "hi本bye"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ReplaceAt(tt.s, tt.start, tt.replacement)
			if got != tt.want {
				t.Errorf("ReplaceAt(%q, %d, %q) = %q, want %q",
					tt.s, tt.start, tt.replacement, got, tt.want)
			}
		})
	}
}
