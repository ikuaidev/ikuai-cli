BINARY = ikuai-cli
CMD    = ./cmd/ikuai-cli
DIST   = dist
VERSION ?= dev
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE    ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS = -X github.com/ikuaidev/ikuai-cli/internal/buildinfo.Version=$(VERSION) -X github.com/ikuaidev/ikuai-cli/internal/buildinfo.Commit=$(COMMIT) -X github.com/ikuaidev/ikuai-cli/internal/buildinfo.Date=$(DATE)

all: test build

fmt:
	gofmt -w ./cmd ./internal

lint:
	golangci-lint run

test:
	go test ./...

build:
	go build -ldflags="$(LDFLAGS)" -o $(BINARY) $(CMD)

linux-amd64:
	@mkdir -p $(DIST)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w $(LDFLAGS)" -o $(DIST)/$(BINARY)-linux-amd64 $(CMD)

linux-arm64:
	@mkdir -p $(DIST)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w $(LDFLAGS)" -o $(DIST)/$(BINARY)-linux-arm64 $(CMD)

linux-armv7:
	@mkdir -p $(DIST)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -ldflags="-s -w $(LDFLAGS)" -o $(DIST)/$(BINARY)-linux-armv7 $(CMD)

darwin-arm64:
	@mkdir -p $(DIST)
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w $(LDFLAGS)" -o $(DIST)/$(BINARY)-darwin-arm64 $(CMD)

darwin-amd64:
	@mkdir -p $(DIST)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w $(LDFLAGS)" -o $(DIST)/$(BINARY)-darwin-amd64 $(CMD)

smoke:
	sh ./scripts/smoke.sh

clean:
	rm -rf $(DIST) $(BINARY)

.PHONY: all fmt lint test build linux-amd64 linux-arm64 linux-armv7 darwin-arm64 darwin-amd64 smoke clean
