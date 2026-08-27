package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/Pastalikek65/txs/store"
)

func TestRootTab(t *testing.T) {
	s, _ := store.Open(t.TempDir()+"/db")
	defer s.Close()
	m := NewRoot(s)
	if len(m.tabs) != 5 { t.Fatalf("want 5 tabs got %d", len(m.tabs)) }
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
	m = newM.(RootModel)
	if m.active != 1 { t.Fatalf("want active 1 got %d", m.active) }
	v := m.View()
	if v == "" { t.Fatal("empty view") }
}
