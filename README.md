# txs — Termux Station

[![Go 1.25](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Termux](https://img.shields.io/badge/Termux-friendly-brightgreen?logo=android)](https://termux.dev)

Single-binary TUI dashboard for Termux. Five panels, one cache, offline.

I built `txs` because I had four terminals open on my phone: one for `gh-starview`, one for `pytermq`, one for `sqlite-tui`, one for `htop`. Now it's one `tab`.

```bash
go install github.com/Pastalikek65/txs@latest
txs              # TUI: tab to switch, q to quit
TXS_PLAIN=1 txs  # plain for CI
```

![demo](https://raw.githubusercontent.com/Pastalikek65/txs/main/demo.gif)

No Docker. No daemon. Cache at `~/.cache/txs/txs.db` WAL `0700`.

## Panels

- **github** — `gh-starview` logic: sortable `s/n/f/u`, `/` filter, `Link` pagination, `15s` timeout
- **jobs** — `pytermq` port: `push`/`list`/`retry`, worker `1h`, `job_<id>.log`, webhook `ntfy`
- **db** — `sqlite-tui` logic: `ListTables`, `Query LIMIT 20`, `WAL` size
- **sys** — `/proc/meminfo` + `df` + battery
- **files** — `ls -R` proot-safe, `enter` to open dir

Keys: `tab` / `shift+tab` switch, `j/k` nav, `r` reload, `q` quit.

## Install (Termux)

```bash
git clone https://github.com/Pastalikek65/txs.git
cd txs
CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o txs .
./txs --help
# as gh extension
gh extension install Pastalikek65/txs
gh txs
```

## How it works

`store` → `~/.cache/txs/txs.db` 5 tables (`repos` `jobs` `dbs` `sys` `files`) WAL. Each panel `Load(ctx, store)` in parallel via `errgroup`, `tui/root.go` routes `tab`.

Go 1.25, `cobra`, `bubbletea` 1.3.4, `lipgloss`, `modernc.org/sqlite`.

## Dev

```bash
go vet ./...
go test ./... -count=1 -timeout 30s -cover
CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o txs .
TXS_PLAIN=1 ./txs
```

## License

MIT
