VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "0.1.0")
LDFLAGS := -s -w -X main.version=$(VERSION)
.PHONY: test vet cover build clean
test:
	go test ./... -count=1 -timeout 30s -cover
vet:
	go vet ./...
cover:
	go test ./... -coverprofile=coverage.out -covermode=atomic && go tool cover -func=coverage.out | tail -5
build:
	CGO_ENABLED=0 go build -trimpath -ldflags="$(LDFLAGS)" -o txs .
clean:
	rm -f txs coverage.out
