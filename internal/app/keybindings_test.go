package app

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestApp_OpenKey(t *testing.T) {
	m := setupModel(t)

	_, cmd := m.Update(tea.KeyPressMsg{Text: "o", Code: 'o'})
	if cmd == nil {
		t.Fatal("expected quit command from o key")
	}
}
