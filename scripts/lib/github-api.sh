#!/bin/bash
# SCANOSS Code Compare - GitHub API Helper Functions

# GitHub repository details
readonly GITHUB_REPO="scanoss/scanoss.cc"
readonly GITHUB_API_URL="https://api.github.com/repos/$GITHUB_REPO"

# Get release asset download URL for a specific platform
get_asset_url() {
    local version="$1"
    local platform="$2"  # macos, linux, windows
    local arch="${3:-amd64}"

    log_debug "Fetching asset URL for platform: $platform, arch: $arch, version: $version"

    # Determine the asset name pattern based on platform
    local asset_pattern=""
    case "$platform" in
        macos)
            asset_pattern="scanoss-cc-mac.zip"
            ;;
        linux)
            asset_pattern="scanoss-cc-linux.zip"
            ;;
        windows)
            asset_pattern="scanoss-cc-win.zip"
            ;;
        *)
            abort "Unknown platform: $platform"
            ;;
    esac

    local release_url="$GITHUB_API_URL/releases/tags/v$version"
    local url=""

    if command_exists curl; then
        url=$(curl -fsSL "$release_url" | \
            grep "browser_download_url.*$asset_pattern" | \
            head -n 1 | \
            cut -d '"' -f 4)
    elif command_exists wget; then
        url=$(wget -qO- "$release_url" | \
            grep "browser_download_url.*$asset_pattern" | \
            head -n 1 | \
            cut -d '"' -f 4)
    else
        abort "Neither curl nor wget is available"
    fi

    if [ -z "$url" ]; then
        abort "Failed to find download URL for $asset_pattern in release v$version"
    fi

    log_debug "Found asset URL: $url"
    echo "$url"
}

# Get checksum for a specific asset
get_asset_checksum() {
    local version="$1"
    local asset_name="$2"

    log_debug "Fetching checksum for: $asset_name"

    # Try to fetch SHA256SUMS file if it exists
    local checksums_url="https://github.com/$GITHUB_REPO/releases/download/v$version/SHA256SUMS"
    local checksum=""

    # Temporarily disable errexit to handle missing checksums gracefully
    set +e

    if command_exists curl; then
        checksum=$(curl -sSL "$checksums_url" 2>/dev/null | grep "$asset_name" | awk '{print $1}')
    elif command_exists wget; then
        checksum=$(wget -qO- "$checksums_url" 2>/dev/null | grep "$asset_name" | awk '{print $1}')
    fi

    # Restore errexit
    set -e

    if [ -z "$checksum" ]; then
        log_warn "No checksum file found for this release"
        log_warn "Installation will proceed without checksum verification"
        return 1
    fi

    log_debug "Found checksum: $checksum"
    echo "$checksum"
}

# Check if a specific version exists
version_exists() {
    local version="$1"
    local release_url="$GITHUB_API_URL/releases/tags/v$version"
    local http_code=""

    # Temporarily disable errexit to handle HTTP errors gracefully
    set +e

    if command_exists curl; then
        # Get HTTP status code without failing on 4xx/5xx
        http_code=$(curl -sSL -o /dev/null -w "%{http_code}" "$release_url" 2>/dev/null)
    elif command_exists wget; then
        # Get HTTP status code from wget
        http_code=$(wget --spider -S "$release_url" 2>&1 | grep "HTTP/" | tail -n 1 | awk '{print $2}')
    fi

    # Restore errexit
    set -e

    # Check if we got a 200 response
    if [ "$http_code" = "200" ]; then
        return 0
    else
        return 1
    fi
}

# List available versions
list_versions() {
    local limit="${1:-10}"

    log_info "Fetching available versions..."

    if command_exists curl; then
        curl -fsSL "$GITHUB_API_URL/releases?per_page=$limit" | \
            grep '"tag_name"' | \
            sed -E 's/.*"v?([^"]+)".*/\1/' | \
            head -n "$limit"
    elif command_exists wget; then
        wget -qO- "$GITHUB_API_URL/releases?per_page=$limit" | \
            grep '"tag_name"' | \
            sed -E 's/.*"v?([^"]+)".*/\1/' | \
            head -n "$limit"
    else
        abort "Neither curl nor wget is available"
    fi
}

# Get release notes for a version
get_release_notes() {
    local version="$1"
    local release_url="$GITHUB_API_URL/releases/tags/v$version"

    if command_exists curl; then
        curl -fsSL "$release_url" | \
            grep -A 1000 '"body":' | \
            head -n 50 | \
            sed 's/"body": "//g' | \
            sed 's/\\n/\n/g' | \
            sed 's/\\"/"/g'
    elif command_exists wget; then
        wget -qO- "$release_url" | \
            grep -A 1000 '"body":' | \
            head -n 50 | \
            sed 's/"body": "//g' | \
            sed 's/\\n/\n/g' | \
            sed 's/\\"/"/g'
    fi
}

# Download and verify a release asset
download_and_verify_asset() {
    local version="$1"
    local platform="$2"
    local output_path="$3"
    local arch="${4:-amd64}"

    # Get download URL
    local download_url=$(get_asset_url "$version" "$platform" "$arch")
    local asset_name=$(basename "$download_url")

    # Download the file
    download_file "$download_url" "$output_path"

    # Try to verify checksum
    local expected_checksum=""
    if expected_checksum=$(get_asset_checksum "$version" "$asset_name" 2>/dev/null); then
        if ! verify_checksum "$output_path" "$expected_checksum"; then
            abort "Checksum verification failed! The downloaded file may be corrupted or tampered with."
        fi
    fi

    log_info "Download completed: $output_path"
}

# Export functions
export -f get_asset_url get_asset_checksum
export -f version_exists list_versions get_release_notes
export -f download_and_verify_asset
