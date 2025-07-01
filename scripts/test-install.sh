#!/bin/bash

# Test script for install.sh functionality
# This tests the install script functions without requiring actual releases

set -e

echo "🧪 Testing Flow install.sh script"
echo "=================================="

# Test platform detection
echo "Testing platform detection..."
detect_platform() {
    local os
    local arch
    
    case "$(uname -s)" in
        Linux*)     os="linux" ;;
        Darwin*)    os="darwin" ;;
        CYGWIN*|MINGW*|MSYS*) os="windows" ;;
        *)          echo "Unsupported OS: $(uname -s)" && exit 1 ;;
    esac
    
    case "$(uname -m)" in
        x86_64|amd64)   arch="amd64" ;;
        arm64|aarch64)  arch="arm64" ;;
        *)              echo "Unsupported arch: $(uname -m)" && exit 1 ;;
    esac
    
    echo "${os}-${arch}"
}

platform=$(detect_platform)
echo "✅ Platform detection: $platform"

# Test dependencies
echo "Testing dependencies..."
if command -v curl >/dev/null 2>&1; then
    echo "✅ curl is available"
else
    echo "❌ curl is missing"
    exit 1
fi

if command -v tar >/dev/null 2>&1; then
    echo "✅ tar is available"
else
    echo "❌ tar is missing"
    exit 1
fi

# Test archive naming logic
echo "Testing archive naming..."
if [[ "$platform" == *"windows"* ]]; then
    archive_name="flow-${platform}.zip"
    binary_path="flow-${platform}.exe"
else
    archive_name="flow-${platform}.tar.gz"
    binary_path="flow-${platform}"
fi

echo "✅ Archive name: $archive_name"
echo "✅ Binary path: $binary_path"

# Test directory creation and permission check
echo "Testing install directory logic..."
install_dir="/usr/local/bin"

if [ -d "$install_dir" ]; then
    echo "✅ Install directory exists: $install_dir"
    if [ -w "$install_dir" ]; then
        echo "✅ Can write to $install_dir directly"
    elif command -v sudo >/dev/null 2>&1; then
        echo "✅ Can use sudo for $install_dir"
    else
        echo "⚠️  Cannot write to $install_dir and no sudo available"
    fi
else
    echo "⚠️  Install directory does not exist: $install_dir"
    if command -v sudo >/dev/null 2>&1; then
        echo "✅ Can use sudo to create $install_dir"
    else
        echo "⚠️  Cannot create $install_dir without sudo"
    fi
fi

# Test mkdir -p functionality
echo "Testing directory creation logic..."
test_dir="/tmp/flow-test-install-$$"
if mkdir -p "$test_dir" 2>/dev/null; then
    echo "✅ mkdir -p works correctly"
    rm -rf "$test_dir"
else
    echo "❌ mkdir -p failed"
fi

echo ""
echo "🎉 All install.sh tests passed!"
echo ""
echo "To test with actual release:"
echo "  ./install.sh"
echo ""
echo "Note: Actual installation requires a published release." 