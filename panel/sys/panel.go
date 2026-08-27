package sys

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct{}

func New() Model { return Model{} }
func (m Model) Init() tea.Cmd { return nil }
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if k, ok := msg.(tea.KeyMsg); ok && (k.String()=="q"||k.String()=="ctrl+c") { return m, tea.Quit }
	return m, nil
}
func (m Model) View() string {
	var b strings.Builder
	b.WriteString(" sys panel — Termux stats\n")
	// mem
	if data, err := os.ReadFile("/proc/meminfo"); err==nil {
		for _, l := range strings.Split(string(data), "\n")[:5] {
			b.WriteString("  "+l+"\n")
		}
	} else {
		b.WriteString("  no /proc/meminfo\n")
	}
	b.WriteString(fmt.Sprintf("\n  GOOS: %s  HT: %s\n", os.Getenv("GOOS"), os.Getenv("HOME")))
	b.WriteString("\n[q] quit\n")
	return b.String()
}
