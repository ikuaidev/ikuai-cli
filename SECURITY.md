# Security Policy

## Reporting a Vulnerability

**Do not open a public GitHub Issue for security vulnerabilities.**

Report via [GitHub Security Advisories](https://github.com/ikuaidev/ikuai-cli/security/advisories/new).

- Acknowledgement within **3 business days**
- Fix or mitigation plan within **30 days** for confirmed issues

## Scope

Security-sensitive areas:

- Config storage and credential handling (`~/.ikuai-cli/config.json`)
- API authentication flows
- SSH-based recovery flows

## Out of Scope

- Vulnerabilities in iKuai router firmware itself
- Issues requiring physical access to the router
