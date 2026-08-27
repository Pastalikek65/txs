package files

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	dir string
	entries []string
	cursor int
}

func New(dir string) Model {
	if dir=="" { dir="." }
	entries, _ := os.ReadDir(dir)
	var names []string
	for _, e:= range entries { names = append(names, e.Name()) }
	return Model{dir: dir, entries: names}
}
func (m Model) Init() tea.Cmd { return nil }
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q","ctrl+c": return m, tea.Quit
		case "j","down": if m.cursor < len(m.entries)-1 { m.cursor++ }
		case "k","up": if m.cursor>0 { m.cursor-- }
		case "enter": {
			sel := m.entries[m.cursor]
			p := filepath.Join(m.dir, sel)
			if fi, err:=os.Stat(p); err==nil && fi.IsDir() {
				entries,_:=os.ReadDir(p)
				var names []string
				for _,e:=range entries { names=append(names, e.Name()) }
				m.dir=p
				m.entries=names
				m.cursor=0
			}
		}
		}
	}
	return m, nil
}
func (m Model) View() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(" files — %s\n", m.dir))
	for i,e := range m.entries {
		prefix:="  "
		if i==m.cursor { prefix="> " }
		b.WriteString(prefix+e+"\n")
		if i>20 { b.WriteString("  ...\n"); break }
	}
	b.WriteString("\n[j/k] nav [enter] open [q] quit\n")
	b.WriteString(fmt.Sprintf(" %d entries\n", len(m.entries)))
	_ = strings.Contains
	return b.String()
}
