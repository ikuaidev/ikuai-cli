# Output Modes

`ikuai-cli` supports four output modes.

Example fixtures in `docs/examples/` are sanitized static samples used by the docs.

## Table (default)

Default output is a human-readable table for terminal readability.

Example:

```bash
ikuai-cli auth status
```

## JSON

Pass `--format json` to emit compact single-line JSON for scripts and agents.

Example:

```bash
ikuai-cli auth status --format json
```

Representative fixtures:

- `docs/examples/auth-status.compact.json`
- `docs/examples/system-set-response.compact.json`
- `docs/examples/users-accounts-list.compact.json`
- `docs/examples/version.compact.json`

## YAML

Pass `--format yaml` to emit YAML output, useful for configuration files and token-efficient agent consumption.

Example:

```bash
ikuai-cli auth status --format yaml
```

## Raw

Pass `--raw` to emit the full API envelope including metadata, pagination info, and status codes. Useful for debugging.

Example:

```bash
ikuai-cli auth status --raw
```

## Response Shapes

Different command types produce different JSON shapes:

- **Read/list commands:** return the data payload directly (e.g., `{"items":[...]}` or `{"sysinfo":{...}}`)
- **Write/update commands:** return `{"message":"saved"}` or similar confirmation
- **Create commands:** return `{"message":"success","rowid":42}` — the `rowid` field contains the new resource ID

## Rules

- Successful command output goes to stdout.
- Errors go to stderr.
- `--format` changes formatting, not payload semantics.
- `--raw` includes the full API response envelope.
- New commands should follow the same output expectations as existing commands.
- Example fixtures in `docs/examples/` should be updated when output contracts change.
