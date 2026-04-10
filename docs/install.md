# Install

## From Release Binaries (Recommended)

Download a prebuilt binary from the [Releases page](https://github.com/ikuaidev/ikuai-cli/releases).

### Linux / macOS

```bash
# Replace VERSION and ARCH as needed (amd64 or arm64)
VERSION=0.1.0
ARCH=linux_amd64

curl -fsSL "https://github.com/ikuaidev/ikuai-cli/releases/download/v${VERSION}/ikuai-cli_${ARCH}.tar.gz" \
  -o ikuai-cli.tar.gz

# Verify checksum
curl -fsSL "https://github.com/ikuaidev/ikuai-cli/releases/download/v${VERSION}/checksums.txt" \
  -o checksums.txt
# Linux:
sha256sum --check --ignore-missing checksums.txt
# macOS:
# shasum -a 256 --check --ignore-missing checksums.txt

# Extract and install
tar -xzf ikuai-cli.tar.gz
sudo mv ikuai-cli /usr/local/bin/
ikuai-cli version
```

### Windows

Download `ikuai-cli_windows_amd64.zip` from the [Releases page](https://github.com/ikuaidev/ikuai-cli/releases),
extract it, and add `ikuai-cli.exe` to your `PATH`.

## From Source (Go)

Requires Go 1.21 or later:

```bash
go install github.com/ikuaidev/ikuai-cli/cmd/ikuai-cli@latest
```

## Build Locally

```bash
git clone https://github.com/ikuaidev/ikuai-cli.git
cd ikuai-cli
make build
./ikuai-cli version
```

## Shell Completion

Generate the completion script for your shell:

```bash
ikuai-cli completion bash
ikuai-cli completion zsh
ikuai-cli completion fish
ikuai-cli completion powershell
```

### Bash

```bash
mkdir -p ~/.local/share/bash-completion/completions
ikuai-cli completion bash > ~/.local/share/bash-completion/completions/ikuai-cli
```

### Zsh

```bash
mkdir -p ~/.zsh/completions
ikuai-cli completion zsh > ~/.zsh/completions/_ikuai-cli
# Ensure ~/.zsh/completions is in your $fpath
```

### Fish

```bash
mkdir -p ~/.config/fish/completions
ikuai-cli completion fish > ~/.config/fish/completions/ikuai-cli.fish
```

### PowerShell

```powershell
ikuai-cli completion powershell > ikuai-cli.ps1
. .\ikuai-cli.ps1
```
