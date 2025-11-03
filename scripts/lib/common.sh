#!/bin/bash
# SCANOSS Code Compare - Common Installation Functions
# This library provides shared functionality for installation scripts

set -e  # Exit on error

# Colors for output
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
readonly NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${GREEN}==>${NC} $1" >&2
}

log_warn() {
    echo -e "${YELLOW}Warning:${NC} $1" >&2
}

log_error() {
    echo -e "${RED}Error:${NC} $1" >&2
}

log_debug() {
    if [[ "${DEBUG:-0}" == "1" ]]; then
        echo -e "${BLUE}Debug:${NC} $1" >&2
    fi
}

# Abort with error message
abort() {
    log_error "$1"
    exit 1
}

# Detect CPU architecture
detect_architecture() {
    local arch="$(uname -m)"

    case "$arch" in
        x86_64|amd64)
            echo "amd64"
            ;;
        aarch64|arm64)
            echo "arm64"
            ;;
        *)
            log_warn "Unknown architecture: $arch, defaulting to amd64"
            echo "amd64"
            ;;
    esac
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

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Download file with progress bar
download_file() {
    local url="$1"
    local output="$2"
    local retries="${3:-3}"
    local attempt=1

    log_info "Downloading: $url"

    while [ $attempt -le $retries ]; do
        if command_exists curl; then
            if curl -fSL --progress-bar "$url" -o "$output"; then
                return 0
            fi
        elif command_exists wget; then
            if wget -q --show-progress "$url" -O "$output"; then
                return 0
            fi
        else
            abort "Neither curl nor wget is available. Please install one of them."
        fi

        log_warn "Download attempt $attempt/$retries failed"
        attempt=$((attempt + 1))

        if [ $attempt -le $retries ]; then
            log_info "Retrying in 2 seconds..."
            sleep 2
        fi
    done

    abort "Failed to download after $retries attempts"
}

# Verify SHA256 checksum
verify_checksum() {
    local file="$1"
    local expected_checksum="$2"

    log_info "Verifying checksum..."

    local actual_checksum=""
    if command_exists shasum; then
        actual_checksum=$(shasum -a 256 "$file" | awk '{print $1}')
    elif command_exists sha256sum; then
        actual_checksum=$(sha256sum "$file" | awk '{print $1}')
    else
        log_warn "No checksum tool available (shasum or sha256sum)"
        log_warn "Skipping checksum verification"
        return 0
    fi

    # Convert to lowercase for comparison
    actual_checksum=$(echo "$actual_checksum" | tr '[:upper:]' '[:lower:]')
    expected_checksum=$(echo "$expected_checksum" | tr '[:upper:]' '[:lower:]')

    if [ "$actual_checksum" == "$expected_checksum" ]; then
        log_info "Checksum verified: ✓"
        return 0
    else
        log_error "Checksum mismatch!"
        log_error "Expected: $expected_checksum"
        log_error "Got:      $actual_checksum"
        return 1
    fi
}

# Create temporary directory and set trap for cleanup
# NOTE: The caller should set up the cleanup trap after calling this function
# Example: temp_dir=$(setup_temp_dir) && trap "cleanup_temp_dir '$temp_dir'" EXIT INT TERM
setup_temp_dir() {
    # macOS and Linux have different mktemp behavior
    # Use a portable approach that works on both
    mktemp -d "${TMPDIR:-/tmp}/scanoss-install.XXXXXXXXXX"
}

# Cleanup temporary directory
cleanup_temp_dir() {
    local temp_dir="$1"

    if [ -d "$temp_dir" ]; then
        log_debug "Cleaning up temporary directory: $temp_dir"
        rm -rf "$temp_dir"
    fi
}

# Confirm with user before proceeding
confirm() {
    local prompt="$1"
    local default="${2:-n}"

    if [ "$default" == "y" ]; then
        prompt="$prompt [Y/n]: "
    else
        prompt="$prompt [y/N]: "
    fi

    read -p "$prompt" response
    response=${response:-$default}

    case "$response" in
        [yY][eE][sS]|[yY])
            return 0
            ;;
        *)
            return 1
            ;;
    esac
}

# Check for required permissions
check_permissions() {
    local test_dir="$1"

    if [ ! -w "$test_dir" ]; then
        log_error "No write permission for: $test_dir"
        log_error "You may need to run this script with elevated permissions (sudo)"
        return 1
    fi

    return 0
}

# Get latest GitHub release version
get_latest_version() {
    local repo="$1"

    log_info "Fetching latest version..."

    local version=""
    if command_exists curl; then
        version=$(curl -fsSL "https://api.github.com/repos/$repo/releases/latest" | grep '"tag_name"' | sed -E 's/.*"v?([^"]+)".*/\1/')
    elif command_exists wget; then
        version=$(wget -qO- "https://api.github.com/repos/$repo/releases/latest" | grep '"tag_name"' | sed -E 's/.*"v?([^"]+)".*/\1/')
    else
        abort "Neither curl nor wget is available"
    fi

    if [ -z "$version" ]; then
        abort "Failed to fetch latest version"
    fi

    log_info "Latest version: v$version"
    echo "$version"
}

# Show completion message
show_completion() {
    echo >&2
    log_info "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    log_info "✨ SCANOSS Code Compare has been successfully installed!"
    log_info "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo >&2
    log_info "Verify installation:"
    echo "    scanoss-cc --version" >&2
    echo >&2
    log_info "Get started:"
    echo "    scanoss-cc --help" >&2
    echo >&2
    log_info "Documentation:"
    echo "    https://github.com/scanoss/scanoss.cc" >&2
    echo >&2
}

# Export functions for use in other scripts
export -f log_info log_warn log_error log_debug abort
export -f detect_architecture detect_os command_exists
export -f download_file verify_checksum
export -f setup_temp_dir cleanup_temp_dir
export -f confirm check_permissions
export -f get_latest_version show_completion
