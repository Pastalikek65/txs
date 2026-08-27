package db

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/Pastalikek65/txs/store"
)

type Model struct {
	store *store.Store
	tables []string
	cursor int
	rows []map[string]string
}

func New(s *store.Store) Model {
	tables, _ := s.ListTables()
	m := Model{store: s, tables: tables}
	if len(tables)>0 { m.load(tables[0]) }
	return m
}

func (m *Model) load(table string) {
	// use store's db via ListTables? For txs.db, we can query via store
	// simple: try to get rows via direct query
	m.rows = nil
	// fallback: list repos/jobs as tables
}

func (m Model) Tables() []string { return m.tables }
func (m Model) Init() tea.Cmd { return nil }
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q","ctrl+c": return m, tea.Quit
		case "j","down": if m.cursor < len(m.tables)-1 { m.cursor++; m.load(m.tables[m.cursor]) }
		case "k","up": if m.cursor>0 { m.cursor-- ; m.load(m.tables[m.cursor]) }
		}
	}
	return m, nil
}
func (m Model) View() string {
	var b strings.Builder
	b.WriteString(" db panel — tables\n")
	for i,t := range m.tables {
		prefix:="  "
		if i==m.cursor { prefix="> " }
		b.WriteString(fmt.Sprintf("%s%s\n", prefix, t))
	}
	b.WriteString("\n[j/k] nav [q] quit\n")
	return b.String()
}
