#!/bin/bash
# SCANOSS Code Compare - Linux Installer
# This script installs SCANOSS Code Compare on Linux

set -e

# Repository and app details
readonly REPO="scanoss/scanoss.cc"
readonly APP_NAME="scanoss-cc"
readonly INSTALL_DIR="/usr/local/bin"
BINARY_NAME="$APP_NAME-linux-amd64"

# Source common functions
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
if [ -f "$SCRIPT_DIR/lib/common.sh" ]; then
    source "$SCRIPT_DIR/lib/common.sh"
    source "$SCRIPT_DIR/lib/github-api.sh"
else
    # If running from curl, download the library
    echo "Downloading installation libraries..." >&2
    TEMP_LIB_DIR=$(mktemp -d)
    curl -fsSL "https://raw.githubusercontent.com/$REPO/main/scripts/lib/common.sh" -o "$TEMP_LIB_DIR/common.sh"
    curl -fsSL "https://raw.githubusercontent.com/$REPO/main/scripts/lib/github-api.sh" -o "$TEMP_LIB_DIR/github-api.sh"
    source "$TEMP_LIB_DIR/common.sh"
    source "$TEMP_LIB_DIR/github-api.sh"
    trap 'rm -rf "$TEMP_LIB_DIR"' EXIT
fi

# Detect Linux distribution
detect_distro() {
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        echo "$ID"
    elif [ -f /etc/debian_version ]; then
        echo "debian"
    elif [ -f /etc/fedora-release ]; then
        echo "fedora"
    elif [ -f /etc/arch-release ]; then
        echo "arch"
    else
        echo "unknown"
    fi
}

detect_distro_version(){
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        echo "$VERSION_ID"
    else
        echo "unknown"
    fi
}

# Install dependencies based on distribution
install_dependencies() {
    local distro=$(detect_distro)
    local distro_version=$(detect_distro_version)

    log_info "Detected distribution: $distro $distro_version"
    log_info "Installing dependencies..."
    echo >&2

    # As default we use webkit version for ubuntu <=22.04 and debian <=13
    local webkit_package="libwebkit2gtk-4.0-dev"

    # If it's either ubuntu >=24 or debian >=13 we use webkit 4.1
    if [ "$distro_version" != "unknown" ] && { [ "$distro" = "ubuntu" ] && [ "${distro_version%%.*}" -ge 24 ] || [ "$distro" = "debian" ] && [ "${distro_version%%.*}" -ge 13 ]; }; then
        webkit_package="libwebkit2gtk-4.1-dev"
        BINARY_NAME="$APP_NAME-linux-amd64-webkit41"
    fi

    case "$distro" in
        ubuntu|debian|pop|linuxmint)
            log_info "Using apt package manager..."
            $SUDO apt-get update
            $SUDO apt-get install -y libgtk-3-0 $webkit_package || {
                log_warn "Some dependencies may not be available"
                log_warn "The application may still work"
            }
            ;;
        fedora|rhel|centos)
            log_info "Using dnf package manager..."
            $SUDO dnf install -y gtk3 webkit2gtk4.0 || {
                log_warn "Some dependencies may not be available"
                log_warn "The application may still work"
            }
            ;;
        arch|manjaro)
            log_info "Using pacman package manager..."
            $SUDO pacman -S --noconfirm gtk3 webkit2gtk || {
                log_warn "Some dependencies may not be available"
                log_warn "The application may still work"
            }
            ;;
        opensuse*|sles)
            log_info "Using zypper package manager..."
            $SUDO zypper install -y gtk3 webkit2gtk3 || {
                log_warn "Some dependencies may not be available"
                log_warn "The application may still work"
            }
            ;;
        *)
            log_warn "Unknown distribution: $distro"
            log_warn "You may need to manually install GTK3 and WebKit2GTK"
            log_warn "Required packages:"
            log_warn "  - GTK 3.0"
            log_warn "  - WebKit2GTK 4.0"
            log_warn "Continuing with installation..."
            echo >&2
            ;;
    esac

    echo >&2
    log_info "✓ Dependencies check complete"
}


# Main installation logic
install_linux() {
    log_info "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    log_info "  SCANOSS Code Compare - Linux Installation"
    log_info "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo >&2

    # Check if running on Linux
    if [ "$(uname -s)" != "Linux" ]; then
        abort "This installer is for Linux only. Please use the appropriate installer for your platform."
    fi

    setup_sudo
    echo >&2

    # Install dependencies automatically
    install_dependencies

    # Setup temporary directory
    local temp_dir=$(setup_temp_dir)
    trap 'cleanup_temp_dir "$temp_dir"' EXIT INT TERM

    # Get latest version
    local version=$(get_latest_version "$REPO")

    # Download the Linux binary
    local zip_path="$temp_dir/$BINARY_NAME.zip"

    echo >&2
    log_info "Downloading SCANOSS Code Compare v$version..."
    download_and_verify_asset "$version" "linux" "$zip_path"

    # Extract the archive
    log_info "Extracting..."
    if command_exists unzip; then
        unzip -q "$zip_path" -d "$temp_dir"
    else
        abort "unzip command not found. Please install unzip and try again."
    fi

    # Find the binary
    local binary_path="$temp_dir/$BINARY_NAME"
    if [ ! -f "$binary_path" ]; then
        abort "Binary not found in archive: $BINARY_NAME"
    fi

    # Make it executable
    chmod +x "$binary_path"

    # Install to system
    log_info "Installing to $INSTALL_DIR..."
    $SUDO install -m 755 "$binary_path" "$INSTALL_DIR/$APP_NAME"

    echo >&2
    log_info "✓ Binary installed successfully"
    echo >&2
    log_info "✓ Installation complete!"
}

# Main function
main() {
    install_linux
    show_completion

    # Verify installation
    echo >&2
    log_info "Verifying installation..."
    if command_exists "$APP_NAME"; then
        "$APP_NAME" --version || log_warn "Could not get version"
        echo >&2
        log_info "Installation verified successfully!"
    else
        log_warn "Command '$APP_NAME' not found in PATH"
        log_warn "You may need to:"
        log_warn "  1. Open a new terminal window"
        log_warn "  2. Check that $INSTALL_DIR is in your PATH"
        echo >&2
        log_info "To add to PATH, run:"
        echo "    export PATH=\"$INSTALL_DIR:\$PATH\"" >&2
        echo >&2
        log_info "To make it permanent, add the above line to your ~/.bashrc or ~/.zshrc"
    fi

    # Show post-installation notes
    echo >&2
    log_info "Post-Installation Notes:"
    echo >&2
    echo "  • GUI Mode: Launch from application menu or run: $APP_NAME" >&2
    echo "  • CLI Mode: Run with arguments: $APP_NAME --help" >&2
    echo >&2
    echo "  • The app requires GTK3 and WebKit2GTK to display GUI dialogs" >&2
    echo "  • If you experience issues, ensure dependencies are installed" >&2
    echo >&2
}

main "$@"
