# Examples

This directory contains sanitized example outputs used by the project docs.

These files are intentionally static:

- they show output shape, not a live router snapshot
- they are safe to copy into docs and issues
- they should stay aligned with CLI output contracts

## First-Release Baseline

The first release keeps a minimal fixture set that covers the main CLI usage patterns:

- auth/session
- read-only inspection
- list-style queries
- write/update responses
- build metadata

Current files:

- `auth-status.pretty.json`
- `auth-status.compact.json`
- `monitor-system.pretty.json`
- `system-set-response.compact.json`
- `users-accounts-list.compact.json`
- `version.compact.json`

## Example Paths

### Auth and session

Command:

```bash
ikuai-cli auth status
ikuai-cli auth status --format json
```

Fixtures:

- `auth-status.pretty.json`
- `auth-status.compact.json`

### Read-only inspection

Command:

```bash
ikuai-cli monitor system
```

Fixture:

- `monitor-system.pretty.json`

### List-style query

Command:

```bash
ikuai-cli users accounts list --format json
```

Fixture:

- `users-accounts-list.compact.json`

### Write/update response

Command:

```bash
ikuai-cli system set --hostname "ikuai-gw" --format json
```

Fixture:

- `system-set-response.compact.json`

### Build metadata

Command:

```bash
ikuai-cli version --format json
```

Fixture:

- `version.compact.json`

## Maintenance Rule

When a documented output contract changes, update the corresponding fixture in this directory in the same change.
