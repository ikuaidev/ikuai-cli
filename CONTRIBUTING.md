# Contributing

Thank you for your interest in contributing to ikuai-cli.

## Prerequisites

- Go 1.21+
- [golangci-lint](https://golangci-lint.run/usage/install/) v2.x
- Make

## Getting Started

```bash
git clone https://github.com/ikuaidev/ikuai-cli.git
cd ikuai-cli
make test && make build
./ikuai-cli version
```

## Before Submitting a PR

```bash
make fmt       # format code
make lint      # golangci-lint
make test      # all tests
make build     # verify binary compiles
```

## Guidelines

- Keep changes focused — one concern per PR
- Update tests when behavior changes
- Follow existing command shape: `<resource> <action> [args] [flags]`
- Default table output + `--format table|json|yaml` parity

## Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md).
