#!/bin/bash
# SCANOSS Code Compare - macOS Installer
# This script installs SCANOSS Code Compare on macOS

set -e

# Repository and app details
readonly REPO="scanoss/scanoss.cc"
readonly APP_NAME="scanoss-cc"
readonly APP_DISPLAY_NAME="SCANOSS Code Compare"
readonly INSTALL_DIR="/Applications"
readonly BIN_DIR="/usr/local/bin"

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

# Check if Homebrew is installed
check_homebrew() {
    if command_exists brew; then
        return 0
    else
        return 1
    fi
}

# Install via Homebrew
install_via_homebrew() {
    log_info "Installing via Homebrew..."
    echo >&2

    if ! check_homebrew; then
        log_warn "Homebrew is not installed."
        echo >&2
        if confirm "Would you like to install Homebrew first?" "n"; then
            log_info "Installing Homebrew..."
            /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)" || abort "Failed to install Homebrew"
        else
            log_error "Homebrew is required for this installation method."
            return 1
        fi
    fi

    log_info "Adding SCANOSS tap..."
    brew tap scanoss/dist || log_warn "Tap might already exist"

    log_info "Installing SCANOSS Code Compare..."
    brew install scanoss-code-compare || abort "Failed to install via Homebrew"

    echo >&2
    log_info "✓ Installation complete via Homebrew!"
    log_info "Homebrew will automatically keep SCANOSS Code Compare updated."
    echo >&2

    return 0
}

# Direct installation from DMG
install_direct() {
    log_info "Direct Installation"
    echo >&2

    # Check if we have write permissions to /Applications and /usr/local/bin
    local need_sudo=false
    if [ ! -w "$INSTALL_DIR" ]; then
        log_warn "No write permission to $INSTALL_DIR"
        need_sudo=true
    fi

    # Check if /usr/local/bin exists and if we can write to it
    if [ -d "$BIN_DIR" ] && [ ! -w "$BIN_DIR" ]; then
        log_warn "No write permission to $BIN_DIR"
        need_sudo=true
    fi

    if [ "$need_sudo" = true ]; then
        log_error "This installation requires administrator privileges."
        echo >&2
        if confirm "Run installation with sudo?" "y"; then
            # Re-run this script with sudo
            exec sudo -E bash "$0" "$@"
        else
            abort "Installation cancelled"
        fi
    fi

    # Setup temporary directory
    local temp_dir=$(setup_temp_dir)
    trap "cleanup_temp_dir '$temp_dir'" EXIT INT TERM

    # Get latest version
    local version=$(get_latest_version "$REPO")

    local zip_path="$temp_dir/scanoss-cc-mac.zip"

    echo >&2
    log_info "Downloading SCANOSS Code Compare v$version..."
    download_and_verify_asset "$version" "macos" "$zip_path"

    # Extract the DMG from the ZIP
    log_info "Extracting DMG..."
    if ! command_exists unzip; then
        abort "unzip command not found. Please install unzip and try again."
    fi

    unzip -q "$zip_path" -d "$temp_dir"

    # Find the extracted DMG
    local dmg_path=$(find "$temp_dir" -name "*.dmg" | head -n 1)

    if [ -z "$dmg_path" ]; then
        abort "Could not find DMG file in the downloaded archive"
    fi

    # Mount the DMG
    log_info "Mounting DMG..."

    # Mount the DMG and capture the mount point
    # hdiutil attach outputs lines like: /dev/disk4s2          Apple_HFS                       /Volumes/scanoss-cc Installer
    local mount_output=$(hdiutil attach "$dmg_path" -nobrowse 2>/dev/null)
    local mount_point=$(echo "$mount_output" | grep "/Volumes/" | sed 's/.*\(\/Volumes\/.*\)/\1/')

    if [ -z "$mount_point" ]; then
        abort "Failed to mount DMG or detect mount point"
    fi

    log_info "DMG mounted at: $mount_point"

    # Wait a moment for the mount to complete
    sleep 2

    # Find the .app in the mounted volume
    local app_path=$(find "$mount_point" -name "$APP_NAME.app" -maxdepth 2 2>/dev/null | head -n 1)

    if [ -z "$app_path" ]; then
        hdiutil detach "$mount_point" >/dev/null 2>&1
        abort "Could not find $APP_NAME.app in the DMG"
    fi

    # Remove existing installation if present
    if [ -d "$INSTALL_DIR/$APP_NAME.app" ]; then
        log_info "Removing existing installation..."
        rm -rf "$INSTALL_DIR/$APP_NAME.app"
    fi

    # Copy the app to Applications
    log_info "Installing to $INSTALL_DIR..."
    cp -R "$app_path" "$INSTALL_DIR/"

    # Unmount the DMG
    hdiutil detach "$mount_point" >/dev/null 2>&1

    # Create symlink for CLI access
    log_info "Creating command-line symlink..."

    # Create bin directory if it doesn't exist
    if [ ! -d "$BIN_DIR" ]; then
        sudo mkdir -p "$BIN_DIR" 2>/dev/null || mkdir -p "$BIN_DIR"
    fi

    # Remove existing symlink if present
    if [ -L "$BIN_DIR/$APP_NAME" ] || [ -f "$BIN_DIR/$APP_NAME" ]; then
        sudo rm -f "$BIN_DIR/$APP_NAME" 2>/dev/null || rm -f "$BIN_DIR/$APP_NAME"
    fi

    # Create new symlink to the binary inside the .app bundle
    if sudo ln -s "$INSTALL_DIR/$APP_NAME.app/Contents/MacOS/$APP_NAME" "$BIN_DIR/$APP_NAME" 2>/dev/null; then
        log_info "Symlink created successfully"
    elif ln -s "$INSTALL_DIR/$APP_NAME.app/Contents/MacOS/$APP_NAME" "$BIN_DIR/$APP_NAME" 2>/dev/null; then
        log_info "Symlink created successfully"
    else
        log_warn "Failed to create symlink at $BIN_DIR/$APP_NAME"
        log_warn "You can manually add to PATH or create the symlink later"
    fi

    # Set permissions
    chmod +x "$INSTALL_DIR/$APP_NAME.app/Contents/MacOS/$APP_NAME"

    # Clear macOS quarantine attribute
    log_info "Clearing quarantine attributes..."
    xattr -r -d com.apple.quarantine "$INSTALL_DIR/$APP_NAME.app" 2>/dev/null || true

    echo >&2
    log_info "✓ Installation complete!"
    echo >&2
}

# Show installation method selection
show_installation_menu() {
    echo >&2
    log_info "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    log_info "  SCANOSS Code Compare - macOS Installation"
    log_info "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo >&2
    echo "Choose installation method:" >&2
    echo >&2
    echo "  1) Homebrew (Recommended)" >&2
    echo "     • Automatic updates" >&2
    echo "     • Easy uninstallation" >&2
    echo "     • Managed by package manager" >&2
    echo >&2
    echo "  2) Direct Installation" >&2
    echo "     • Download and install DMG" >&2
    echo "     • Manual updates" >&2
    echo "     • Full app bundle in /Applications" >&2
    echo >&2
    echo "  3) Cancel" >&2
    echo >&2
}

# Main installation logic
main() {
    # Check if running on macOS
    if [ "$(uname -s)" != "Darwin" ]; then
        abort "This installer is for macOS only. Please use the appropriate installer for your platform."
    fi

    # Check architecture
    local arch=$(uname -m)
    if [ "$arch" != "x86_64" ] && [ "$arch" != "arm64" ]; then
        log_warn "Unsupported architecture: $arch"
        log_warn "This app is built for Intel (x86_64) and Apple Silicon (arm64)"
    fi

    # If running non-interactively or with --homebrew flag, skip menu
    if [ ! -t 0 ] || [ "$1" == "--homebrew" ]; then
        log_info "Running in non-interactive mode (installing via Homebrew)"
        echo >&2
        install_via_homebrew
        show_completion
        return 0
    fi

    # If running with --direct flag, skip menu
    if [ "$1" == "--direct" ]; then
        install_direct
        show_completion
        return 0
    fi

    # Show menu
    show_installation_menu

    # Get user choice
    local choice
    while true; do
        read -p "Enter your choice (1-3): " choice </dev/tty
        case $choice in
            1)
                if install_via_homebrew; then
                    break
                else
                    echo >&2
                    log_error "Homebrew installation failed or was cancelled."
                    if confirm "Try direct installation instead?" "y"; then
                        install_direct
                        break
                    else
                        abort "Installation cancelled"
                    fi
                fi
                ;;
            2)
                install_direct
                break
                ;;
            3)
                log_info "Installation cancelled by user"
                exit 0
                ;;
            *)
                log_error "Invalid choice. Please enter 1, 2, or 3."
                ;;
        esac
    done

    # Show completion message
    show_completion

    # Verify installation
    echo >&2
    log_info "Verifying installation..."
    if command_exists "$APP_NAME"; then
        "$APP_NAME" --version || log_warn "Could not get version"
    else
        log_warn "Command '$APP_NAME' not found in PATH"
        log_warn "You may need to open a new terminal window"
    fi

    # Offer to open the app
    echo >&2
    if confirm "Would you like to open SCANOSS Code Compare now?" "y"; then
        log_info "Opening $APP_DISPLAY_NAME..."
        open -a "$APP_NAME" 2>/dev/null || log_warn "Could not open app automatically"
    fi
}

# Prevent partial execution if script is interrupted during download
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
