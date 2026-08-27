package db

import (
	"path/filepath"
	"testing"

	"github.com/Pastalikek65/txs/store"
)

func TestDBPanel(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.db")
	s, _ := store.Open(path)
	defer s.Close()
	// create a table via store's db
	s.UpsertRepos([]store.Repo{{Name: "a", Stars: 1}})
	m := New(s)
	if len(m.Tables()) == 0 {
		t.Fatal("no tables")
	}
	v := m.View()
	if v == "" { t.Fatal("empty view") }
}
