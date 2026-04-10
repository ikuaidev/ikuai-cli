#!/bin/sh
# ikuai-cli installer — downloads the latest (or pinned) release binary.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/ikuaidev/ikuai-cli/main/scripts/install.sh | sh
#   curl -fsSL ... | sh -s -- --version v0.1.0
#   curl -fsSL ... | sh -s -- --dir ~/.local/bin

set -eu

REPO="ikuaidev/ikuai-cli"
BINARY="ikuai-cli"
INSTALL_DIR="${HOME}/.local/bin"
VERSION=""

# ── Parse flags ──────────────────────────────────────────────────────
while [ $# -gt 0 ]; do
  case "$1" in
    --version)  VERSION="$2"; shift 2 ;;
    --dir)      INSTALL_DIR="$2"; shift 2 ;;
    -h|--help)
      echo "Usage: install.sh [--version vX.Y.Z] [--dir /path]"
      exit 0 ;;
    *) echo "Unknown flag: $1"; exit 1 ;;
  esac
done

# ── Detect OS / Arch ────────────────────────────────────────────────
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$OS" in
  linux)  ;;
  darwin) ;;
  *)      echo "Error: unsupported OS: $OS"; exit 1 ;;
esac

case "$ARCH" in
  x86_64|amd64)   ARCH="amd64" ;;
  aarch64|arm64)   ARCH="arm64" ;;
  *)               echo "Error: unsupported architecture: $ARCH"; exit 1 ;;
esac

# ── Resolve version ─────────────────────────────────────────────────
if [ -z "$VERSION" ]; then
  VERSION="$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
    | grep '"tag_name"' | head -1 | cut -d'"' -f4)"
  if [ -z "$VERSION" ]; then
    echo "Error: failed to fetch latest release version"
    exit 1
  fi
fi

# ── Download and install ────────────────────────────────────────────
ARCHIVE="${BINARY}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${VERSION}/${ARCHIVE}"
CHECKSUM_URL="https://github.com/${REPO}/releases/download/${VERSION}/checksums.txt"

TMPDIR="$(mktemp -d)"
trap 'rm -rf "$TMPDIR"' EXIT

echo "Downloading ${BINARY} ${VERSION} (${OS}/${ARCH})..."
curl -fsSL "$URL" -o "${TMPDIR}/${ARCHIVE}"
curl -fsSL "$CHECKSUM_URL" -o "${TMPDIR}/checksums.txt"

# Verify checksum
cd "$TMPDIR"
if command -v sha256sum >/dev/null 2>&1; then
  sha256sum --check --ignore-missing checksums.txt >/dev/null 2>&1
elif command -v shasum >/dev/null 2>&1; then
  shasum -a 256 --check --ignore-missing checksums.txt >/dev/null 2>&1
else
  echo "Warning: no sha256sum or shasum found, skipping checksum verification"
fi

# Extract
tar -xzf "${ARCHIVE}"

# Install
mkdir -p "$INSTALL_DIR"
mv "${BINARY}" "${INSTALL_DIR}/${BINARY}"
chmod +x "${INSTALL_DIR}/${BINARY}"

echo ""
echo "✓ ${BINARY} ${VERSION} installed to ${INSTALL_DIR}/${BINARY}"
echo ""

# Remind user to add to PATH if needed
case ":${PATH}:" in
  *":${INSTALL_DIR}:"*) ;;
  *) echo "Add to PATH: export PATH=\"${INSTALL_DIR}:\$PATH\""
     echo "" ;;
esac

echo "Get started:"
echo "  ${BINARY} auth set-url https://192.168.1.1"
echo "  ${BINARY} auth set-token <your-token>"
echo "  ${BINARY} auth status"
