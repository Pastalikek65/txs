test:
	go test ./... -count=1 -timeout 30s -cover

build:
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -X main.version=$(shell git describe --tags --always 2>/dev/null || echo dev)" -o txs .

vet:
	go vet ./...
