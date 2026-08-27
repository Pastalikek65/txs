package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	ghpanel "github.com/Pastalikek65/txs/panel/github"
	jobspanel "github.com/Pastalikek65/txs/panel/jobs"
	dbpanel "github.com/Pastalikek65/txs/panel/db"
	syspanel "github.com/Pastalikek65/txs/panel/sys"
	filepanel "github.com/Pastalikek65/txs/panel/files"
	"github.com/Pastalikek65/txs/store"
)

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("99")).Padding(0,1)
	faintStyle = lipgloss.NewStyle().Faint(true)
	activeStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
)

type RootModel struct {
	tabs []string
	active int
	models []tea.Model
	store *store.Store
}

func NewRoot(s *store.Store) RootModel {
	repos,_ := s.ListRepos("stars")
	jobs,_ := s.ListJobs("")
	return RootModel{
		tabs: []string{"github","jobs","db","sys","files"},
		active: 0,
		models: []tea.Model{
			ghpanel.NewModel(repos),
			jobspanel.NewModel(jobs),
			dbpanel.New(s),
			syspanel.New(),
			filepanel.New("."),
		},
		store: s,
	}
}

func (m RootModel) Init() tea.Cmd { return nil }

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if k, ok := msg.(tea.KeyMsg); ok {
		switch k.String() {
		case "q","ctrl+c": return m, tea.Quit
		case "tab","l","right": m.active = (m.active+1)%len(m.tabs); return m, nil
		case "shift+tab","h","left": m.active = (m.active-1+len(m.tabs))%len(m.tabs); return m, nil
		}
	}
	// delegate to active panel
	newM, cmd := m.models[m.active].Update(msg)
	m.models[m.active] = newM
	return m, cmd
}

func (m RootModel) View() string {
	var b strings.Builder
	// tab bar
	for i,t := range m.tabs {
		if i==m.active { b.WriteString(activeStyle.Render("["+t+"]")+" ") } else { b.WriteString(faintStyle.Render(t)+" ") }
	}
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", 60)+"\n")
	b.WriteString(titleStyle.Render(" txs — Termux Station ")+"\n")
	b.WriteString(m.models[m.active].View())
	b.WriteString("\n"+faintStyle.Render("tab: switch  q: quit")+"\n")
	b.WriteString(fmt.Sprintf(" txs %d tabs\n", len(m.tabs)))
	return b.String()
}
