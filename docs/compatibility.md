# Compatibility

This document records the current compatibility contract for `ikuai-cli`.

It separates what is verified in this repository from what still needs real-device confirmation.

## API Contract

`ikuai-cli` currently targets the local iKuai HTTP API under:

```text
/api/v4.0
```

Current contract assumptions:

- the router exposes the local v4.0 API
- authentication is Bearer token-based (`Authorization: Bearer <token>`)
- many resource collections use shared pagination fields:
  - `page`
  - `page_size`
  - `filter`
  - `order`
  - `order_by`

## Runtime Compatibility Handling

### Self-signed TLS certificates

Routers commonly use self-signed certificates on the local management interface.

The built-in HTTP client intentionally skips TLS verification for local API access.

### Bare `nil` values in firmware responses

Some firmware versions can emit invalid JSON containing bare `nil`.

The client normalizes bare `nil` to `null` before JSON parsing.

Reference implementation:

- `internal/api/client.go`

## Verification Matrix

### Repository-verified matrix

| Resource Group | Status | Notes |
| --- | --- | --- |
| `auth` | covered | set-url/set-token/clear/session persistence paths covered |
| `monitor` | covered | read-only monitoring endpoint behavior covered |
| `network` | covered | collection query params and create/update style requests covered |
| `system` | covered | collection query params and config update request body covered |
| `users` | covered | account listing and creation covered |
| `objects` | covered | object listing and creation covered |
| `routing` | covered | static route listing and stream rule creation covered |
| `qos` | covered | QoS listing and creation covered |
| `security` | covered | ACL listing and MAC mode update covered |
| `vpn` | covered | OpenVPN client listing and WireGuard peer creation covered |
| `wireless` | covered | blacklist listing and AC service update covered |
| `advanced` | covered | advanced-service user listing and SNMPD update covered |
| `auth-server` | covered | auth web service get covered |
| `log` | covered | log listing and clear behavior covered |

These tests validate request method, path, query parameters, JSON body shape, and output behavior at the CLI layer.

### Environment and firmware matrix

| Dimension | Evidence Type | Current Status | Notes |
| --- | --- | --- | --- |
| API namespace | repository tests + code inspection | assumed and covered | current code targets `/api/v4.0` only |
| Authentication | repository tests + code inspection | assumed and covered | Bearer token via `Authorization` header |
| TLS | code inspection | handled | self-signed local certs are accepted by the built-in client |
| Response JSON normalization | repository tests + code inspection | handled | bare `nil` is normalized before JSON parsing |
| Firmware versions | manual device validation required | not yet verified | no real-device version matrix is declared yet |
| Localization variants | manual device validation required | not yet verified | no language-specific firmware matrix is declared yet |
| Undocumented fields | manual device validation required | not yet verified | CLI only guarantees the request/response shapes covered in tests |

## Validation Workflow

To promote a compatibility assumption into a declared support statement:

1. Record the router firmware version and locale.
2. Run the relevant CLI commands against a real device.
3. Compare request shape and response shape against repository behavior tests.
4. Add or update a behavior test when a new endpoint shape or field contract appears.
5. Update this matrix only after the device behavior is confirmed.

## Not Yet Declared Stable

The project does not yet claim a fully verified support matrix across:

- all iKuai firmware releases
- all localized firmware variants
- every documented or undocumented API field

If you hit an incompatibility, include the following in the issue:

- `ikuai-cli version --format json`
- router firmware version
- the exact command used
- sanitized response body if available

## Recommended Upgrade Policy

When adding new commands or changing request shapes:

1. Prefer patterns already used by nearby resource groups.
2. Add at least one behavior test for the new endpoint shape.
3. Add a device-validation note here if a firmware-specific assumption is discovered.
