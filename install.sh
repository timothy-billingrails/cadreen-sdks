#!/bin/sh
set -e

# Cadreen CLI installer
# Usage: curl -fsSL https://raw.githubusercontent.com/timothy-billingrails/cadreen-sdks/master/install.sh | sh

REPO="timothy-billingrails/cadreen-sdks"
BINARY="cadreen"

# Detect OS
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$OS" in
    linux)  OS="linux" ;;
    darwin) OS="darwin" ;;
    mingw*|msys*|cygwin*) OS="windows" ;;
    *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

# Detect architecture
ARCH="$(uname -m)"
case "$ARCH" in
    x86_64|amd64)  ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Determine version
VERSION="${CADREEN_VERSION:-latest}"
if [ "$VERSION" = "latest" ]; then
    echo "Fetching latest version..."
    VERSION=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed 's/.*"cli-v\(.*\)".*/\1/')
    if [ -z "$VERSION" ]; then
        echo "Failed to fetch latest version."
        exit 1
    fi
fi

echo "Installing Cadreen CLI ${VERSION} for ${OS}/${ARCH}..."

# Determine install directory
INSTALL_DIR="${CADREEN_INSTALL_DIR:-/usr/local/bin}"
if [ ! -d "$INSTALL_DIR" ]; then
    echo "Creating install directory: $INSTALL_DIR"
    mkdir -p "$INSTALL_DIR"
fi

# Download URL
EXT=""
if [ "$OS" = "windows" ]; then
    EXT=".exe"
fi
FILENAME="cadreen_${OS}_${ARCH}${EXT}"
URL="https://github.com/${REPO}/releases/download/cli-v${VERSION}/${FILENAME}"

# Download
TMPDIR="$(mktemp -d)"
trap 'rm -rf "$TMPDIR"' EXIT

echo "Downloading from: $URL"
curl -fsSL "$URL" -o "${TMPDIR}/${FILENAME}"

# Verify download
if [ ! -s "${TMPDIR}/${FILENAME}" ]; then
    echo "Download failed or file is empty."
    exit 1
fi

# Install
INSTALL_PATH="${INSTALL_DIR}/${BINARY}${EXT}"
mv "${TMPDIR}/${FILENAME}" "$INSTALL_PATH"
chmod +x "$INSTALL_PATH"

echo ""
echo "Cadreen CLI ${VERSION} installed to: ${INSTALL_PATH}"
echo ""
echo "Next steps:"
echo "  cadreen init    — Set up your account"
echo "  cadreen --help  — See all commands"
echo ""
echo "To add shell completions:"
echo "  cadreen completion bash >> ~/.bashrc"
echo "  cadreen completion zsh >> ~/.zshrc"
