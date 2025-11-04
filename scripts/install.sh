#!/bin/bash
# SCANOSS Code Compare - Universal Installer
# This script detects your platform and installs SCANOSS Code Compare

set -e

# Repository
readonly REPO="scanoss/scanoss.cc"
readonly RAW_BASE_URL="https://raw.githubusercontent.com/$REPO/main/scripts"

# Colors for output
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
readonly NC='\033[0m'

# Simple logging (don't rely on common.sh yet)
log_info() {
    echo -e "${GREEN}==>${NC} $1" >&2
}

log_warn() {
    echo -e "${YELLOW}Warning:${NC} $1" >&2
}

log_error() {
    echo -e "${RED}Error:${NC} $1" >&2
}

abort() {
    log_error "$1"
    exit 1
}

# Detect operating system
detect_os() {
    local os="$(uname -s)"

    case "$os" in
        Darwin*)
            echo "macos"
            ;;
        Linux*)
            echo "linux"
            ;;
        MINGW*|MSYS*|CYGWIN*)
            echo "windows"
            ;;
        *)
            abort "Unsupported operating system: $os"
            ;;
    esac
}

# Detect if running in Git Bash on Windows
is_git_bash_windows() {
    if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" ]]; then
        return 0
    fi
    return 1
}

# Show banner
show_banner() {
    echo "" >&2
    log_info "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    log_info "  SCANOSS Code Compare - Universal Installer"
    log_info "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "" >&2
}

# Download and execute platform-specific installer
install_for_platform() {
    local platform="$1"

    case "$platform" in
        macos)
            log_info "Detected: macOS"
            log_info "Downloading macOS installer..."
            echo "" >&2

            # Download and execute macOS installer
            if command -v curl >/dev/null 2>&1; then
                bash <(curl -fsSL "$RAW_BASE_URL/install-macos.sh") "$@"
            elif command -v wget >/dev/null 2>&1; then
                bash <(wget -qO- "$RAW_BASE_URL/install-macos.sh") "$@"
            else
                abort "Neither curl nor wget is available. Please install one of them."
            fi
            ;;

        linux)
            log_info "Detected: Linux"
            log_info "Downloading Linux installer..."
            echo "" >&2

            # Download and execute Linux installer
            if command -v curl >/dev/null 2>&1; then
                bash <(curl -fsSL "$RAW_BASE_URL/install-linux.sh") "$@"
            elif command -v wget >/dev/null 2>&1; then
                bash <(wget -qO- "$RAW_BASE_URL/install-linux.sh") "$@"
            else
                abort "Neither curl nor wget is available. Please install one of them."
            fi
            ;;

        windows)
            log_error "Windows detected"
            echo "" >&2
            log_info "Please use the PowerShell installer for Windows:"
            echo "" >&2
            echo "  1. Open PowerShell as Administrator" >&2
            echo "  2. Run:" >&2
            echo "" >&2
            echo "     irm https://raw.githubusercontent.com/$REPO/main/scripts/install-windows.ps1 | iex" >&2
            echo "" >&2
            log_warn "If you're in Git Bash, please switch to PowerShell"
            echo "" >&2
            exit 1
            ;;

        *)
            abort "Unsupported platform: $platform"
            ;;
    esac
}

# Check prerequisites
check_prerequisites() {
    # Check for curl or wget
    if ! command -v curl >/dev/null 2>&1 && ! command -v wget >/dev/null 2>&1; then
        abort "Neither curl nor wget is available. Please install one of them and try again."
    fi
}

# Main function
main() {
    show_banner

    # Check prerequisites
    check_prerequisites

    # Detect platform
    local platform=$(detect_os)

    # Special handling for Git Bash on Windows
    if is_git_bash_windows; then
        log_warn "Detected Git Bash on Windows"
        platform="windows"
    fi

    # Install for platform
    install_for_platform "$platform" "$@"
}

main "$@"
