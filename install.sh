#!/bin/bash
# InstantTLS Installer
# Usage: curl -fsSL https://raw.githubusercontent.com/CyberWarBaby/Instant-TLS/main/install.sh | bash

set -e

REPO="github.com/CyberWarBaby/Instant-TLS/cli/cmd/instanttls"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="instanttls"

echo ""
echo "  ██╗███╗   ██╗███████╗████████╗ █████╗ ███╗   ██╗████████╗████████╗██╗     ███████╗"
echo "  ██║████╗  ██║██╔════╝╚══██╔══╝██╔══██╗████╗  ██║╚══██╔══╝╚══██╔══╝██║     ██╔════╝"
echo "  ██║██╔██╗ ██║███████╗   ██║   ███████║██╔██╗ ██║   ██║      ██║   ██║     ███████╗"
echo "  ██║██║╚██╗██║╚════██║   ██║   ██╔══██║██║╚██╗██║   ██║      ██║   ██║     ╚════██║"
echo "  ██║██║ ╚████║███████║   ██║   ██║  ██║██║ ╚████║   ██║      ██║   ███████╗███████║"
echo "  ╚═╝╚═╝  ╚═══╝╚══════╝   ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═══╝   ╚═╝      ╚═╝   ╚══════╝╚══════╝"
echo ""
echo "  Installing InstantTLS..."
echo ""

# Check for Go
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first:"
    echo "   https://golang.org/dl/"
    exit 1
fi

echo "📦 Installing via go install..."
go install ${REPO}@latest

# Find the installed binary
GO_BIN=$(go env GOPATH)/bin/${BINARY_NAME}

if [ ! -f "$GO_BIN" ]; then
    echo "❌ Installation failed. Binary not found at $GO_BIN"
    exit 1
fi

echo "🔗 Making instanttls available system-wide..."

# Copy to /usr/local/bin (requires sudo)
if [ -w "$INSTALL_DIR" ]; then
    cp "$GO_BIN" "$INSTALL_DIR/$BINARY_NAME"
else
    echo "   (requires sudo)"
    sudo cp "$GO_BIN" "$INSTALL_DIR/$BINARY_NAME"
    sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
fi

echo ""
echo "✅ InstantTLS installed successfully!"
echo ""
echo "   Now run:"
echo "   sudo instanttls setup"
echo ""
