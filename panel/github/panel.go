package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/Pastalikek65/txs/store"
)

// ghRepo maps GitHub API response.
type ghRepo struct {
	Name      string `json:"name"`
	Stars     int    `json:"stargazers_count"`
	Language  string `json:"language"`
	UpdatedAt string `json:"updated_at"`
	HTMLURL   string `json:"html_url"`
	URL       string `json:"url"`
}

// Load fetches repos for user, handles pagination, stores via UpsertRepos, and returns sorted repos.
// It validates user, uses PathEscape, respects GH_TOKEN, follows Link pagination, and uses 15s timeout.
func Load(ctx context.Context, s *store.Store, user, sortBy string) ([]store.Repo, error) {
	if strings.TrimSpace(user) == "" {
		return nil, fmt.Errorf("invalid user: empty user")
	}
	// additional validation: user should not contain only whitespace, otherwise PathEscape would still work
	// but we treat empty after trim as invalid.

	base := os.Getenv("GITHUB_API_BASE")
	if base == "" {
		base = "https://api.github.com"
	}
	base = strings.TrimSuffix(base, "/")

	escaped := url.PathEscape(user)
	// initial URL
	nextURL := fmt.Sprintf("%s/users/%s/repos?per_page=100&sort=updated", base, escaped)

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	// ensure context has timeout if not already? Use 15s as fallback via client timeout, not context.
	// If ctx has no deadline, we still have client timeout.

	var all []store.Repo

	for nextURL != "" {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, nextURL, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("User-Agent", "txs")
		// GH_TOKEN precedence
		token := os.Getenv("GH_TOKEN")
		if token == "" {
			token = os.Getenv("GITHUB_TOKEN")
		}
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		// ensure body closed
		func() {
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				// try to read error? minimal
				err = fmt.Errorf("github api %d for %s", resp.StatusCode, nextURL)
				return
			}
			var ghRepos []ghRepo
			if decErr := json.NewDecoder(resp.Body).Decode(&ghRepos); decErr != nil {
				err = decErr
				return
			}
			for _, gr := range ghRepos {
				u := gr.HTMLURL
				if u == "" {
					u = gr.URL
				}
				all = append(all, store.Repo{
					Name:      gr.Name,
					Stars:     gr.Stars,
					Language:  gr.Language,
					UpdatedAt: gr.UpdatedAt,
					URL:       u,
				})
			}
			// parse Link header for next
			link := resp.Header.Get("Link")
			nextURL = parseNextLink(link)
		}()
		if err != nil {
			return nil, err
		}
		if respStatusNotOK(nextURL, err) {
			// already handled above, but keep
		}
	}

	// upsert
	if err := s.UpsertRepos(all); err != nil {
		return nil, err
	}
	return s.ListRepos(sortBy)
}

func respStatusNotOK(nextURL string, err error) bool {
	return false
}

// parseNextLink extracts URL with rel="next" from Link header.
func parseNextLink(link string) string {
	if link == "" {
		return ""
	}
	// Link: <url>; rel="next", <url>; rel="last"
	parts := strings.Split(link, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if strings.Contains(p, `rel="next"`) {
			// find <...>
			start := strings.Index(p, "<")
			end := strings.Index(p, ">")
			if start != -1 && end != -1 && end > start {
				return p[start+1 : end]
			}
		}
	}
	return ""
}

// Model implements tea.Model for github panel.

type model struct {
	repos  []store.Repo
	cursor int
}

func NewModel(repos []store.Repo) tea.Model {
	// copy to avoid aliasing
	cp := make([]store.Repo, len(repos))
	copy(cp, repos)
	return model{repos: cp}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "j", "down":
			if m.cursor < len(m.repos)-1 {
				m.cursor++
			}
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	if len(m.repos) == 0 {
		return "no repos"
	}
	var b strings.Builder
	for i, r := range m.repos {
		prefix := "  "
		if i == m.cursor {
			prefix = "> "
		}
		// rune-aware truncate would be used here if needed
		b.WriteString(fmt.Sprintf("%s%s ★%d %s %s\n", prefix, r.Name, r.Stars, r.Language, r.URL))
	}
	b.WriteString("\n[j/k] move  [q] quit  [/] filter  [s/n/f/u] sort\n")
	return b.String()
}
