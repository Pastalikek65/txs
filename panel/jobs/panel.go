package jobs

import (
	"fmt"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/Pastalikek65/txs/store"
)

func Push(s *store.Store, cmd string) (int64, error) {
	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return 0, fmt.Errorf("empty command")
	}
	return s.PushJob(cmd)
}

func RunOnce(s *store.Store, webhook string) (int, error) {
	pending, err := s.ListJobs("pending")
	if err != nil {
		return -1, err
	}
	if len(pending) == 0 {
		// retry failed with <3 retries
		failed, _ := s.ListJobs("failed")
		for _, j := range failed {
			if j.Retries < 3 {
				pending = []store.Job{j}
				break
			}
		}
		if len(pending) == 0 {
			return -1, nil
		}
	}
	job := pending[0]
	if err := s.SetJobStatus(job.ID, "running"); err != nil {
		return -1, err
	}
	cmd := exec.Command("sh", "-c", job.Cmd)
	// timeout handled by store? use 1h via RunOnce caller
	err = cmd.Run()
	if err != nil {
		_ = s.SetJobStatus(job.ID, "failed")
		_ = s.IncrementRetries(job.ID)
		if webhook != "" {
			notify(webhook, fmt.Sprintf("job %d failed: %s", job.ID, job.Cmd))
		}
		return 1, nil
	}
	_ = s.SetJobStatus(job.ID, "done")
	return 0, nil
}

func notify(webhook, msg string) {
	// minimal: try curl via exec, ignore errors
	_ = exec.Command("sh", "-c", fmt.Sprintf("curl -s -X POST -H 'Content-Type: application/json' -d '{\"text\":\"%s\"}' %q >/dev/null 2>&1 || true", msg, webhook)).Run()
}

type model struct {
	jobs []store.Job
	cursor int
}

func NewModel(jobs []store.Job) tea.Model {
	cp := make([]store.Job, len(jobs))
	copy(cp, jobs)
	return model{jobs: cp}
}

func (m model) Init() tea.Cmd { return nil }
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c": return m, tea.Quit
		case "j","down": if m.cursor < len(m.jobs)-1 { m.cursor++ }
		case "k","up": if m.cursor>0 { m.cursor-- }
		}
	}
	return m, nil
}
func (m model) View() string {
	if len(m.jobs)==0 { return "no jobs — use push" }
	var b strings.Builder
	for i,j := range m.jobs {
		prefix:="  "
		if i==m.cursor { prefix="> " }
		b.WriteString(fmt.Sprintf("%s%d [%s] %s (retry %d)\n", prefix, j.ID, j.Status, j.Cmd, j.Retries))
	}
	b.WriteString("\n[j/k] move [q] quit\n")
	return b.String()
}
