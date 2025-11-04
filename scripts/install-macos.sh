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


# Direct installation from DMG
install_direct() {
    log_info "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    log_info "  SCANOSS Code Compare - macOS Installation"
    log_info "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo >&2

    setup_sudo
    echo >&2

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
        $SUDO rm -rf "$INSTALL_DIR/$APP_NAME.app"
    fi

    # Copy the app to Applications
    log_info "Installing to $INSTALL_DIR..."
    $SUDO cp -R "$app_path" "$INSTALL_DIR/"

    # Unmount the DMG
    hdiutil detach "$mount_point" >/dev/null 2>&1

    # Create wrapper script for CLI access
    log_info "Creating command-line wrapper..."

    # Create bin directory if it doesn't exist
    if [ ! -d "$BIN_DIR" ]; then
        $SUDO mkdir -p "$BIN_DIR"
    fi

    # Remove existing symlink or wrapper script if present
    if [ -L "$BIN_DIR/$APP_NAME" ] || [ -f "$BIN_DIR/$APP_NAME" ]; then
        $SUDO rm -f "$BIN_DIR/$APP_NAME"
    fi

    # Create wrapper script that launches the app with proper bundle context
    # This ensures macOS GUI features (like file picker dialogs) work correctly
    log_info "Creating wrapper script..."
    WRAPPER_SCRIPT="$BIN_DIR/$APP_NAME"
    $SUDO tee "$WRAPPER_SCRIPT" > /dev/null << 'EOF'
#!/bin/bash

APP_PATH="/Applications/scanoss-cc.app"
CURRENT_DIR=$(pwd)

# Handle version flag directly without launching GUI
if [ "$1" = "--version" ] || [ "$1" = "-v" ]; then
  exec "$APP_PATH/Contents/MacOS/scanoss-cc" "$@"
  exit 0
fi

# If no arguments are provided, launch GUI application
if [ "$#" -eq 0 ]; then
  exec open -a "$APP_PATH"
  exit 0
fi

# If the first argument starts with a dash, assume scanning options
first_arg="$1"
if [[ "$first_arg" == "-"* ]]; then
  # Check if --scan-root is already provided
  found=0
  for arg in "$@"; do
    if [ "$arg" = "--scan-root" ]; then
      found=1
      break
    fi
  done
  if [ $found -eq 0 ]; then
    # Add --scan-root with current directory if not present
    exec open -a "$APP_PATH" --args --scan-root "$CURRENT_DIR" "$@"
  else
    exec open -a "$APP_PATH" --args "$@"
  fi
else
  # Otherwise, assume subcommand mode and pass arguments as-is
  exec open -a "$APP_PATH" --args "$@"
fi
EOF

    # Make the wrapper script executable
    $SUDO chmod +x "$WRAPPER_SCRIPT"
    log_info "Wrapper script created successfully"

    # Set permissions
    $SUDO chmod +x "$INSTALL_DIR/$APP_NAME.app/Contents/MacOS/$APP_NAME"

    # Clear macOS quarantine attribute
    log_info "Clearing quarantine attributes..."
    $SUDO xattr -r -d com.apple.quarantine "$INSTALL_DIR/$APP_NAME.app" 2>/dev/null || true

    echo >&2
    log_info "✓ Installation complete!"
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

    # Always use direct installation
    install_direct

    # Show completion message
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
        log_warn "You may need to open a new terminal window"
    fi
}

main "$@"
