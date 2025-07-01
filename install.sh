#!/bin/bash

set -e

# Flow installer script
# Usage: curl -sSL https://raw.githubusercontent.com/e6a5/flow/main/install.sh | bash

REPO="e6a5/flow"
BINARY_NAME="flow"
INSTALL_DIR="/usr/local/bin"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions
log() {
    echo -e "${BLUE}[Flow]${NC} $1"
}

success() {
    echo -e "${GREEN}[Flow]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[Flow]${NC} $1"
}

error() {
    echo -e "${RED}[Flow]${NC} $1" >&2
    exit 1
}

# Detect OS and architecture
detect_platform() {
    local os
    local arch
    
    case "$(uname -s)" in
        Linux*)     os="linux" ;;
        Darwin*)    os="darwin" ;;
        CYGWIN*|MINGW*|MSYS*) os="windows" ;;
        *)          error "Unsupported operating system: $(uname -s)" ;;
    esac
    
    case "$(uname -m)" in
        x86_64|amd64)   arch="amd64" ;;
        arm64|aarch64)  arch="arm64" ;;
        *)              error "Unsupported architecture: $(uname -m)" ;;
    esac
    
    echo "${os}-${arch}"
}

# Get latest release version
get_latest_version() {
    local version
    version=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    if [ -z "$version" ]; then
        error "Failed to get latest version"
    fi
    echo "$version"
}

# Download and install
install_flow() {
    local platform
    local version
    local download_url
    local archive_name
    local binary_path
    
    platform=$(detect_platform)
    version=$(get_latest_version)
    
    log "Detected platform: $platform"
    log "Latest version: $version"
    
    # Determine archive format
    if [[ "$platform" == *"windows"* ]]; then
        archive_name="${BINARY_NAME}-${platform}.zip"
        binary_path="${BINARY_NAME}-${platform}.exe"
    else
        archive_name="${BINARY_NAME}-${platform}.tar.gz"
        binary_path="${BINARY_NAME}-${platform}"
    fi
    
    download_url="https://github.com/${REPO}/releases/download/${version}/${archive_name}"
    
    log "Downloading from: $download_url"
    
    # Create temporary directory
    local tmp_dir
    tmp_dir=$(mktemp -d)
    trap "rm -rf $tmp_dir" EXIT
    
    # Download
    if ! curl -sL "$download_url" -o "$tmp_dir/$archive_name"; then
        error "Failed to download $archive_name"
    fi
    
    # Extract
    cd "$tmp_dir"
    if [[ "$archive_name" == *.zip ]]; then
        if ! command -v unzip >/dev/null 2>&1; then
            error "unzip is required but not installed"
        fi
        unzip -q "$archive_name"
    else
        tar -xzf "$archive_name"
    fi
    
    # Check if binary exists
    if [ ! -f "$binary_path" ]; then
        error "Binary not found in archive: $binary_path"
    fi
    
    # Install binary
    log "Installing to $INSTALL_DIR/$BINARY_NAME"
    
    # Create install directory if it doesn't exist
    if [ ! -d "$INSTALL_DIR" ]; then
        log "Creating install directory: $INSTALL_DIR"
        if command -v sudo >/dev/null 2>&1; then
            sudo mkdir -p "$INSTALL_DIR"
        else
            mkdir -p "$INSTALL_DIR" 2>/dev/null || error "Cannot create $INSTALL_DIR and sudo is not available"
        fi
    fi
    
    # Check if we need sudo for writing
    if [ ! -w "$INSTALL_DIR" ]; then
        if command -v sudo >/dev/null 2>&1; then
            sudo mv "$binary_path" "$INSTALL_DIR/$BINARY_NAME"
            sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
        else
            error "Cannot write to $INSTALL_DIR and sudo is not available"
        fi
    else
        mv "$binary_path" "$INSTALL_DIR/$BINARY_NAME"
        chmod +x "$INSTALL_DIR/$BINARY_NAME"
    fi
    
    success "Flow $version installed successfully!"
    
    # Verify installation
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        success "Verification: $(flow --version)"
        echo
        log "🌊 Flow is ready! Try:"
        log "  flow start --tag \"writing docs\""
        log "  flow status"
        log "  flow end"
        echo
        log "For help: flow --help"
        log "Documentation: https://github.com/$REPO"
    else
        warn "Installation completed, but '$BINARY_NAME' is not in PATH"
        warn "You may need to:"
        warn "  - Restart your terminal"
        warn "  - Add $INSTALL_DIR to your PATH"
        warn "  - Or run: export PATH=\"$INSTALL_DIR:\$PATH\""
    fi

    # Shell completion instructions
    SHELL_TYPE=$(basename "$SHELL")
    log "To enable shell completion, add the following to your shell's config file:"

    case "$SHELL_TYPE" in
        bash)
            log "\n# For Bash:"
            log 'echo ''eval "$(flow completion bash)"'' >> ~/.bashrc'
            warn "You may need to restart your shell for changes to take effect."
            ;;
        zsh)
            log "\n# For Zsh:"
            log 'echo ''eval "$(flow completion zsh)"'' >> ~/.zshrc'
            warn "You may need to restart your shell for changes to take effect."
            ;;
        *)
            warn "Unsupported shell for automatic completion setup: ${SHELL_TYPE}"
            warn "Run 'flow completion --help' for manual instructions."
            ;;
    esac
}

# Main
main() {
    echo "🌊 Flow Installer"
    echo "=================="
    
    # Check dependencies
    if ! command -v curl >/dev/null 2>&1; then
        error "curl is required but not installed"
    fi
    
    if ! command -v tar >/dev/null 2>&1; then
        error "tar is required but not installed"
    fi
    
    # Install
    install_flow
}

# Run main function
main "$@" 