#!/bin/bash

# Smart Todo CLI Installer
# This script installs the latest version of Smart Todo CLI

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default values
INSTALL_DIR="/usr/local/bin"
REPO="AhmedYacineAbdelmalek/todo"
VERSION="latest"

print_usage() {
    echo "Smart Todo CLI Installer"
    echo ""
    echo "Usage: $0 [options]"
    echo ""
    echo "Options:"
    echo "  -d, --dir DIR     Installation directory (default: /usr/local/bin)"
    echo "  -v, --version VER Version to install (default: latest)"
    echo "  -h, --help        Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                          # Install latest version to /usr/local/bin"
    echo "  $0 -d ~/.local/bin          # Install to custom directory"
    echo "  $0 -v v1.0.0               # Install specific version"
}

detect_platform() {
    local os=$(uname -s | tr '[:upper:]' '[:lower:]')
    local arch=$(uname -m)
    
    case $os in
        linux*)
            OS="linux"
            ;;
        darwin*)
            OS="darwin"
            ;;
        msys*|mingw*|cygwin*)
            OS="windows"
            ;;
        *)
            echo -e "${RED}Error: Unsupported operating system: $os${NC}"
            exit 1
            ;;
    esac
    
    case $arch in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        *)
            echo -e "${RED}Error: Unsupported architecture: $arch${NC}"
            exit 1
            ;;
    esac
    
    BINARY_NAME="todo-${OS}-${ARCH}"
    if [ "$OS" = "windows" ]; then
        BINARY_NAME="${BINARY_NAME}.exe"
    fi
}

get_latest_version() {
    echo -e "${YELLOW}Fetching latest version...${NC}"
    VERSION=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    if [ -z "$VERSION" ]; then
        echo -e "${RED}Error: Could not fetch latest version${NC}"
        exit 1
    fi
    echo -e "${GREEN}Latest version: $VERSION${NC}"
}

download_binary() {
    local url="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY_NAME}"
    local temp_file="/tmp/${BINARY_NAME}"
    
    echo -e "${YELLOW}Downloading from: $url${NC}"
    
    if command -v curl >/dev/null 2>&1; then
        curl -L -o "$temp_file" "$url"
    elif command -v wget >/dev/null 2>&1; then
        wget -O "$temp_file" "$url"
    else
        echo -e "${RED}Error: curl or wget is required${NC}"
        exit 1
    fi
    
    if [ ! -f "$temp_file" ]; then
        echo -e "${RED}Error: Download failed${NC}"
        exit 1
    fi
    
    echo "$temp_file"
}

install_binary() {
    local temp_file="$1"
    local target_file="$INSTALL_DIR/todo"
    
    # Create directory if it doesn't exist
    mkdir -p "$INSTALL_DIR"
    
    # Move binary to installation directory
    if [ -w "$INSTALL_DIR" ]; then
        mv "$temp_file" "$target_file"
    else
        echo -e "${YELLOW}Requesting sudo privileges to install to $INSTALL_DIR${NC}"
        sudo mv "$temp_file" "$target_file"
    fi
    
    # Make executable
    if [ -w "$target_file" ]; then
        chmod +x "$target_file"
    else
        sudo chmod +x "$target_file"
    fi
    
    echo -e "${GREEN}✅ Smart Todo CLI installed successfully to $target_file${NC}"
}

verify_installation() {
    if command -v todo >/dev/null 2>&1; then
        echo -e "${GREEN}✅ Installation verified!${NC}"
        echo ""
        todo version
        echo ""
        echo "🎉 You can now use 'todo' command!"
        echo "📖 Get started: todo add \"My first task\""
        echo "📚 Help: todo --help"
    else
        echo -e "${YELLOW}⚠️  Installation complete, but 'todo' not found in PATH${NC}"
        echo "Please add $INSTALL_DIR to your PATH or use the full path: $INSTALL_DIR/todo"
    fi
}

main() {
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -d|--dir)
                INSTALL_DIR="$2"
                shift 2
                ;;
            -v|--version)
                VERSION="$2"
                shift 2
                ;;
            -h|--help)
                print_usage
                exit 0
                ;;
            *)
                echo -e "${RED}Error: Unknown option $1${NC}"
                print_usage
                exit 1
                ;;
        esac
    done
    
    echo "🚀 Smart Todo CLI Installer"
    echo "=========================="
    
    detect_platform
    echo -e "${GREEN}Detected platform: $OS $ARCH${NC}"
    
    if [ "$VERSION" = "latest" ]; then
        get_latest_version
    fi
    
    temp_file=$(download_binary)
    install_binary "$temp_file"
    verify_installation
    
    echo ""
    echo "🎯 Happy productivity!"
}

main "$@"
