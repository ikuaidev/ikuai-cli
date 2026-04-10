# Releasing

This repository uses GoReleaser for tagged releases.

## Pre-release Checklist

Run the standard checks locally:

```bash
make fmt
make lint
make test
make build
```

Review before tagging:

- `README.md`
- `.goreleaser.yml`

## Build Metadata

The CLI embeds build metadata through `ldflags`:

- `version`
- `commit`
- `date`

Local example:

```bash
make build VERSION=0.1.0
./ikuai-cli version --format json
```

## Cut a Release

1. Commit any release prep changes.
2. Create and push a tag:

```bash
git tag v0.1.0
git push origin v0.1.0
```

4. Run GoReleaser locally if needed:

```bash
goreleaser release --clean
```

Or dry-run first:

```bash
goreleaser release --snapshot --clean
```

## Expected Artifacts

GoReleaser is configured to publish:

- `tar.gz` archives
- `checksums.txt`

Supported targets:

- Linux `amd64`
- Linux `arm64`
- Darwin `amd64`
- Darwin `arm64`
