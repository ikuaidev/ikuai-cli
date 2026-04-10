#!/usr/bin/env sh
set -eu

go test ./...
go build ./cmd/ikuai-cli
go run ./cmd/ikuai-cli auth status --format json >/dev/null
