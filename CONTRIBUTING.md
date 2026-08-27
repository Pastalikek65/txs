# Contributing to txs

```bash
git clone https://github.com/Pastalikek65/txs.git
cd txs
go vet ./...
go test ./... -count=1 -timeout 30s -cover
CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o txs .
TXS_PLAIN=1 ./txs
```

PRs: fork → feat/name → test → open PR.
