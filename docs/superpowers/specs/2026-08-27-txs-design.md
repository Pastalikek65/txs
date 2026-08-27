# Termux Station (txs) — Design Spec

> **Status:** Approved 2026-08-27 (yolo) | **For:** `Pastalikek65/txs` | **Constraints:** Termux proot Debian, aarch64, `CGO_ENABLED=0`, single binary <20M, <50MB RAM, offline

## 1. Goal & Success Criteria

**Goal:** Single-binary TUI dashboard for Termux that shows GitHub repos, job queue, SQLite DBs, system stats, and files in 5 panels. Built on phone, `go build` 4s.

**Success (verifiable):**
- `CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o txs .` <10s, binary <20M, `go test ./...` 70%+ in <30s, no CGO
- `txs --help` shows 5 panels, `txs --version` injected via `git describe`
- `txs` TUI: `tab` switches panels, `q` quits, `GH_TUI_PLAIN=1 txs` prints plain without bubbletea
- `store.Open` creates `~/.cache/txs/txs.db` WAL `0700`, 5 tables, survives reboot
- `gh extension install Pastalikek65/txs` works (binary `txs` at repo root, goreleaser `android-arm64`+`linux-arm64`)
- `README` humanized (first-person), `demo.gif` 4-frame, `topics:8`

## 2. Architecture

Single Go binary, `cobra` + `bubbletea` + `lipgloss`, `modernc.org/sqlite` WAL. 5 isolated panels, each with `Load(ctx)`, `Init()`, `Update(msg)`, `View()`, share nothing except `store.Store`.

```
txs (cmd/root.go)
 ├─ store/store.go         — Open, 5 tables (repos, jobs, dbs, sys, files), WAL, 0700
 ├─ panel/github/panel.go  — fetchRepos with Link pagination, PathEscape, 10s timeout → store
 ├─ panel/jobs/panel.go    — push/list/retry, worker subprocess 1h timeout, job_<id>.log, webhook
 ├─ panel/db/panel.go      — ListTables/Query/WalSize, preview 20 rows
 ├─ panel/sys/panel.go     — /proc/meminfo, /proc/stat, df, battery via termux-battery-status
 └─ panel/files/panel.go   — ls -R, cat, rm (proot safe, no rm -rf /)
tui/root.go                — tab router, 5 panels as tea.Model, faint help line
```

**Isolation:** Each panel can be understood via `Load` signature alone. Changing `github` fetch does not affect `db` View. Store is the only shared dep, via `Store` interface.

## 3. Components & Interfaces

### 3.1 store.Store
```go
func Open(path string) (*Store, error) // WAL, 0700, 5 tables, non-destructive migration
func (s *Store) UpsertRepos([]Repo) error
func (s *Store) ListRepos(sortBy string) ([]Repo, error) // sortBy validated inside, not concat
func (s *Store) PushJob(cmd string) (int64, error)
func (s *Store) ListJobs(status string) ([]Job, error)
func (s *Store) Close() error
```
`Repo` mirrors `cache.Repo` but with `url TEXT PRIMARY KEY`. `Job` mirrors `pytermq` `jobs` table.

### 3.2 panel/github
```go
func Load(ctx context.Context, store *store.Store, user, sortBy string) ([]Repo, error) // respects GH_TOKEN, 15s ctx
func NewModel(repos []Repo) tea.Model // sort s/n/f/u, filter /, j/k, q
```

### 3.3 panel/jobs
```go
func Push(cmd string) (int64, error)
func RunOnce(storePath string, webhook string) (int, error) // 1 job, retry up to 3, ntfy POST
```

### 3.4 panel/db
Already spec'd in `sqlite-tui`: `ListTables() ([]string,error)`, `Query(table,limit)` etc., plus `WalSize()`.

### 3.5 panel/sys & files
- `sys`: `ReadMem()`, `ReadCPU()`, `ReadDisk()`, `ReadBattery()` — parses `/proc`, `df`, `termux-battery-status` JSON
- `files`: `List(dir) ([]Entry,error)`, `Cat(path) (string,error)` — `filepath.Join` safe

## 4. Data Flow & Error Handling

1. `txs` start → `store.Open("~/.cache/txs/txs.db")` → 5 panels `Load` in parallel via `errgroup`
2. `github` → `GET /users/:user/repos?per_page=100&sort=updated` with `Link rel="next"` pagination → `UpsertRepos` → `store`
3. `jobs` → `sqlite` pending → `subprocess.run(shell=True)` 1h → log → `done`/`failed`+`retries=3` → webhook `POST`
4. `db` → `sqlite_master` → `Query LIMIT 20` → `WalSize` via `stat -wal`
5. Errors: `ErrRateLimited` → stale cache + `⚠️` on stderr; `ErrNetwork` → offline cache; `invalid user` → `fmt.Errorf` before fetch; `403` non-limit → generic, not rate-limit

## 5. Phone Constraints

- `CGO_ENABLED=0`, `modernc.org/sqlite` pure Go, `trimpath`, `0700` cache, `android-arm64` goreleaser (gh reports `android-arm64` on Termux)
- `go build` <10s on 8-core aarch64, `go test` <30s, binary <20M, RAM <50M (5 panels × 20 rows)
- No Docker, no daemon, single `txs.db`, `term.IsTerminal` plain fallback

## 6. Testing

- `store`: `TestOpenAndUpsertRepos`, `TestPushAndListJobs`, `TestMigration`
- `panel/github`: `TestFetchMockServer` (httptest + Link), `TestRateLimit429`, `TestInvalidUser`
- `panel/db`: `TestListTables`, `TestQuery`, `TestWalSize`
- `panel/sys`: `TestReadMem` with fixture `/proc/meminfo` sample
- `tui/root`: `TestTabSwitch`, `TestFilter`, `TestQuit`
- `integration`: `TestLoadAllPanels` with temp `XDG_CACHE_HOME` + mock GitHub

## 7. Distribution

- `gh extension` — binary `txs` at repo root, `goreleaser` `goos: [linux,android,darwin]` `goarch: [amd64,arm64]` `android/amd64` ignored, `ldflags -X main.version={{.Version}}`, `checksum: checksums.txt`
- Release footer: `gh extension install` + `go install` + `CGO_ENABLED=0 go build`
- `README` first-person, `demo.gif` 4-frame (900x700), `CHANGELOG` keep-a-changelog, `topics:8`

## 8. Non-Goals (YAGNI)

- No `termtunnel` WebRTC share in v0.1 (defer to v0.2)
- No `rclone` sync, no `docker`, no `systemd`, no private repos (`/user/repos`)
- No `sindresorhus/awesome` list (single tool, not list)
